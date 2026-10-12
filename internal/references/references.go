package references

import (
	"strings"
	"unicode"

	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"

	"github.com/ts91/urscript-ls/internal/parser"
)

// Pos is a found occurrence of an identifier.
type Pos struct {
	Line   int
	Column int
	Len    int
}

// FindAll walks the AST and returns every position where name appears as an identifier.
func FindAll(prog *parser.Program, name string) []Pos {
	f := &finder{name: name}
	f.walkStmts(prog.Stmts)
	return f.results
}

// FindAllInSource uses simple text scanning to find all identifier occurrences
// of name. This catches occurrences the AST walker might miss (e.g. in
// keyword argument names).
func FindAllInSource(source, name string) []Pos {
	var results []Pos
	lines := strings.Split(source, "\n")
	nameLen := len(name)

	for lineNo, line := range lines {
		off := 0
		for off < len(line) {
			idx := strings.Index(line[off:], name)
			if idx < 0 {
				break
			}
			col := off + idx
			// Ensure it's a whole-word match.
			before := col - 1
			after := col + nameLen
			wordBefore := before >= 0 && isIdent(rune(line[before]))
			wordAfter := after < len(line) && isIdent(rune(line[after]))
			if !wordBefore && !wordAfter {
				results = append(results, Pos{Line: lineNo, Column: col, Len: nameLen})
			}
			off = col + nameLen
		}
	}
	return results
}

// ToLocations converts found positions to LSP locations.
func ToLocations(positions []Pos, docURI uri.URI) []protocol.Location {
	locs := make([]protocol.Location, len(positions))
	for i, p := range positions {
		locs[i] = protocol.Location{
			URI: docURI,
			Range: protocol.Range{
				Start: protocol.Position{Line: uint32(p.Line), Character: uint32(p.Column)},
				End:   protocol.Position{Line: uint32(p.Line), Character: uint32(p.Column + p.Len)},
			},
		}
	}
	return locs
}

// ToTextEdits creates text edits to rename all occurrences.
func ToTextEdits(positions []Pos, newName string) []protocol.TextEdit {
	edits := make([]protocol.TextEdit, len(positions))
	for i, p := range positions {
		edits[i] = protocol.TextEdit{
			Range: protocol.Range{
				Start: protocol.Position{Line: uint32(p.Line), Character: uint32(p.Column)},
				End:   protocol.Position{Line: uint32(p.Line), Character: uint32(p.Column + p.Len)},
			},
			NewText: newName,
		}
	}
	return edits
}

func isIdent(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

type finder struct {
	name    string
	results []Pos
}

func (f *finder) add(pos parser.Pos, nameLen int) {
	f.results = append(f.results, Pos{Line: pos.Line, Column: pos.Column, Len: nameLen})
}

func (f *finder) walkStmts(stmts []parser.Node) {
	for _, stmt := range stmts {
		f.walkNode(stmt)
	}
}

func (f *finder) walkNode(n parser.Node) {
	if n == nil {
		return
	}
	switch node := n.(type) {
	case *parser.FunctionDef:
		if node.Name == f.name {
			// Position is at the 'def' keyword; name starts after "def ".
			f.add(parser.Pos{Line: node.Pos.Line, Column: node.Pos.Column + 4}, len(node.Name))
		}
		f.walkStmts(node.Body)

	case *parser.ThreadDef:
		if node.Name == f.name {
			f.add(parser.Pos{Line: node.Pos.Line, Column: node.Pos.Column + 7}, len(node.Name))
		}
		f.walkStmts(node.Body)

	case *parser.IfStmt:
		f.walkNode(node.Condition)
		f.walkStmts(node.Body)
		for _, elif := range node.ElseIfs {
			f.walkNode(elif.Condition)
			f.walkStmts(elif.Body)
		}
		f.walkStmts(node.Else)

	case *parser.WhileStmt:
		f.walkNode(node.Condition)
		f.walkStmts(node.Body)

	case *parser.ReturnStmt:
		f.walkNode(node.Value)

	case *parser.Assignment:
		f.walkNode(node.Target)
		f.walkNode(node.Value)

	case *parser.ExprStmt:
		f.walkNode(node.Expr)

	case *parser.GlobalStmt:
		if node.Name == f.name {
			f.add(parser.Pos{Line: node.Pos.Line, Column: node.Pos.Column + 7}, len(node.Name))
		}
		f.walkNode(node.Init)

	case *parser.LocalStmt:
		if node.Name == f.name {
			f.add(parser.Pos{Line: node.Pos.Line, Column: node.Pos.Column + 6}, len(node.Name))
		}
		f.walkNode(node.Init)

	case *parser.Identifier:
		if node.Name == f.name {
			f.add(node.Pos, len(node.Name))
		}

	case *parser.CallExpr:
		f.walkNode(node.Func)
		for _, arg := range node.Args {
			f.walkNode(arg)
		}

	case *parser.BinaryExpr:
		f.walkNode(node.Left)
		f.walkNode(node.Right)

	case *parser.UnaryExpr:
		f.walkNode(node.Operand)

	case *parser.IndexExpr:
		f.walkNode(node.Expr)
		f.walkNode(node.Index)

	case *parser.ListLit:
		for _, el := range node.Elements {
			f.walkNode(el)
		}
	}
}
