package completion

import (
	"fmt"
	"strings"

	"go.lsp.dev/protocol"

	"github.com/ts91/urscript-ls/internal/builtins"
	"github.com/ts91/urscript-ls/internal/symbols"
)

var keywords = []string{
	"def", "end", "if", "elif", "else", "while", "return",
	"global", "local", "True", "False", "and", "or", "not", "thread",
}

// Provide returns completion items for the given document context.
func Provide(idx *symbols.Index, registry *builtins.Registry) protocol.CompletionItemSlice {
	seen := make(map[string]bool)
	var items protocol.CompletionItemSlice

	// Built-in functions.
	for _, b := range registry.All() {
		if seen[b.Name] {
			continue
		}
		seen[b.Name] = true
		items = append(items, protocol.CompletionItem{
			Label:  b.Name,
			Kind:   protocol.CompletionItemKindFunction,
			Detail: protocol.NewOptional[string](b.Signature),
			Documentation: &protocol.MarkupContent{
				Kind:  protocol.MarkupKindMarkdown,
				Value: b.Doc,
			},
		})
	}

	// User-defined functions.
	if idx != nil {
		for _, fn := range idx.Functions() {
			if seen[fn.Name] {
				continue
			}
			seen[fn.Name] = true
			sig := fmt.Sprintf("%s(%s)", fn.Name, strings.Join(fn.Params, ", "))
			item := protocol.CompletionItem{
				Label:  fn.Name,
				Kind:   protocol.CompletionItemKindFunction,
				Detail: protocol.NewOptional[string](sig),
			}
			if fn.Doc != "" {
				item.Documentation = &protocol.MarkupContent{
					Kind:  protocol.MarkupKindMarkdown,
					Value: fn.Doc,
				}
			}
			items = append(items, item)
		}

		// User-defined variables.
		for _, v := range idx.Variables() {
			if seen[v.Name] {
				continue
			}
			seen[v.Name] = true
			kind := protocol.CompletionItemKindVariable
			detail := "variable"
			switch v.Kind {
			case symbols.KindGlobal:
				detail = "global"
			case symbols.KindLocal:
				detail = "local"
			}
			items = append(items, protocol.CompletionItem{
				Label:  v.Name,
				Kind:   kind,
				Detail: protocol.NewOptional[string](detail),
			})
		}

		// Threads.
		for _, th := range idx.Threads() {
			if seen[th.Name] {
				continue
			}
			seen[th.Name] = true
			items = append(items, protocol.CompletionItem{
				Label:  th.Name,
				Kind:   protocol.CompletionItemKindFunction,
				Detail: protocol.NewOptional[string]("thread " + th.Name + "()"),
			})
		}
	}

	// Keywords.
	for _, kw := range keywords {
		items = append(items, protocol.CompletionItem{
			Label: kw,
			Kind:  protocol.CompletionItemKindKeyword,
		})
	}

	return items
}
