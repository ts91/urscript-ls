package lsp

import (
	"strings"
	"testing"

	"go.lsp.dev/protocol"

	"github.com/ts91/urscript-ls/internal/builtins"
	"github.com/ts91/urscript-ls/internal/parser"
	"github.com/ts91/urscript-ls/internal/symbols"
)

func TestHoverBuiltin(t *testing.T) {
	src := "movej(q1)"
	reg := builtins.NewRegistry()
	pos := protocol.Position{Line: 0, Character: 2}

	result := Hover(src, pos, reg, nil)
	if result == nil {
		t.Fatal("expected hover result for builtin")
	}
	mc := result.Contents.(*protocol.MarkupContent)
	if !strings.Contains(mc.Value, "movej") {
		t.Errorf("expected movej in hover, got: %s", mc.Value)
	}
}

func TestHoverUserFunction(t *testing.T) {
	src := `# Moves robot to home position
def go_home():
  movej(q_home)
end

go_home()`
	prog, errs := parser.Parse(src)
	if len(errs) != 0 {
		t.Fatalf("unexpected parse errors: %v", errs)
	}
	idx := symbols.Build(prog, src)
	reg := builtins.NewRegistry()

	// Hover on go_home() call on line 5.
	pos := protocol.Position{Line: 5, Character: 2}
	result := Hover(src, pos, reg, idx)
	if result == nil {
		t.Fatal("expected hover result for user function")
	}
	mc := result.Contents.(*protocol.MarkupContent)
	if !strings.Contains(mc.Value, "def go_home()") {
		t.Errorf("expected 'def go_home()' in hover, got: %s", mc.Value)
	}
	if !strings.Contains(mc.Value, "Moves robot to home position") {
		t.Errorf("expected doc comment in hover, got: %s", mc.Value)
	}
}

func TestHoverUserVariable(t *testing.T) {
	src := `global speed = 1.5
movej(q, v=speed)`
	prog, errs := parser.Parse(src)
	if len(errs) != 0 {
		t.Fatalf("unexpected parse errors: %v", errs)
	}
	idx := symbols.Build(prog, src)
	reg := builtins.NewRegistry()

	pos := protocol.Position{Line: 1, Character: 12}
	result := Hover(src, pos, reg, idx)
	if result == nil {
		t.Fatal("expected hover result for global variable")
	}
	mc := result.Contents.(*protocol.MarkupContent)
	if !strings.Contains(mc.Value, "global speed") {
		t.Errorf("expected 'global speed' in hover, got: %s", mc.Value)
	}
}

func TestHoverBuiltinFormatsParams(t *testing.T) {
	src := "movej(q1)"
	reg := builtins.NewRegistry()
	pos := protocol.Position{Line: 0, Character: 2}

	result := Hover(src, pos, reg, nil)
	if result == nil {
		t.Fatal("expected hover result for builtin")
	}
	mc := result.Contents.(*protocol.MarkupContent)
	if !strings.Contains(mc.Value, "**@param**") {
		t.Errorf("expected formatted @param in builtin hover, got: %s", mc.Value)
	}
	if !strings.Contains(mc.Value, "`q`") {
		t.Errorf("expected backtick-wrapped param name in builtin hover, got: %s", mc.Value)
	}
}

func TestHoverUnknown(t *testing.T) {
	src := "unknown_thing"
	reg := builtins.NewRegistry()
	pos := protocol.Position{Line: 0, Character: 3}

	result := Hover(src, pos, reg, nil)
	if result != nil {
		t.Errorf("expected nil for unknown symbol, got: %v", result)
	}
}
