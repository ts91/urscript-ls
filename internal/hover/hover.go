package hover

import (
	"fmt"
	"strings"
	"unicode"

	"go.lsp.dev/protocol"

	"github.com/ts91/urscript-ls/internal/builtins"
	"github.com/ts91/urscript-ls/internal/symbols"
)

// Provide returns hover information for the word at the given position.
func Provide(content string, pos protocol.Position, registry *builtins.Registry, idx *symbols.Index) *protocol.Hover {
	word, wordRange := wordAtPosition(content, pos)
	if word == "" {
		return nil
	}

	// Try builtins first.
	if b, ok := registry.Lookup(word); ok {
		md := fmt.Sprintf("```urscript\n%s\n```\n\n%s", b.Signature, formatDoc(b.Doc))
		return &protocol.Hover{
			Contents: &protocol.MarkupContent{
				Kind:  protocol.MarkupKindMarkdown,
				Value: md,
			},
			Range: &wordRange,
		}
	}

	// Fall back to user-defined symbols.
	if idx == nil {
		return nil
	}
	syms := idx.Lookup(word)
	if len(syms) == 0 {
		return nil
	}

	md := formatUserSymbol(syms[0])
	return &protocol.Hover{
		Contents: &protocol.MarkupContent{
			Kind:  protocol.MarkupKindMarkdown,
			Value: md,
		},
		Range: &wordRange,
	}
}

func formatUserSymbol(sym symbols.Symbol) string {
	var parts []string

	switch sym.Kind {
	case symbols.KindFunction:
		sig := fmt.Sprintf("def %s(%s)", sym.Name, strings.Join(sym.Params, ", "))
		parts = append(parts, fmt.Sprintf("```urscript\n%s\n```", sig))
	case symbols.KindThread:
		parts = append(parts, fmt.Sprintf("```urscript\nthread %s()\n```", sym.Name))
	case symbols.KindGlobal:
		parts = append(parts, fmt.Sprintf("```urscript\nglobal %s\n```", sym.Name))
	case symbols.KindLocal:
		parts = append(parts, fmt.Sprintf("```urscript\nlocal %s\n```", sym.Name))
	default:
		parts = append(parts, fmt.Sprintf("```urscript\n%s\n```", sym.Name))
	}

	if sym.Doc != "" {
		parts = append(parts, formatDoc(sym.Doc))
	}

	return strings.Join(parts, "\n\n")
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

// formatDoc converts raw docstring content into readable Markdown.
func formatDoc(doc string) string {
	lines := strings.Split(doc, "\n")
	var out []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "@param ") {
			parts := strings.Fields(trimmed)
			// @param name type description...
			if len(parts) >= 3 {
				desc := strings.Join(parts[3:], " ")
				out = append(out, fmt.Sprintf("**@param** `%s` *%s* %s", parts[1], parts[2], desc))
			} else if len(parts) == 2 {
				out = append(out, fmt.Sprintf("**@param** `%s`", parts[1]))
			}
		} else if strings.HasPrefix(trimmed, "@returns ") {
			parts := strings.Fields(trimmed)
			// @returns type description...
			if len(parts) >= 2 {
				desc := strings.Join(parts[2:], " ")
				out = append(out, fmt.Sprintf("**@returns** *%s* %s", parts[1], desc))
			}
		} else if strings.HasPrefix(trimmed, "#") {
			out = append(out, `\`+trimmed)
		} else {
			out = append(out, trimmed)
		}
	}
	return strings.Join(out, "  \n")
}
