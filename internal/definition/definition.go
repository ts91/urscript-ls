package definition

import (
	"strings"
	"unicode"

	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"

	"github.com/ts91/urscript-ls/internal/symbols"
)

// Provide returns the definition location for the symbol at the given position.
func Provide(content string, pos protocol.Position, docURI uri.URI, idx *symbols.Index) protocol.LocationSlice {
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
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
