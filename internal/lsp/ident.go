package lsp

import (
	"strings"
	"unicode"

	"go.lsp.dev/protocol"
)

func isIdent(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
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
