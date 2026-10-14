package lsp

import (
	"strings"

	"go.lsp.dev/protocol"

	"github.com/ts91/urscript-ls/internal/builtins"
	"github.com/ts91/urscript-ls/internal/parser"
	"github.com/ts91/urscript-ls/internal/symbols"
)

// InlayHints returns inlay hints for parameter names at call sites.
func InlayHints(prog *parser.Program, rng protocol.Range, registry *builtins.Registry, idx *symbols.Index) (hints []protocol.InlayHint) {
	if prog == nil {
		return nil
	}
	defer func() {
		if r := recover(); r != nil {
			hints = nil
		}
	}()
	inlayWalkStmts(prog.Stmts, rng, registry, idx, &hints)
	return hints
}

func inlayWalkStmts(stmts []parser.Node, rng protocol.Range, registry *builtins.Registry, idx *symbols.Index, hints *[]protocol.InlayHint) {
	for _, stmt := range stmts {
		inlayWalkNode(stmt, rng, registry, idx, hints)
	}
}

func inlayWalkNode(n parser.Node, rng protocol.Range, registry *builtins.Registry, idx *symbols.Index, hints *[]protocol.InlayHint) {
	if n == nil {
		return
	}
	switch node := n.(type) {
	case *parser.FunctionDef:
		inlayWalkStmts(node.Body, rng, registry, idx, hints)
	case *parser.ThreadDef:
		inlayWalkStmts(node.Body, rng, registry, idx, hints)
	case *parser.IfStmt:
		inlayWalkNode(node.Condition, rng, registry, idx, hints)
		inlayWalkStmts(node.Body, rng, registry, idx, hints)
		for _, elif := range node.ElseIfs {
			inlayWalkNode(elif.Condition, rng, registry, idx, hints)
			inlayWalkStmts(elif.Body, rng, registry, idx, hints)
		}
		inlayWalkStmts(node.Else, rng, registry, idx, hints)
	case *parser.WhileStmt:
		inlayWalkNode(node.Condition, rng, registry, idx, hints)
		inlayWalkStmts(node.Body, rng, registry, idx, hints)
	case *parser.ReturnStmt:
		inlayWalkNode(node.Value, rng, registry, idx, hints)
	case *parser.Assignment:
		inlayWalkNode(node.Value, rng, registry, idx, hints)
	case *parser.ExprStmt:
		inlayWalkNode(node.Expr, rng, registry, idx, hints)
	case *parser.GlobalStmt:
		inlayWalkNode(node.Init, rng, registry, idx, hints)
	case *parser.LocalStmt:
		inlayWalkNode(node.Init, rng, registry, idx, hints)
	case *parser.BinaryExpr:
		inlayWalkNode(node.Left, rng, registry, idx, hints)
		inlayWalkNode(node.Right, rng, registry, idx, hints)
	case *parser.UnaryExpr:
		inlayWalkNode(node.Operand, rng, registry, idx, hints)
	case *parser.IndexExpr:
		inlayWalkNode(node.Expr, rng, registry, idx, hints)
		inlayWalkNode(node.Index, rng, registry, idx, hints)
	case *parser.ListLit:
		for _, el := range node.Elements {
			inlayWalkNode(el, rng, registry, idx, hints)
		}
	case *parser.CallExpr:
		processCall(node, rng, registry, idx, hints)
		for _, arg := range node.Args {
			inlayWalkNode(arg, rng, registry, idx, hints)
		}
	case *parser.HaltStmt:
		// Nothing to walk.
	}
}

func processCall(call *parser.CallExpr, rng protocol.Range, registry *builtins.Registry, idx *symbols.Index, hints *[]protocol.InlayHint) {
	pos := parser.NodePos(call)
	line := uint32(pos.Line)
	if line < rng.Start.Line || line > rng.End.Line {
		return
	}

	id, ok := call.Func.(*parser.Identifier)
	if !ok {
		return
	}

	params := lookupParams(id.Name, registry, idx)
	if len(params) == 0 {
		return
	}

	kind := protocol.InlayHintKindParameter
	for i, arg := range call.Args {
		if i >= len(params) {
			break
		}
		if arg == nil {
			continue
		}
		// Skip keyword arguments — they already show the name.
		if isKeywordArg(arg) {
			continue
		}
		argPos := parser.NodePos(arg)
		*hints = append(*hints, protocol.InlayHint{
			Position: protocol.Position{Line: uint32(argPos.Line), Character: uint32(argPos.Column)},
			Label:    protocol.String(params[i] + ":"),
			Kind:     kind,
		})
	}
}

func lookupParams(name string, registry *builtins.Registry, idx *symbols.Index) []string {
	if b, ok := registry.Lookup(name); ok {
		return extractParams(b.Signature)
	}
	if idx != nil {
		for _, s := range idx.Lookup(name) {
			if s.Kind == symbols.KindFunction {
				return s.Params
			}
		}
	}
	return nil
}

func extractParams(sig string) []string {
	start := strings.Index(sig, "(")
	end := strings.LastIndex(sig, ")")
	if start < 0 || end < 0 || end <= start+1 {
		return nil
	}
	body := sig[start+1 : end]
	parts := strings.Split(body, ",")
	params := make([]string, len(parts))
	for i, p := range parts {
		p = strings.TrimSpace(p)
		if eq := strings.Index(p, "="); eq >= 0 {
			p = strings.TrimSpace(p[:eq])
		}
		params[i] = p
	}
	return params
}

func isKeywordArg(n parser.Node) bool {
	if bin, ok := n.(*parser.BinaryExpr); ok && bin.Op == "=" {
		if _, ok := bin.Left.(*parser.Identifier); ok {
			return true
		}
	}
	return false
}
