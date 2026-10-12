package lsp

import (
	"fmt"
	"strings"

	"go.lsp.dev/protocol"

	"github.com/ts91/urscript-ls/internal/builtins"
	"github.com/ts91/urscript-ls/internal/symbols"
)

// SignatureHelp returns signature help for the function call at the given position.
func SignatureHelp(content string, pos protocol.Position, registry *builtins.Registry, idx *symbols.Index) *protocol.SignatureHelp {
	line := int(pos.Line)
	col := int(pos.Character)

	lines := strings.Split(content, "\n")
	if line >= len(lines) {
		return nil
	}
	text := lines[line]
	if col > len(text) {
		col = len(text)
	}

	// Walk backwards from cursor to find the function name and count commas
	// to determine the active parameter.
	funcName, activeParam := findCallContext(text[:col])
	if funcName == "" {
		return nil
	}

	// Look up in builtins first, then user symbols.
	var sig string
	var params []string
	var doc string

	if b, ok := registry.Lookup(funcName); ok {
		sig = b.Signature
		params = extractParams(b.Signature)
		doc = b.Doc
	} else if idx != nil {
		for _, s := range idx.Lookup(funcName) {
			if s.Kind == symbols.KindFunction {
				params = s.Params
				sig = fmt.Sprintf("%s(%s)", s.Name, strings.Join(s.Params, ", "))
				doc = s.Doc
				break
			}
		}
	}

	if sig == "" {
		return nil
	}

	paramInfos := make([]protocol.ParameterInformation, len(params))
	for i, p := range params {
		paramInfos[i] = protocol.ParameterInformation{
			Label: protocol.String(p),
		}
	}

	activeSig := uint32(0)

	sigInfo := protocol.SignatureInformation{
		Label:      sig,
		Parameters: paramInfos,
	}
	if doc != "" {
		sigInfo.Documentation = &protocol.MarkupContent{
			Kind:  protocol.MarkupKindMarkdown,
			Value: doc,
		}
	}

	return &protocol.SignatureHelp{
		Signatures:      []protocol.SignatureInformation{sigInfo},
		ActiveSignature: &activeSig,
		ActiveParameter: protocol.NewNullable(uint32(activeParam)),
	}
}

// findCallContext walks backwards from the cursor to find the enclosing
// function call name and the current parameter index (comma count).
func findCallContext(text string) (funcName string, paramIndex int) {
	depth := 0
	commas := 0
	i := len(text) - 1

	for i >= 0 {
		ch := text[i]
		switch ch {
		case ')':
			depth++
		case '(':
			if depth > 0 {
				depth--
			} else {
				// Found the opening paren of the call we're in.
				name := extractIdentBefore(text, i)
				return name, commas
			}
		case ',':
			if depth == 0 {
				commas++
			}
		case '"':
			// Skip backwards over string literal.
			i--
			for i >= 0 && text[i] != '"' {
				if i > 0 && text[i-1] == '\\' {
					i--
				}
				i--
			}
		}
		i--
	}
	return "", 0
}

// extractIdentBefore returns the identifier immediately before position pos.
func extractIdentBefore(text string, pos int) string {
	end := pos
	start := end - 1
	for start >= 0 && isIdentByte(text[start]) {
		start--
	}
	start++
	if start >= end {
		return ""
	}
	return text[start:end]
}

func isIdentByte(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_'
}

