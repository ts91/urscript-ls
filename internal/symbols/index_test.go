package symbols

import (
	"testing"

	"github.com/ts91/urscript-ls/internal/parser"
)

func TestBuildFunctions(t *testing.T) {
	src := `# Moves to home position
def go_home():
  movej(home_pos)
end

def add(a, b):
  return a + b
end`
	prog, errs := parser.Parse(src)
	if len(errs) != 0 {
		t.Fatalf("parse errors: %v", errs)
	}

	idx := Build(prog, src)
	fns := idx.Functions()
	if len(fns) != 2 {
		t.Fatalf("expected 2 functions, got %d", len(fns))
	}

	if fns[0].Name != "go_home" {
		t.Errorf("expected 'go_home', got %q", fns[0].Name)
	}
	if fns[0].Doc != "Moves to home position" {
		t.Errorf("expected doc 'Moves to home position', got %q", fns[0].Doc)
	}
	if len(fns[0].Params) != 0 {
		t.Errorf("expected 0 params, got %d", len(fns[0].Params))
	}

	if fns[1].Name != "add" {
		t.Errorf("expected 'add', got %q", fns[1].Name)
	}
	if len(fns[1].Params) != 2 {
		t.Errorf("expected 2 params, got %d", len(fns[1].Params))
	}
}

func TestBuildVariables(t *testing.T) {
	src := `global speed = 1.05
local offset = 0.1
x = 42`
	prog, errs := parser.Parse(src)
	if len(errs) != 0 {
		t.Fatalf("parse errors: %v", errs)
	}

	idx := Build(prog, src)
	vars := idx.Variables()
	if len(vars) != 3 {
		t.Fatalf("expected 3 variables, got %d", len(vars))
	}

	kinds := map[string]SymbolKind{}
	for _, v := range vars {
		kinds[v.Name] = v.Kind
	}
	if kinds["speed"] != KindGlobal {
		t.Errorf("expected 'speed' to be global")
	}
	if kinds["offset"] != KindLocal {
		t.Errorf("expected 'offset' to be local")
	}
	if kinds["x"] != KindVariable {
		t.Errorf("expected 'x' to be variable")
	}
}

func TestBuildThread(t *testing.T) {
	src := `thread monitor():
  while True:
    sleep(0.1)
  end
end`
	prog, errs := parser.Parse(src)
	if len(errs) != 0 {
		t.Fatalf("parse errors: %v", errs)
	}

	idx := Build(prog, src)
	threads := idx.Threads()
	if len(threads) != 1 {
		t.Fatalf("expected 1 thread, got %d", len(threads))
	}
	if threads[0].Name != "monitor" {
		t.Errorf("expected 'monitor', got %q", threads[0].Name)
	}
}

func TestLookup(t *testing.T) {
	src := `def foo():
  x = 1
end
def bar():
  x = 2
end`
	prog, errs := parser.Parse(src)
	if len(errs) != 0 {
		t.Fatalf("parse errors: %v", errs)
	}

	idx := Build(prog, src)

	foos := idx.Lookup("foo")
	if len(foos) != 1 {
		t.Fatalf("expected 1 'foo', got %d", len(foos))
	}
	if foos[0].Kind != KindFunction {
		t.Errorf("expected function, got %v", foos[0].Kind)
	}

	xs := idx.Lookup("x")
	if len(xs) != 2 {
		t.Errorf("expected 2 'x' symbols, got %d", len(xs))
	}
}

func TestDocExtraction(t *testing.T) {
	src := `# Line 1 of doc
# Line 2 of doc
def documented():
  pass
end`
	prog, errs := parser.Parse(src)
	if len(errs) != 0 {
		t.Fatalf("parse errors: %v", errs)
	}

	idx := Build(prog, src)
	fns := idx.Functions()
	if len(fns) != 1 {
		t.Fatalf("expected 1 function, got %d", len(fns))
	}
	expected := "Line 1 of doc\nLine 2 of doc"
	if fns[0].Doc != expected {
		t.Errorf("expected doc %q, got %q", expected, fns[0].Doc)
	}
}
