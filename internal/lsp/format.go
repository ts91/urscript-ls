package lsp

import (
	"strings"

	"go.lsp.dev/protocol"
)

// Format reformats URScript source code and returns a single TextEdit
// replacing the entire document.
func Format(content string, opts protocol.FormattingOptions) []protocol.TextEdit {
	indent := "  "
	if !opts.InsertSpaces {
		indent = "\t"
	} else if opts.TabSize > 0 {
		indent = strings.Repeat(" ", int(opts.TabSize))
	}

	lines := strings.Split(content, "\n")
	formatted := make([]string, 0, len(lines))

	level := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip empty lines — preserve them but without trailing spaces.
		if trimmed == "" {
			formatted = append(formatted, "")
			continue
		}

		// Dedent before the line for closers and mid-block keywords.
		if fmtIsCloser(trimmed) || fmtIsMidBlock(trimmed) {
			level--
			if level < 0 {
				level = 0
			}
		}

		formatted = append(formatted, strings.Repeat(indent, level)+trimmed)

		// Indent after openers and mid-block keywords (elif/else open a
		// new block at the same level as their parent if/elif).
		if fmtIsMidBlock(trimmed) || fmtIsOpener(trimmed) {
			level++
		}
	}

	// Trim trailing empty lines, then ensure exactly one final newline.
	for len(formatted) > 0 && formatted[len(formatted)-1] == "" {
		formatted = formatted[:len(formatted)-1]
	}

	result := strings.Join(formatted, "\n") + "\n"

	if result == content {
		return nil
	}

	// Return a single edit replacing the entire document.
	lastLine := len(lines) - 1
	lastChar := len(lines[lastLine])

	return []protocol.TextEdit{
		{
			Range: protocol.Range{
				Start: protocol.Position{Line: 0, Character: 0},
				End:   protocol.Position{Line: uint32(lastLine), Character: uint32(lastChar)},
			},
			NewText: result,
		},
	}
}

func fmtIsOpener(line string) bool {
	if fmtIsMidBlock(line) {
		return false
	}
	return fmtLineStartsBlock(line, "def") ||
		fmtLineStartsBlock(line, "thread") ||
		fmtLineStartsBlock(line, "if") ||
		fmtLineStartsBlock(line, "while")
}

// fmtLineStartsBlock reports whether line begins with keyword followed by
// space or '(', e.g. "if(" and "while (True):".
func fmtLineStartsBlock(line, keyword string) bool {
	if !strings.HasPrefix(line, keyword) {
		return false
	}
	rest := line[len(keyword):]
	if rest == "" {
		return false
	}
	return rest[0] == ' ' || rest[0] == '('
}

func fmtIsCloser(line string) bool {
	return line == "end"
}

func fmtIsMidBlock(line string) bool {
	if strings.HasPrefix(line, "elif ") || strings.HasPrefix(line, "else:") || line == "else" {
		return true
	}
	return false
}
