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

	// Compute byte offset of the cursor so we can feed findCallContext
	// all text up to the cursor, enabling multiline call support.
	offset := 0
	for i, l := range strings.Split(content, "\n") {
		if i == line {
			if col > len(l) {
				col = len(l)
			}
			offset += col
			break
		}
		offset += len(l) + 1 // +1 for newline
	}

	funcName, activeParam := findCallContext(content[:offset])
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

	paramDocs := extractParamDocs(doc)
	paramInfos := make([]protocol.ParameterInformation, len(params))
	for i, p := range params {
		pi := protocol.ParameterInformation{
			Label: protocol.String(p),
		}
		if d, ok := paramDocs[p]; ok {
			pi.Documentation = &protocol.MarkupContent{
				Kind:  protocol.MarkupKindMarkdown,
				Value: d,
			}
		}
		paramInfos[i] = pi
	}

	activeSig := uint32(0)

	sigInfo := protocol.SignatureInformation{
		Label:      sig,
		Parameters: paramInfos,
	}
	if doc != "" {
		sigInfo.Documentation = &protocol.MarkupContent{
			Kind:  protocol.MarkupKindMarkdown,
			Value: formatDoc(doc),
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

// extractParamDocs parses @param lines from a doc string and returns a map
// from parameter name to a formatted markdown description.
func extractParamDocs(doc string) map[string]string {
	if doc == "" {
		return nil
	}
	result := make(map[string]string)
	for _, line := range strings.Split(doc, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "@param ") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		name := fields[1]
		typ := fields[2]
		var desc string
		if len(fields) >= 4 {
			desc = strings.Join(fields[3:], " ")
		}
		if desc != "" {
			result[name] = fmt.Sprintf("`%s` — %s", typ, desc)
		} else {
			result[name] = fmt.Sprintf("`%s`", typ)
		}
	}
	return result
}

