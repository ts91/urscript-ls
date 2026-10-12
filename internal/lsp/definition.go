package lsp

import (
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"

	"github.com/ts91/urscript-ls/internal/symbols"
)

// Definition returns the definition location for the symbol at the given position.
func Definition(content string, pos protocol.Position, docURI uri.URI, idx *symbols.Index) protocol.LocationSlice {
	if idx == nil {
		return nil
	}

	word := wordAt(content, pos)
	if word == "" {
		return nil
	}

	syms := idx.Lookup(word)
	if len(syms) == 0 {
		return nil
	}

	// Return the first definition-level symbol (function, thread, global, local)
	// rather than usage-level (variable assignments, params).
	best := syms[0]
	for _, s := range syms {
		if s.Kind == symbols.KindFunction || s.Kind == symbols.KindThread ||
			s.Kind == symbols.KindGlobal || s.Kind == symbols.KindLocal {
			best = s
			break
		}
	}

	return protocol.LocationSlice{
		{
			URI: docURI,
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      uint32(best.Pos.Line),
					Character: uint32(best.Pos.Column),
				},
				End: protocol.Position{
					Line:      uint32(best.Pos.Line),
					Character: uint32(best.Pos.Column + len(best.Name)),
				},
			},
		},
	}
}

