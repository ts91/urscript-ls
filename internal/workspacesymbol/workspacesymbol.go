package workspacesymbol

import (
	"strings"

	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"

	"github.com/ts91/urscript-ls/internal/symbols"
)

// Provide returns workspace symbols matching the query across all indexed files.
func Provide(query string, indices map[uri.URI]*symbols.Index) protocol.WorkspaceSymbolSlice {
	query = strings.ToLower(query)
	var result protocol.WorkspaceSymbolSlice

	for docURI, idx := range indices {
		if idx == nil {
			continue
		}
		for _, sym := range idx.Symbols {
			if sym.Kind == symbols.KindParam {
				continue
			}
			if query != "" && !strings.Contains(strings.ToLower(sym.Name), query) {
				continue
			}
			result = append(result, protocol.WorkspaceSymbol{
				BaseSymbolInformation: protocol.BaseSymbolInformation{
					Name: sym.Name,
					Kind: mapKind(sym.Kind),
				},
				Location: &protocol.Location{
					URI: docURI,
					Range: protocol.Range{
						Start: protocol.Position{Line: uint32(sym.Pos.Line), Character: uint32(sym.Pos.Column)},
						End:   protocol.Position{Line: uint32(sym.Pos.Line), Character: uint32(sym.Pos.Column + len(sym.Name))},
					},
				},
			})
		}
	}
	return result
}

func mapKind(k symbols.SymbolKind) protocol.SymbolKind {
	switch k {
	case symbols.KindFunction, symbols.KindThread:
		return protocol.SymbolKindFunction
	case symbols.KindGlobal, symbols.KindLocal, symbols.KindVariable:
		return protocol.SymbolKindVariable
	default:
		return protocol.SymbolKindVariable
	}
}
