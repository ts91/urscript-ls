package lsp

import (
	"strings"

	"go.lsp.dev/protocol"

	"github.com/ts91/urscript-ls/internal/parser"
)

// FoldingRanges returns folding ranges for block structures and comment runs.
func FoldingRanges(prog *parser.Program, source string) []protocol.FoldingRange {
	var ranges []protocol.FoldingRange
	foldingWalkStmts(prog.Stmts, &ranges)
	ranges = append(ranges, commentRanges(source)...)
	return ranges
}

func foldingWalkStmts(stmts []parser.Node, ranges *[]protocol.FoldingRange) {
	for _, stmt := range stmts {
		foldingWalkNode(stmt, ranges)
	}
}

func foldingWalkNode(n parser.Node, ranges *[]protocol.FoldingRange) {
	switch node := n.(type) {
	case *parser.FunctionDef:
		addBlockFold(node.Pos, node.End, ranges)
		foldingWalkStmts(node.Body, ranges)

	case *parser.ThreadDef:
		addBlockFold(node.Pos, node.End, ranges)
		foldingWalkStmts(node.Body, ranges)

	case *parser.IfStmt:
		addBlockFold(node.Pos, node.End, ranges)
		foldingWalkStmts(node.Body, ranges)
		for _, elif := range node.ElseIfs {
			foldingWalkStmts(elif.Body, ranges)
		}
		foldingWalkStmts(node.Else, ranges)

	case *parser.WhileStmt:
		addBlockFold(node.Pos, node.End, ranges)
		foldingWalkStmts(node.Body, ranges)
	}
}

func addBlockFold(start, end parser.Pos, ranges *[]protocol.FoldingRange) {
	if end.Line > start.Line {
		*ranges = append(*ranges, protocol.FoldingRange{
			StartLine: uint32(start.Line),
			EndLine:   uint32(end.Line),
			Kind:      protocol.FoldingRangeKindRegion,
		})
	}
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
