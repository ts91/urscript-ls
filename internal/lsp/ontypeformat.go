package lsp

import (
	"strings"

	"go.lsp.dev/protocol"
)

// OnTypeFormat returns text edits for on-type formatting.
// Triggered on newline to auto-indent after block openers.
func OnTypeFormat(content string, pos protocol.Position, ch string, opts protocol.FormattingOptions) []protocol.TextEdit {
	if ch != "\n" {
		return nil
	}

	line := int(pos.Line)
	lines := strings.Split(content, "\n")

	// The newline was just typed, so we look at the previous line.
	prevLine := line - 1
	if prevLine < 0 || prevLine >= len(lines) {
		return nil
	}

	trimmed := strings.TrimSpace(lines[prevLine])

	indent := "  "
	if !opts.InsertSpaces {
		indent = "\t"
	} else if opts.TabSize > 0 {
		indent = strings.Repeat(" ", int(opts.TabSize))
	}

	// Count existing indentation of the previous line.
	prevIndent := leadingWhitespace(lines[prevLine])

	if onTypeIsOpener(trimmed) {
		// After a block opener, indent one level deeper.
		newIndent := prevIndent + indent
		return []protocol.TextEdit{{
			Range: protocol.Range{
				Start: protocol.Position{Line: uint32(line), Character: 0},
				End:   protocol.Position{Line: uint32(line), Character: uint32(currentIndentLen(lines, line))},
			},
			NewText: newIndent,
		}}
	}

	return nil
}

func onTypeIsOpener(line string) bool {
	return fmtIsOpener(line)
}

func leadingWhitespace(line string) string {
	i := 0
	for i < len(line) && (line[i] == ' ' || line[i] == '\t') {
		i++
	}
	return line[:i]
}

func currentIndentLen(lines []string, line int) int {
	if line >= len(lines) {
		return 0
	}
	return len(leadingWhitespace(lines[line]))
}
