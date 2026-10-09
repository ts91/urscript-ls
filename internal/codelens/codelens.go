package codelens

import (
	"fmt"

	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"

	"github.com/ts91/urscript-ls/internal/references"
	"github.com/ts91/urscript-ls/internal/symbols"
)

// Provide returns code lenses showing reference counts for functions and threads.
func Provide(docURI uri.URI, idx *symbols.Index, allDocs map[uri.URI]string) []protocol.CodeLens {
	if idx == nil {
		return nil
	}

	var lenses []protocol.CodeLens
	for _, sym := range idx.Symbols {
		if sym.Kind != symbols.KindFunction && sym.Kind != symbols.KindThread {
			continue
		}

		count := 0
		for _, content := range allDocs {
			positions := references.FindAllInSource(content, sym.Name)
			count += len(positions)
		}
		// Subtract the definition itself.
		count--
		if count < 0 {
			count = 0
		}

		title := fmt.Sprintf("%d references", count)
		if count == 1 {
			title = "1 reference"
		}

		lenses = append(lenses, protocol.CodeLens{
			Range: protocol.Range{
				Start: protocol.Position{Line: uint32(sym.Pos.Line), Character: 0},
				End:   protocol.Position{Line: uint32(sym.Pos.Line), Character: 0},
			},
			Command: protocol.Command{
				Title: title,
			},
		})
	}

	return lenses
}
