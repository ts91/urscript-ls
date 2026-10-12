package lsp

import (
	"go.lsp.dev/protocol"
)

// Highlights returns document highlights for all occurrences of the word
// at the given position in the current file.
func Highlights(content string, pos protocol.Position, word string) []protocol.DocumentHighlight {
	if word == "" {
		return nil
	}

	positions := FindAllInSource(content, word)
	if len(positions) == 0 {
		return nil
	}

	highlights := make([]protocol.DocumentHighlight, len(positions))
	for i, p := range positions {
		highlights[i] = protocol.DocumentHighlight{
			Kind: protocol.DocumentHighlightKindText,
			Range: protocol.Range{
				Start: protocol.Position{Line: uint32(p.Line), Character: uint32(p.Column)},
				End:   protocol.Position{Line: uint32(p.Line), Character: uint32(p.Column + p.Len)},
			},
		}
	}
	return highlights
}
