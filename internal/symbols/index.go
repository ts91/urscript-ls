package symbols

import (
	"github.com/ts91/urscript-ls/internal/parser"
)

// SymbolKind classifies a symbol.
type SymbolKind int

const (
	KindFunction SymbolKind = iota
	KindVariable
	KindGlobal
	KindLocal
	KindThread
	KindParam
)

// Symbol represents a named entity defined in the source.
type Symbol struct {
	Name   string
	Kind   SymbolKind
	Pos    parser.Pos
	EndPos parser.Pos // end of the definition block (for functions/threads)

	// Function-specific fields.
	Params []string
	Doc    string // preceding comment block, if any
}

// Index holds all symbols extracted from a single document.
type Index struct {
	Symbols []Symbol
	byName  map[string][]Symbol
}

// Lookup returns all symbols with the given name.
func (idx *Index) Lookup(name string) []Symbol {
	if idx == nil || idx.byName == nil {
		return nil
	}
	return idx.byName[name]
}

// Functions returns all function symbols.
func (idx *Index) Functions() []Symbol {
	return idx.filter(KindFunction)
}

// Variables returns all variable symbols (excluding params).
func (idx *Index) Variables() []Symbol {
	var out []Symbol
	for _, s := range idx.Symbols {
		if s.Kind == KindVariable || s.Kind == KindGlobal || s.Kind == KindLocal {
			out = append(out, s)
		}
	}
	return out
}

// Threads returns all thread symbols.
func (idx *Index) Threads() []Symbol {
	return idx.filter(KindThread)
}

func (idx *Index) filter(kind SymbolKind) []Symbol {
	var out []Symbol
	for _, s := range idx.Symbols {
		if s.Kind == kind {
			out = append(out, s)
		}
	}
	return out
}

// Build walks the AST and extracts all symbols.
func Build(prog *parser.Program, source string) *Index {
	idx := &Index{
		byName: make(map[string][]Symbol),
	}
	b := &builder{idx: idx, source: source}
	b.walkStmts(prog.Stmts)
	return idx
}

func (idx *Index) add(s Symbol) {
	idx.Symbols = append(idx.Symbols, s)
	idx.byName[s.Name] = append(idx.byName[s.Name], s)
}

type builder struct {
	idx    *Index
	source string
}

func (b *builder) walkStmts(stmts []parser.Node) {
	for _, stmt := range stmts {
		b.walkNode(stmt)
	}
}

func (b *builder) walkNode(n parser.Node) {
	switch node := n.(type) {
	case *parser.FunctionDef:
		endPos := node.Pos
		if len(node.Body) > 0 {
			endPos = lastStmtPos(node.Body)
		}
		b.idx.add(Symbol{
			Name:   node.Name,
			Kind:   KindFunction,
			Pos:    node.Pos,
			EndPos: endPos,
			Params: node.Params,
			Doc:    extractDoc(b.source, node.Pos.Line),
		})
		for _, p := range node.Params {
			b.idx.add(Symbol{
				Name: p,
				Kind: KindParam,
				Pos:  node.Pos,
			})
		}
		b.walkStmts(node.Body)

	case *parser.ThreadDef:
		endPos := node.Pos
		if len(node.Body) > 0 {
			endPos = lastStmtPos(node.Body)
		}
		b.idx.add(Symbol{
			Name:   node.Name,
			Kind:   KindThread,
			Pos:    node.Pos,
			EndPos: endPos,
			Doc:    extractDoc(b.source, node.Pos.Line),
		})
		b.walkStmts(node.Body)

	case *parser.GlobalStmt:
		b.idx.add(Symbol{
			Name: node.Name,
			Kind: KindGlobal,
			Pos:  node.Pos,
		})

	case *parser.LocalStmt:
		b.idx.add(Symbol{
			Name: node.Name,
			Kind: KindLocal,
			Pos:  node.Pos,
		})

	case *parser.Assignment:
		if id, ok := node.Target.(*parser.Identifier); ok {
			b.idx.add(Symbol{
				Name: id.Name,
				Kind: KindVariable,
				Pos:  id.Pos,
			})
		}

	case *parser.IfStmt:
		b.walkStmts(node.Body)
		for _, elif := range node.ElseIfs {
			b.walkStmts(elif.Body)
		}
		b.walkStmts(node.Else)

	case *parser.WhileStmt:
		b.walkStmts(node.Body)

	default:
		// ExprStmt, ReturnStmt, etc. — no symbols to extract
	}
}

// extractDoc collects documentation above defLine.
// It first checks for a ### block docstring, then falls back to
// contiguous # line comments.
func extractDoc(source string, defLine int) string {
	lines := splitLines(source)
	if defLine <= 0 || defLine > len(lines) {
		return ""
	}

	// Try ### block docstring first.
	if doc := extractBlockDoc(lines, defLine); doc != "" {
		return doc
	}

	// Fall back to contiguous # line comments.
	return extractLineDoc(lines, defLine)
}

// extractBlockDoc looks for a ### ... ### block immediately above defLine.
func extractBlockDoc(lines []string, defLine int) string {
	// Find the closing ### (walking up, skipping blank lines).
	closeIdx := -1
	for i := defLine - 1; i >= 0; i-- {
		trimmed := trimSpace(lines[i])
		if len(trimmed) == 0 {
			continue
		}
		if trimmed == "###" {
			closeIdx = i
		}
		break
	}
	if closeIdx < 0 {
		return ""
	}

	// Find the matching opening ###.
	openIdx := -1
	for i := closeIdx - 1; i >= 0; i-- {
		if trimSpace(lines[i]) == "###" {
			openIdx = i
			break
		}
	}
	if openIdx < 0 {
		return ""
	}

	// Collect the content between the delimiters, stripping leading "# " or "#".
	var docLines []string
	for i := openIdx + 1; i < closeIdx; i++ {
		line := trimSpace(lines[i])
		if len(line) > 0 && line[0] == '#' {
			line = line[1:]
			if len(line) > 0 && line[0] == ' ' {
				line = line[1:]
			}
		}
		docLines = append(docLines, line)
	}

	result := ""
	for i, line := range docLines {
		if i > 0 {
			result += "\n"
		}
		result += line
	}
	return result
}

// extractLineDoc collects contiguous # comment lines above defLine.
func extractLineDoc(lines []string, defLine int) string {
	var docLines []string
	for i := defLine - 1; i >= 0; i-- {
		trimmed := trimSpace(lines[i])
		if len(trimmed) == 0 {
			break
		}
		if trimmed[0] != '#' {
			break
		}
		// Stop at ### block delimiters.
		if trimmed == "###" {
			break
		}
		// Strip leading "# " or "#"
		comment := trimmed[1:]
		if len(comment) > 0 && comment[0] == ' ' {
			comment = comment[1:]
		}
		docLines = append(docLines, comment)
	}

	// Reverse to get top-to-bottom order.
	for i, j := 0, len(docLines)-1; i < j; i, j = i+1, j-1 {
		docLines[i], docLines[j] = docLines[j], docLines[i]
	}

	result := ""
	for i, line := range docLines {
		if i > 0 {
			result += "\n"
		}
		result += line
	}
	return result
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	lines = append(lines, s[start:])
	return lines
}

// lastStmtPos extracts the position from the last statement in a body.
func lastStmtPos(stmts []parser.Node) parser.Pos {
	return parser.NodePos(stmts[len(stmts)-1])
}

// Merge combines multiple indices into one, giving priority to the primary index.
func Merge(primary *Index, others ...*Index) *Index {
	merged := &Index{
		byName: make(map[string][]Symbol),
	}
	if primary != nil {
		for _, s := range primary.Symbols {
			merged.add(s)
		}
	}
	seen := make(map[string]bool)
	for _, s := range merged.Symbols {
		if s.Kind == KindFunction || s.Kind == KindThread {
			seen[s.Name] = true
		}
	}
	for _, idx := range others {
		if idx == nil {
			continue
		}
		for _, s := range idx.Symbols {
			if s.Kind == KindParam {
				continue
			}
			// Don't duplicate functions/threads already in the primary file.
			if (s.Kind == KindFunction || s.Kind == KindThread) && seen[s.Name] {
				continue
			}
			merged.add(s)
			if s.Kind == KindFunction || s.Kind == KindThread {
				seen[s.Name] = true
			}
		}
	}
	return merged
}

func trimSpace(s string) string {
	i := 0
	for i < len(s) && (s[i] == ' ' || s[i] == '\t' || s[i] == '\r') {
		i++
	}
	j := len(s)
	for j > i && (s[j-1] == ' ' || s[j-1] == '\t' || s[j-1] == '\r') {
		j--
	}
	return s[i:j]
}
