package lsp

import (
	"fmt"
	"strings"

	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"

	"github.com/ts91/urscript-ls/internal/builtins"
	"github.com/ts91/urscript-ls/internal/symbols"
)

// CodeActions returns code actions for the given diagnostics.
func CodeActions(
	params *protocol.CodeActionParams,
	registry *builtins.Registry,
	idx *symbols.Index,
) []protocol.CommandOrCodeAction {
	var actions []protocol.CommandOrCodeAction

	for _, diag := range params.Context.Diagnostics {
		msg := diagMessage(diag)

		if strings.Contains(msg, "is not defined") {
			acts := suggestSimilar(params.TextDocument.URI, diag, msg, registry, idx)
			actions = append(actions, acts...)
		}
	}

	return actions
}

func suggestSimilar(
	docURI uri.URI,
	diag protocol.Diagnostic,
	msg string,
	registry *builtins.Registry,
	idx *symbols.Index,
) []protocol.CommandOrCodeAction {
	// Extract the undefined name from "'foo' is not defined".
	name := extractQuoted(msg)
	if name == "" {
		return nil
	}

	// Collect all known names.
	var candidates []string
	for _, b := range registry.All() {
		candidates = append(candidates, b.Name)
	}
	if idx != nil {
		for _, s := range idx.Symbols {
			if s.Kind == symbols.KindFunction || s.Kind == symbols.KindThread {
				candidates = append(candidates, s.Name)
			}
		}
	}

	// Find close matches.
	var actions []protocol.CommandOrCodeAction
	for _, candidate := range candidates {
		if levenshtein(name, candidate) <= 2 {
			kind := protocol.CodeActionKindQuickFix
			preferred := true
			actions = append(actions, &protocol.CodeAction{
				Title:       fmt.Sprintf("Did you mean '%s'?", candidate),
				Kind:        &kind,
				Diagnostics: []protocol.Diagnostic{diag},
				IsPreferred: &preferred,
				Edit: &protocol.WorkspaceEdit{
					Changes: map[uri.URI][]protocol.TextEdit{
						docURI: {{
							Range:   diag.Range,
							NewText: candidate,
						}},
					},
				},
			})
		}
	}
	return actions
}

func extractQuoted(msg string) string {
	start := strings.Index(msg, "'")
	if start < 0 {
		return ""
	}
	end := strings.Index(msg[start+1:], "'")
	if end < 0 {
		return ""
	}
	return msg[start+1 : start+1+end]
}

func diagMessage(d protocol.Diagnostic) string {
	if s, ok := d.Message.(protocol.String); ok {
		return string(s)
	}
	return ""
}

func levenshtein(a, b string) int {
	la, lb := len(a), len(b)
	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}

	prev := make([]int, lb+1)
	curr := make([]int, lb+1)
	for j := range prev {
		prev[j] = j
	}

	for i := 1; i <= la; i++ {
		curr[0] = i
		for j := 1; j <= lb; j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			curr[j] = min3(curr[j-1]+1, prev[j]+1, prev[j-1]+cost)
		}
		prev, curr = curr, prev
	}
	return prev[lb]
}

func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}
