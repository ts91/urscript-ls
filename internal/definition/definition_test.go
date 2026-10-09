package definition

import (
	"testing"

	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"

	"github.com/ts91/urscript-ls/internal/parser"
	"github.com/ts91/urscript-ls/internal/symbols"
)

func TestDefinitionFunction(t *testing.T) {
	src := `def my_func():
  return 1
end

x = my_func()`
	prog, errs := parser.Parse(src)
	if len(errs) != 0 {
		t.Fatalf("unexpected parse errors: %v", errs)
	}
	idx := symbols.Build(prog, src)
	docURI := uri.URI("file:///test.script")

	// Cursor on "my_func" in the call on line 4.
	pos := protocol.Position{Line: 4, Character: 5}
	result := Provide(src, pos, docURI, idx)

	if len(result) == 0 {
		t.Fatal("expected at least one location")
	}

	loc := result[0]
	if loc.URI != docURI {
		t.Errorf("expected URI %s, got %s", docURI, loc.URI)
	}
	if loc.Range.Start.Line != 0 {
		t.Errorf("expected definition on line 0, got %d", loc.Range.Start.Line)
	}
}

func TestDefinitionVariable(t *testing.T) {
	src := `speed = 1.5
movej(q, v=speed)`
	prog, errs := parser.Parse(src)
	if len(errs) != 0 {
		t.Fatalf("unexpected parse errors: %v", errs)
	}
	idx := symbols.Build(prog, src)
	docURI := uri.URI("file:///test.script")

	// Cursor on "speed" in the call on line 1.
	pos := protocol.Position{Line: 1, Character: 12}
	result := Provide(src, pos, docURI, idx)

	if len(result) == 0 {
		t.Fatal("expected at least one location")
	}
	if result[0].Range.Start.Line != 0 {
		t.Errorf("expected definition on line 0, got %d", result[0].Range.Start.Line)
	}
}

func TestDefinitionNilIndex(t *testing.T) {
	result := Provide("x = 1", protocol.Position{Line: 0, Character: 0}, "file:///x", nil)
	if result != nil {
		t.Errorf("expected nil for nil index, got %v", result)
	}
}

func TestDefinitionUnknown(t *testing.T) {
	src := `x = unknown_var`
	prog, errs := parser.Parse(src)
	if len(errs) != 0 {
		t.Fatalf("unexpected parse errors: %v", errs)
	}
	idx := symbols.Build(prog, src)
	docURI := uri.URI("file:///test.script")

	// "unknown_var" is on the RHS — it won't be in the index as a definition,
	// only 'x' is indexed as a variable.
	pos := protocol.Position{Line: 0, Character: 5}
	result := Provide(src, pos, docURI, idx)

	if len(result) != 0 {
		t.Errorf("expected no results for unknown symbol, got %v", result)
	}
}

func TestDefinitionGlobal(t *testing.T) {
	src := `global speed = 1.0
def set_speed():
  speed = 2.0
end`
	prog, errs := parser.Parse(src)
	if len(errs) != 0 {
		t.Fatalf("unexpected parse errors: %v", errs)
	}
	idx := symbols.Build(prog, src)
	docURI := uri.URI("file:///test.script")

	// Cursor on "speed" assignment inside function (line 2).
	pos := protocol.Position{Line: 2, Character: 2}
	result := Provide(src, pos, docURI, idx)

	if len(result) == 0 {
		t.Fatal("expected at least one location")
	}
	// Should jump to the global declaration on line 0.
	if result[0].Range.Start.Line != 0 {
		t.Errorf("expected definition on line 0 (global), got %d", result[0].Range.Start.Line)
	}
}
