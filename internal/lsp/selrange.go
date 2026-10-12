package lsp

import (
	"strings"

	"go.lsp.dev/protocol"

	"github.com/ts91/urscript-ls/internal/parser"
)

// SelectionRanges returns selection ranges for the given cursor positions.
func SelectionRanges(prog *parser.Program, source string, positions []protocol.Position) []protocol.SelectionRange {
	lines := strings.Split(source, "\n")
	results := make([]protocol.SelectionRange, len(positions))

	for i, pos := range positions {
		results[i] = buildSelectionRange(prog, lines, int(pos.Line), int(pos.Character))
	}
	return results
}

func buildSelectionRange(prog *parser.Program, lines []string, line, col int) protocol.SelectionRange {
	// Start with the word under cursor as the innermost range.
	wordRange := wordRangeAt(lines, line, col)

	// Collect containing AST nodes from innermost to outermost.
	chain := findContaining(prog.Stmts, line)

	// Build the SelectionRange chain from outermost to innermost.
	var result *protocol.SelectionRange

	// Outermost: entire document.
	lastLine := len(lines) - 1
	lastChar := 0
	if lastLine >= 0 {
		lastChar = len(lines[lastLine])
	}
	docRange := protocol.Range{
		Start: protocol.Position{Line: 0, Character: 0},
		End:   protocol.Position{Line: uint32(lastLine), Character: uint32(lastChar)},
	}
	result = &protocol.SelectionRange{Range: docRange}

	// Add AST node ranges from outermost to innermost.
	for _, node := range chain {
		nodeRange := rangeForNode(node, lines)
		sr := &protocol.SelectionRange{Range: nodeRange, Parent: result}
		result = sr
	}

	// Innermost: the word under cursor.
	if wordRange != nil {
		sr := &protocol.SelectionRange{Range: *wordRange, Parent: result}
		result = sr
	}

	return *result
}

// findContaining returns AST nodes that contain the given line,
// ordered from outermost to innermost.
func findContaining(stmts []parser.Node, line int) []parser.Node {
	var chain []parser.Node
	findContainingRec(stmts, line, &chain)
	return chain
}

func findContainingRec(stmts []parser.Node, line int, chain *[]parser.Node) {
	for _, stmt := range stmts {
		if !nodeContainsLine(stmt, line) {
			continue
		}
		*chain = append(*chain, stmt)

		switch node := stmt.(type) {
		case *parser.FunctionDef:
			findContainingRec(node.Body, line, chain)
		case *parser.ThreadDef:
			findContainingRec(node.Body, line, chain)
		case *parser.IfStmt:
			findContainingRec(node.Body, line, chain)
			for _, elif := range node.ElseIfs {
				findContainingRec(elif.Body, line, chain)
			}
			findContainingRec(node.Else, line, chain)
		case *parser.WhileStmt:
			findContainingRec(node.Body, line, chain)
		}
		return
	}
}

func nodeContainsLine(n parser.Node, line int) bool {
	pos := parser.NodePos(n)
	startLine := pos.Line

	switch node := n.(type) {
	case *parser.FunctionDef:
		return line >= startLine && line <= endLine(node.Body, startLine)
	case *parser.ThreadDef:
		return line >= startLine && line <= endLine(node.Body, startLine)
	case *parser.IfStmt:
		end := endLine(node.Body, startLine)
		for _, elif := range node.ElseIfs {
			e := endLine(elif.Body, elif.Pos.Line)
			if e > end {
				end = e
			}
		}
		e := endLine(node.Else, end)
		if e > end {
			end = e
		}
		// +1 for the "end" keyword line
		return line >= startLine && line <= end+1
	case *parser.WhileStmt:
		return line >= startLine && line <= endLine(node.Body, startLine)+1
	default:
		return line == startLine
	}
}

func endLine(body []parser.Node, fallback int) int {
	if len(body) == 0 {
		return fallback
	}
	return parser.NodePos(body[len(body)-1]).Line
}

func rangeForNode(n parser.Node, lines []string) protocol.Range {
	pos := parser.NodePos(n)
	startLine := pos.Line
	startChar := pos.Column

	var end int
	switch node := n.(type) {
	case *parser.FunctionDef:
		end = endLine(node.Body, startLine) + 1
	case *parser.ThreadDef:
		end = endLine(node.Body, startLine) + 1
	case *parser.IfStmt:
		end = endLine(node.Body, startLine)
		for _, elif := range node.ElseIfs {
			e := endLine(elif.Body, elif.Pos.Line)
			if e > end {
				end = e
			}
		}
		e := endLine(node.Else, end)
		if e > end {
			end = e
		}
		end++
	case *parser.WhileStmt:
		end = endLine(node.Body, startLine) + 1
	default:
		end = startLine
	}

	endChar := 0
	if end < len(lines) {
		endChar = len(lines[end])
	}

	return protocol.Range{
		Start: protocol.Position{Line: uint32(startLine), Character: uint32(startChar)},
		End:   protocol.Position{Line: uint32(end), Character: uint32(endChar)},
	}
}

func wordRangeAt(lines []string, line, col int) *protocol.Range {
	if line >= len(lines) {
		return nil
	}
	text := lines[line]
	if col >= len(text) {
		return nil
	}
	if !isIdent(rune(text[col])) {
		return nil
	}
	start := col
	for start > 0 && isIdent(rune(text[start-1])) {
		start--
	}
	end := col
	for end < len(text) && isIdent(rune(text[end])) {
		end++
	}
	r := protocol.Range{
		Start: protocol.Position{Line: uint32(line), Character: uint32(start)},
		End:   protocol.Position{Line: uint32(line), Character: uint32(end)},
	}
	return &r
}

