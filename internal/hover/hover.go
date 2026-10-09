package hover

import (
	"fmt"
	"strings"
	"unicode"

	"go.lsp.dev/protocol"

	"github.com/ts91/urscript-ls/internal/builtins"
)

// Provide returns hover information for the word at the given position.
func Provide(content string, pos protocol.Position, registry *builtins.Registry) *protocol.Hover {
	word, wordRange := wordAtPosition(content, pos)
	if word == "" {
		return nil
	}

	b, ok := registry.Lookup(word)
	if !ok {
		return nil
	}

	md := fmt.Sprintf("```urscript\n%s\n```\n\n%s", b.Signature, b.Doc)
	return &protocol.Hover{
		Contents: &protocol.MarkupContent{
			Kind:  protocol.MarkupKindMarkdown,
			Value: md,
		},
		Range: &wordRange,
	}
}

// wordAtPosition extracts the identifier-like word at the given position.
func wordAtPosition(content string, pos protocol.Position) (string, protocol.Range) {
	lines := strings.Split(content, "\n")
	line := int(pos.Line)
	col := int(pos.Character)

	if line >= len(lines) {
		return "", protocol.Range{}
	}
	text := lines[line]
	if col >= len(text) {
		return "", protocol.Range{}
	}

	if !isIdentChar(rune(text[col])) {
		return "", protocol.Range{}
	}

	start := col
	for start > 0 && isIdentChar(rune(text[start-1])) {
		start--
	}
	end := col
	for end < len(text) && isIdentChar(rune(text[end])) {
		end++
	}

	word := text[start:end]
	r := protocol.Range{
		Start: protocol.Position{Line: uint32(line), Character: uint32(start)},
		End:   protocol.Position{Line: uint32(line), Character: uint32(end)},
	}
	return word, r
}

func isIdentChar(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
