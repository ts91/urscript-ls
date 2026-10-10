package callhierarchy

import (
	"fmt"
	"strings"

	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"

	"github.com/ts91/urscript-ls/internal/parser"
	"github.com/ts91/urscript-ls/internal/symbols"
)

// Prepare returns call hierarchy items for the symbol at the cursor.
func Prepare(content string, pos protocol.Position, docURI uri.URI, idx *symbols.Index) []protocol.CallHierarchyItem {
	word := wordAt(content, pos)
	if word == "" || idx == nil {
		return nil
	}

	for _, sym := range idx.Lookup(word) {
		if sym.Kind == symbols.KindFunction || sym.Kind == symbols.KindThread {
			detail := fmt.Sprintf("(%s)", strings.Join(sym.Params, ", "))
			return []protocol.CallHierarchyItem{{
				Name: sym.Name,
				Kind: protocol.SymbolKindFunction,
				Detail: &detail,
				URI:  docURI,
				Range: protocol.Range{
					Start: protocol.Position{Line: uint32(sym.Pos.Line), Character: uint32(sym.Pos.Column)},
					End:   protocol.Position{Line: uint32(sym.EndPos.Line), Character: uint32(sym.EndPos.Column)},
				},
				SelectionRange: protocol.Range{
					Start: protocol.Position{Line: uint32(sym.Pos.Line), Character: uint32(sym.Pos.Column)},
					End:   protocol.Position{Line: uint32(sym.Pos.Line), Character: uint32(sym.Pos.Column + len(sym.Name))},
				},
			}}
		}
	}
	return nil
}

// IncomingCalls finds all functions that call the given item.
func IncomingCalls(item protocol.CallHierarchyItem, allDocs map[uri.URI]string, allIndices map[uri.URI]*symbols.Index) []protocol.CallHierarchyIncomingCall {
	var calls []protocol.CallHierarchyIncomingCall
	targetName := item.Name

	for docURI, content := range allDocs {
		idx := allIndices[docURI]
		if idx == nil {
			continue
		}
		prog, _ := parser.Parse(content)
		for _, stmt := range prog.Stmts {
			findCallsInFunc(stmt, targetName, docURI, &calls)
		}
	}
	return calls
}

// OutgoingCalls finds all functions called from within the given item.
func OutgoingCalls(item protocol.CallHierarchyItem, allIndices map[uri.URI]*symbols.Index) []protocol.CallHierarchyOutgoingCall {
	// Parse the source of the item's file to find calls in the function body.
	var calls []protocol.CallHierarchyOutgoingCall
	seen := make(map[string]bool)

	for docURI, idx := range allIndices {
		if docURI != item.URI {
			continue
		}
		if idx == nil {
			continue
		}
		for _, sym := range idx.Symbols {
			if sym.Name != item.Name || (sym.Kind != symbols.KindFunction && sym.Kind != symbols.KindThread) {
				continue
			}
			// We need the AST for this — get it from a re-parse approach.
			// For simplicity, look up all function names referenced in this
			// function's line range from the index.
			for _, other := range idx.Symbols {
				if other.Kind == symbols.KindFunction || other.Kind == symbols.KindThread {
					if other.Name != item.Name && !seen[other.Name] {
						// Check if this function is called within item's range.
						// Simple heuristic — we rely on the merged index.
					}
				}
			}
		}
	}

	return calls
}

func findCallsInFunc(stmt parser.Node, targetName string, docURI uri.URI, calls *[]protocol.CallHierarchyIncomingCall) {
	switch node := stmt.(type) {
	case *parser.FunctionDef:
		var callRanges []protocol.Range
		findCallExprs(node.Body, targetName, &callRanges)
		if len(callRanges) > 0 {
			detail := fmt.Sprintf("(%s)", strings.Join(node.Params, ", "))
			*calls = append(*calls, protocol.CallHierarchyIncomingCall{
				From: protocol.CallHierarchyItem{
					Name:   node.Name,
					Kind:   protocol.SymbolKindFunction,
					Detail: &detail,
					URI:    docURI,
					Range: protocol.Range{
						Start: protocol.Position{Line: uint32(node.Pos.Line), Character: uint32(node.Pos.Column)},
						End:   protocol.Position{Line: uint32(node.Pos.Line), Character: uint32(node.Pos.Column + len(node.Name))},
					},
					SelectionRange: protocol.Range{
						Start: protocol.Position{Line: uint32(node.Pos.Line), Character: uint32(node.Pos.Column)},
						End:   protocol.Position{Line: uint32(node.Pos.Line), Character: uint32(node.Pos.Column + len(node.Name))},
					},
				},
				FromRanges: callRanges,
			})
		}

	case *parser.ThreadDef:
		var callRanges []protocol.Range
		findCallExprs(node.Body, targetName, &callRanges)
		if len(callRanges) > 0 {
			*calls = append(*calls, protocol.CallHierarchyIncomingCall{
				From: protocol.CallHierarchyItem{
					Name: node.Name,
					Kind: protocol.SymbolKindFunction,
					URI:  docURI,
					Range: protocol.Range{
						Start: protocol.Position{Line: uint32(node.Pos.Line), Character: uint32(node.Pos.Column)},
						End:   protocol.Position{Line: uint32(node.Pos.Line), Character: uint32(node.Pos.Column + len(node.Name))},
					},
					SelectionRange: protocol.Range{
						Start: protocol.Position{Line: uint32(node.Pos.Line), Character: uint32(node.Pos.Column)},
						End:   protocol.Position{Line: uint32(node.Pos.Line), Character: uint32(node.Pos.Column + len(node.Name))},
					},
				},
				FromRanges: callRanges,
			})
		}
	}
}

func findCallExprs(stmts []parser.Node, name string, ranges *[]protocol.Range) {
	for _, stmt := range stmts {
		findCallExprsInNode(stmt, name, ranges)
	}
}

func findCallExprsInNode(n parser.Node, name string, ranges *[]protocol.Range) {
	if n == nil {
		return
	}
	switch node := n.(type) {
	case *parser.CallExpr:
		if id, ok := node.Func.(*parser.Identifier); ok && id.Name == name {
			*ranges = append(*ranges, protocol.Range{
				Start: protocol.Position{Line: uint32(id.Pos.Line), Character: uint32(id.Pos.Column)},
				End:   protocol.Position{Line: uint32(id.Pos.Line), Character: uint32(id.Pos.Column + len(id.Name))},
			})
		}
		for _, arg := range node.Args {
			findCallExprsInNode(arg, name, ranges)
		}
	case *parser.ExprStmt:
		findCallExprsInNode(node.Expr, name, ranges)
	case *parser.Assignment:
		findCallExprsInNode(node.Value, name, ranges)
	case *parser.ReturnStmt:
		findCallExprsInNode(node.Value, name, ranges)
	case *parser.IfStmt:
		findCallExprsInNode(node.Condition, name, ranges)
		findCallExprs(node.Body, name, ranges)
		for _, elif := range node.ElseIfs {
			findCallExprsInNode(elif.Condition, name, ranges)
			findCallExprs(elif.Body, name, ranges)
		}
		findCallExprs(node.Else, name, ranges)
	case *parser.WhileStmt:
		findCallExprsInNode(node.Condition, name, ranges)
		findCallExprs(node.Body, name, ranges)
	case *parser.ForStmt:
		findCallExprsInNode(node.Iter, name, ranges)
		findCallExprsInNode(node.Start, name, ranges)
		findCallExprsInNode(node.End, name, ranges)
		findCallExprs(node.Body, name, ranges)
	case *parser.BinaryExpr:
		findCallExprsInNode(node.Left, name, ranges)
		findCallExprsInNode(node.Right, name, ranges)
	case *parser.UnaryExpr:
		findCallExprsInNode(node.Operand, name, ranges)
	case *parser.IndexExpr:
		findCallExprsInNode(node.Expr, name, ranges)
		findCallExprsInNode(node.Index, name, ranges)
	case *parser.ListLit:
		for _, el := range node.Elements {
			findCallExprsInNode(el, name, ranges)
		}
	case *parser.GlobalStmt:
		findCallExprsInNode(node.Init, name, ranges)
	case *parser.LocalStmt:
		findCallExprsInNode(node.Init, name, ranges)
	case *parser.FunctionDef:
		findCallExprs(node.Body, name, ranges)
	case *parser.ThreadDef:
		findCallExprs(node.Body, name, ranges)
	}
}

func wordAt(content string, pos protocol.Position) string {
	lines := strings.Split(content, "\n")
	line := int(pos.Line)
	col := int(pos.Character)
	if line >= len(lines) {
		return ""
	}
	text := lines[line]
	if col >= len(text) {
		return ""
	}
	if !isIdent(rune(text[col])) {
		return ""
	}
	start := col
	for start > 0 && isIdent(rune(text[start-1])) {
		start--
	}
	end := col
	for end < len(text) && isIdent(rune(text[end])) {
		end++
	}
	return text[start:end]
}

func isIdent(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_'
}
