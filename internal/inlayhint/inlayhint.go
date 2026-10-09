package inlayhint

import (
	"strings"

	"go.lsp.dev/protocol"

	"github.com/ts91/urscript-ls/internal/builtins"
	"github.com/ts91/urscript-ls/internal/parser"
	"github.com/ts91/urscript-ls/internal/symbols"
)

// Provide returns inlay hints for parameter names at call sites.
func Provide(prog *parser.Program, rng protocol.Range, registry *builtins.Registry, idx *symbols.Index) []protocol.InlayHint {
	var hints []protocol.InlayHint
	walkStmts(prog.Stmts, rng, registry, idx, &hints)
	return hints
}

func walkStmts(stmts []parser.Node, rng protocol.Range, registry *builtins.Registry, idx *symbols.Index, hints *[]protocol.InlayHint) {
	for _, stmt := range stmts {
		walkNode(stmt, rng, registry, idx, hints)
	}
}

func walkNode(n parser.Node, rng protocol.Range, registry *builtins.Registry, idx *symbols.Index, hints *[]protocol.InlayHint) {
	if n == nil {
		return
	}
	switch node := n.(type) {
	case *parser.FunctionDef:
		walkStmts(node.Body, rng, registry, idx, hints)
	case *parser.ThreadDef:
		walkStmts(node.Body, rng, registry, idx, hints)
	case *parser.IfStmt:
		walkNode(node.Condition, rng, registry, idx, hints)
		walkStmts(node.Body, rng, registry, idx, hints)
		for _, elif := range node.ElseIfs {
			walkNode(elif.Condition, rng, registry, idx, hints)
			walkStmts(elif.Body, rng, registry, idx, hints)
		}
		walkStmts(node.Else, rng, registry, idx, hints)
	case *parser.WhileStmt:
		walkNode(node.Condition, rng, registry, idx, hints)
		walkStmts(node.Body, rng, registry, idx, hints)
	case *parser.ForStmt:
		walkNode(node.Iter, rng, registry, idx, hints)
		walkStmts(node.Body, rng, registry, idx, hints)
	case *parser.ReturnStmt:
		walkNode(node.Value, rng, registry, idx, hints)
	case *parser.Assignment:
		walkNode(node.Value, rng, registry, idx, hints)
	case *parser.ExprStmt:
		walkNode(node.Expr, rng, registry, idx, hints)
	case *parser.GlobalStmt:
		walkNode(node.Init, rng, registry, idx, hints)
	case *parser.LocalStmt:
		walkNode(node.Init, rng, registry, idx, hints)
	case *parser.BinaryExpr:
		walkNode(node.Left, rng, registry, idx, hints)
		walkNode(node.Right, rng, registry, idx, hints)
	case *parser.UnaryExpr:
		walkNode(node.Operand, rng, registry, idx, hints)
	case *parser.IndexExpr:
		walkNode(node.Expr, rng, registry, idx, hints)
		walkNode(node.Index, rng, registry, idx, hints)
	case *parser.ListLit:
		for _, el := range node.Elements {
			walkNode(el, rng, registry, idx, hints)
		}
	case *parser.CallExpr:
		processCall(node, rng, registry, idx, hints)
		for _, arg := range node.Args {
			walkNode(arg, rng, registry, idx, hints)
		}
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
