package diagnostics

import (
	"fmt"
	"strings"

	"go.lsp.dev/protocol"

	"github.com/ts91/urscript-ls/internal/builtins"
	"github.com/ts91/urscript-ls/internal/parser"
	"github.com/ts91/urscript-ls/internal/symbols"
)

// Analyze runs semantic checks on the AST and returns diagnostics.
func Analyze(prog *parser.Program, idx *symbols.Index, registry *builtins.Registry) []protocol.Diagnostic {
	a := &analyzer{
		idx:      idx,
		registry: registry,
		defined:  make(map[string]bool),
	}

	// Collect all known names (builtins, user symbols).
	for _, b := range registry.All() {
		a.defined[b.Name] = true
	}
	if idx != nil {
		for _, s := range idx.Symbols {
			a.defined[s.Name] = true
		}
	}

	a.checkDuplicateDefs()
	a.walkStmts(prog.Stmts)
	return a.diags
}

type analyzer struct {
	idx      *symbols.Index
	registry *builtins.Registry
	defined  map[string]bool
	diags    []protocol.Diagnostic
}

func (a *analyzer) warn(pos parser.Pos, msg string) {
	a.diags = append(a.diags, protocol.Diagnostic{
		Range: protocol.Range{
			Start: protocol.Position{Line: uint32(pos.Line), Character: uint32(pos.Column)},
			End:   protocol.Position{Line: uint32(pos.Line), Character: uint32(pos.Column)},
		},
		Severity: protocol.DiagnosticSeverityWarning,
		Source:   protocol.NewOptional[string]("urscript-ls"),
		Message:  protocol.String(msg),
	})
}

func (a *analyzer) checkDuplicateDefs() {
	if a.idx == nil {
		return
	}
	seen := make(map[string]parser.Pos)
	for _, s := range a.idx.Symbols {
		if s.Kind != symbols.KindFunction && s.Kind != symbols.KindThread {
			continue
		}
		if prev, ok := seen[s.Name]; ok {
			a.warn(s.Pos, fmt.Sprintf("'%s' already defined at line %d", s.Name, prev.Line+1))
		} else {
			seen[s.Name] = s.Pos
		}
	}
}

func (a *analyzer) walkStmts(stmts []parser.Node) {
	for _, stmt := range stmts {
		a.walkNode(stmt)
	}
}

func (a *analyzer) walkNode(n parser.Node) {
	switch node := n.(type) {
	case *parser.FunctionDef:
		a.walkStmts(node.Body)
	case *parser.ThreadDef:
		a.walkStmts(node.Body)
	case *parser.IfStmt:
		a.walkExpr(node.Condition)
		a.walkStmts(node.Body)
		for _, elif := range node.ElseIfs {
			a.walkExpr(elif.Condition)
			a.walkStmts(elif.Body)
		}
		a.walkStmts(node.Else)
	case *parser.WhileStmt:
		a.walkExpr(node.Condition)
		a.walkStmts(node.Body)
	case *parser.ForStmt:
		a.walkExpr(node.Iter)
		a.walkStmts(node.Body)
	case *parser.ReturnStmt:
		if node.Value != nil {
			a.walkExpr(node.Value)
		}
	case *parser.Assignment:
		a.walkExpr(node.Value)
	case *parser.ExprStmt:
		a.walkExpr(node.Expr)
	case *parser.GlobalStmt:
		if node.Init != nil {
			a.walkExpr(node.Init)
		}
	case *parser.LocalStmt:
		if node.Init != nil {
			a.walkExpr(node.Init)
		}
	}
}

func (a *analyzer) walkExpr(n parser.Node) {
	if n == nil {
		return
	}
	switch node := n.(type) {
	case *parser.CallExpr:
		a.checkCall(node)
		for _, arg := range node.Args {
			a.walkExpr(arg)
		}
	case *parser.BinaryExpr:
		a.walkExpr(node.Left)
		a.walkExpr(node.Right)
	case *parser.UnaryExpr:
		a.walkExpr(node.Operand)
	case *parser.IndexExpr:
		a.walkExpr(node.Expr)
		a.walkExpr(node.Index)
	case *parser.ListLit:
		for _, el := range node.Elements {
			a.walkExpr(el)
		}
	case *parser.Assignment:
		a.walkExpr(node.Value)
	case *parser.Identifier:
		// Don't warn about identifiers — too noisy without proper scoping.
	}
}

func (a *analyzer) checkCall(call *parser.CallExpr) {
	id, ok := call.Func.(*parser.Identifier)
	if !ok {
		return
	}
	name := id.Name

	// Check argument count against builtins.
	if b, ok := a.registry.Lookup(name); ok {
		expected := countParams(b.Signature)
		actual := len(call.Args)
		// Only warn when we have too many args (keyword args make "too few" unreliable).
		if expected >= 0 && actual > expected {
			a.warn(call.Pos, fmt.Sprintf("'%s' expects at most %d arguments, got %d", name, expected, actual))
		}
		return
	}

	// Check against user-defined functions.
	if a.idx != nil {
		for _, s := range a.idx.Lookup(name) {
			if s.Kind == symbols.KindFunction {
				expected := len(s.Params)
				actual := len(call.Args)
				if actual > expected {
					a.warn(call.Pos, fmt.Sprintf("'%s' expects at most %d arguments, got %d", name, expected, actual))
				}
				return
			}
		}
	}

	// Unknown function call.
	if !a.defined[name] {
		a.warn(id.Pos, fmt.Sprintf("'%s' is not defined", name))
	}
}

// countParams parses parameter count from a signature like "func(a, b, c=1)".
func countParams(sig string) int {
	start := strings.Index(sig, "(")
	end := strings.LastIndex(sig, ")")
	if start < 0 || end < 0 || end <= start+1 {
		return 0
	}
	body := strings.TrimSpace(sig[start+1 : end])
	if body == "" {
		return 0
	}
	return strings.Count(body, ",") + 1
}
