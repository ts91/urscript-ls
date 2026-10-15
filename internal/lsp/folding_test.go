package lsp

import (
	"testing"

	"go.lsp.dev/protocol"

	"github.com/ts91/urscript-ls/internal/parser"
)

func regionFolds(source string) []protocol.FoldingRange {
	prog, _ := parser.Parse(source)
	var regions []protocol.FoldingRange
	for _, r := range FoldingRanges(prog, source) {
		if r.Kind == protocol.FoldingRangeKindRegion {
			regions = append(regions, r)
		}
	}
	return regions
}

func hasRegionFold(regions []protocol.FoldingRange, start, end uint32) bool {
	for _, r := range regions {
		if r.StartLine == start && r.EndLine == end {
			return true
		}
	}
	return false
}

func TestFoldingIfWithCallsAndOr(t *testing.T) {
	source := "if get_tool_analog_in(0) < 2.9 or get_tool_analog_in(0) > 3.6:\n  fold_me = True\nend\n"
	regions := regionFolds(source)
	if !hasRegionFold(regions, 0, 2) {
		t.Fatalf("expected if fold 0-2, got %v", regions)
	}
}

func TestFoldingNestedIfCommentOnlyBranches(t *testing.T) {
	source := "" +
		"if (robot_type == ROBOT_TYPE_UR10):\n" +
		"  if (joint == 1):\n" +
		"    # Shoulder acceleration = reduced_size4_acceleration\n" +
		"  elif (joint == 2):\n" +
		"    # Elbow acceleration = reduced_size3_acceleration\n" +
		"  end\n" +
		"end\n"
	regions := regionFolds(source)
	if !hasRegionFold(regions, 0, 6) {
		t.Fatalf("expected outer if fold 0-6, got %v", regions)
	}
	if !hasRegionFold(regions, 1, 5) {
		t.Fatalf("expected inner if fold 1-5, got %v", regions)
	}
}

func TestFoldingThreadWhileTrue(t *testing.T) {
	source := "" +
		"thread Thread_1():\n" +
		"  while (True):\n" +
		"    sleep(0.1)\n" +
		"    myur_prod_test_time = floor(timer_sec/60)\n" +
		"  end\n" +
		"end\n"
	regions := regionFolds(source)
	if !hasRegionFold(regions, 0, 5) {
		t.Fatalf("expected thread fold 0-5, got %v", regions)
	}
	if !hasRegionFold(regions, 1, 4) {
		t.Fatalf("expected while fold 1-4, got %v", regions)
	}
}

func TestFoldingNestedIfNoSpaceAfterIf(t *testing.T) {
	source := "" +
		"if (server.get_beginning() == 1):\n" +
		"  if(bypass == 0):\n" +
		"    do_check = True\n" +
		"  end\n" +
		"end\n"
	regions := regionFolds(source)
	if !hasRegionFold(regions, 0, 4) {
		t.Fatalf("expected outer if fold 0-4, got %v", regions)
	}
	if !hasRegionFold(regions, 1, 3) {
		t.Fatalf("expected inner if fold 1-3, got %v", regions)
	}
}
