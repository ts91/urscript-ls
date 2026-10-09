package docsymbols

import (
	"fmt"
	"strings"

	"go.lsp.dev/protocol"

	"github.com/ts91/urscript-ls/internal/parser"
	"github.com/ts91/urscript-ls/internal/symbols"
)

// Provide returns the document symbol outline for a file.
func Provide(idx *symbols.Index) protocol.DocumentSymbolSlice {
	if idx == nil {
		return nil
	}

	var result protocol.DocumentSymbolSlice

	for _, sym := range idx.Symbols {
		ds := toDocumentSymbol(sym)
		if ds != nil {
			result = append(result, *ds)
		}
	}

	return result
}

func toDocumentSymbol(sym symbols.Symbol) *protocol.DocumentSymbol {
	// Skip params — they're children of functions, not top-level symbols.
	if sym.Kind == symbols.KindParam {
		return nil
	}

	kind := mapKind(sym)
	selRange := nameRange(sym.Pos, sym.Name)

	fullRange := selRange
	if sym.EndPos != sym.Pos {
		fullRange = protocol.Range{
			Start: protocol.Position{
				Line:      uint32(sym.Pos.Line),
				Character: uint32(sym.Pos.Column),
			},
			End: protocol.Position{
				Line:      uint32(sym.EndPos.Line),
				Character: uint32(sym.EndPos.Column),
			},
		}
	}

	ds := &protocol.DocumentSymbol{
		Name:           sym.Name,
		Kind:           kind,
		Range:          fullRange,
		SelectionRange: selRange,
	}

	if sym.Kind == symbols.KindFunction {
		detail := fmt.Sprintf("(%s)", strings.Join(sym.Params, ", "))
		ds.Detail = &detail
	}
	if sym.Kind == symbols.KindThread {
		detail := "thread"
		ds.Detail = &detail
	}

	return ds
}

func mapKind(sym symbols.Symbol) protocol.SymbolKind {
	switch sym.Kind {
	case symbols.KindFunction:
		return protocol.SymbolKindFunction
	case symbols.KindThread:
		return protocol.SymbolKindFunction
	case symbols.KindVariable:
		return protocol.SymbolKindVariable
	case symbols.KindGlobal:
		return protocol.SymbolKindVariable
	case symbols.KindLocal:
		return protocol.SymbolKindVariable
	default:
		return protocol.SymbolKindVariable
	}
}

func nameRange(pos parser.Pos, name string) protocol.Range {
	p := protocol.Position{
		Line:      uint32(pos.Line),
		Character: uint32(pos.Column),
	}
	end := protocol.Position{
		Line:      uint32(pos.Line),
		Character: uint32(pos.Column + len(name)),
	}
	return protocol.Range{Start: p, End: end}
}
