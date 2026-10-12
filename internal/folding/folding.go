package folding

import (
	"strings"

	"go.lsp.dev/protocol"

	"github.com/ts91/urscript-ls/internal/parser"
)

// Provide returns folding ranges for block structures and comment runs.
func Provide(prog *parser.Program, source string) []protocol.FoldingRange {
	var ranges []protocol.FoldingRange
	walkStmts(prog.Stmts, &ranges)
	ranges = append(ranges, commentRanges(source)...)
	return ranges
}

func walkStmts(stmts []parser.Node, ranges *[]protocol.FoldingRange) {
	for _, stmt := range stmts {
		walkNode(stmt, ranges)
	}
}

func walkNode(n parser.Node, ranges *[]protocol.FoldingRange) {
	switch node := n.(type) {
	case *parser.FunctionDef:
		addBlockRange(node.Pos, node.Body, ranges)
		walkStmts(node.Body, ranges)

	case *parser.ThreadDef:
		addBlockRange(node.Pos, node.Body, ranges)
		walkStmts(node.Body, ranges)

	case *parser.IfStmt:
		// Fold the entire if/elif/else/end as one range.
		endLine := findEndLine(node.Pos, node.Body, node.ElseIfs, node.Else)
		if endLine > node.Pos.Line {
			*ranges = append(*ranges, protocol.FoldingRange{
				StartLine: uint32(node.Pos.Line),
				EndLine:   uint32(endLine),
				Kind:      protocol.FoldingRangeKindRegion,
			})
		}
		walkStmts(node.Body, ranges)
		for _, elif := range node.ElseIfs {
			walkStmts(elif.Body, ranges)
		}
		walkStmts(node.Else, ranges)

	case *parser.WhileStmt:
		addBlockRange(node.Pos, node.Body, ranges)
		walkStmts(node.Body, ranges)

	}
}

func addBlockRange(pos parser.Pos, body []parser.Node, ranges *[]protocol.FoldingRange) {
	if len(body) == 0 {
		return
	}
	last := parser.NodePos(body[len(body)-1])
	// The "end" keyword is on the line after the last statement (approximately).
	// We fold from the opening line to the end line.
	endLine := last.Line + 1
	if endLine > pos.Line {
		*ranges = append(*ranges, protocol.FoldingRange{
			StartLine: uint32(pos.Line),
			EndLine:   uint32(endLine),
			Kind:      protocol.FoldingRangeKindRegion,
		})
	}
}

func findEndLine(pos parser.Pos, body []parser.Node, elseIfs []parser.ElseIfClause, elseBody []parser.Node) int {
	endLine := pos.Line
	if len(elseBody) > 0 {
		last := parser.NodePos(elseBody[len(elseBody)-1])
		endLine = last.Line + 1
	} else if len(elseIfs) > 0 {
		lastElif := elseIfs[len(elseIfs)-1]
		if len(lastElif.Body) > 0 {
			last := parser.NodePos(lastElif.Body[len(lastElif.Body)-1])
			endLine = last.Line + 1
		}
	} else if len(body) > 0 {
		last := parser.NodePos(body[len(body)-1])
		endLine = last.Line + 1
	}
	return endLine
}

func commentRanges(source string) []protocol.FoldingRange {
	lines := strings.Split(source, "\n")
	var ranges []protocol.FoldingRange
	start := -1

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			if start < 0 {
				start = i
			}
		} else {
			if start >= 0 && i-1 > start {
				ranges = append(ranges, protocol.FoldingRange{
					StartLine: uint32(start),
					EndLine:   uint32(i - 1),
					Kind:      protocol.FoldingRangeKindComment,
				})
			}
			start = -1
		}
	}
	// Trailing comment block.
	if start >= 0 && len(lines)-1 > start {
		ranges = append(ranges, protocol.FoldingRange{
			StartLine: uint32(start),
			EndLine:   uint32(len(lines) - 1),
			Kind:      protocol.FoldingRangeKindComment,
		})
	}

	return ranges
}
