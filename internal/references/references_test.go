package references

import (
	"testing"

	"go.lsp.dev/uri"
)

func TestFindAllInSource(t *testing.T) {
	src := `def my_func():
  x = 1
end
my_func()
y = my_func()`

	positions := FindAllInSource(src, "my_func")
	if len(positions) != 3 {
		t.Fatalf("expected 3 occurrences of my_func, got %d", len(positions))
	}

	// Line 0: "def my_func():"
	if positions[0].Line != 0 || positions[0].Column != 4 {
		t.Errorf("expected first occurrence at 0:4, got %d:%d", positions[0].Line, positions[0].Column)
	}
	// Line 3: "my_func()"
	if positions[1].Line != 3 || positions[1].Column != 0 {
		t.Errorf("expected second occurrence at 3:0, got %d:%d", positions[1].Line, positions[1].Column)
	}
	// Line 4: "y = my_func()"
	if positions[2].Line != 4 || positions[2].Column != 4 {
		t.Errorf("expected third occurrence at 4:4, got %d:%d", positions[2].Line, positions[2].Column)
	}
}

func TestFindAllInSourceWholeWord(t *testing.T) {
	src := `speed = 1
speed_max = 2
my_speed = speed`

	positions := FindAllInSource(src, "speed")
	// Should match "speed" on line 0 and the last "speed" on line 2, but NOT
	// "speed_max" or "my_speed".
	if len(positions) != 2 {
		t.Fatalf("expected 2 whole-word occurrences, got %d", len(positions))
	}
}

func TestToLocations(t *testing.T) {
	positions := []Pos{{Line: 0, Column: 4, Len: 5}}
	docURI := uri.URI("file:///test.script")
	locs := ToLocations(positions, docURI)
	if len(locs) != 1 {
		t.Fatalf("expected 1 location, got %d", len(locs))
	}
	if locs[0].URI != docURI {
		t.Errorf("expected URI %s, got %s", docURI, locs[0].URI)
	}
	if locs[0].Range.Start.Character != 4 || locs[0].Range.End.Character != 9 {
		t.Errorf("unexpected range: %v", locs[0].Range)
	}
}

func TestToTextEdits(t *testing.T) {
	positions := []Pos{
		{Line: 0, Column: 4, Len: 5},
		{Line: 3, Column: 0, Len: 5},
	}
	edits := ToTextEdits(positions, "renamed")
	if len(edits) != 2 {
		t.Fatalf("expected 2 edits, got %d", len(edits))
	}
	if edits[0].NewText != "renamed" {
		t.Errorf("expected NewText 'renamed', got '%s'", edits[0].NewText)
	}
}
