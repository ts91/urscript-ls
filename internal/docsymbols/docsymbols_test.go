package docsymbols

import (
	"testing"

	"go.lsp.dev/protocol"

	"github.com/ts91/urscript-ls/internal/parser"
	"github.com/ts91/urscript-ls/internal/symbols"
)

func TestProvideEmpty(t *testing.T) {
	result := Provide(nil)
	if len(result) != 0 {
		t.Fatalf("expected empty slice for nil index, got %v", result)
	}
}

func TestProvideFunctionAndVariable(t *testing.T) {
	src := `def my_func(a, b):
  x = 1
end`
	prog, errs := parser.Parse(src)
	if len(errs) != 0 {
		t.Fatalf("unexpected parse errors: %v", errs)
	}
	idx := symbols.Build(prog, src)
	result := Provide(idx)

	if len(result) == 0 {
		t.Fatal("expected at least one document symbol")
	}

	// Check that my_func is present as a function.
	found := false
	for _, ds := range result {
		if ds.Name == "my_func" && ds.Kind == protocol.SymbolKindFunction {
			found = true
			if ds.Detail == nil || *ds.Detail != "(a, b)" {
				t.Errorf("expected detail '(a, b)', got %v", ds.Detail)
			}
			break
		}
	}
	if !found {
		t.Error("my_func not found in document symbols")
	}

	// Variable x should be present.
	foundVar := false
	for _, ds := range result {
		if ds.Name == "x" && ds.Kind == protocol.SymbolKindVariable {
			foundVar = true
			break
		}
	}
	if !foundVar {
		t.Error("variable x not found in document symbols")
	}
}

func TestProvideThread(t *testing.T) {
	src := `thread worker():
  sync()
end`
	prog, errs := parser.Parse(src)
	if len(errs) != 0 {
		t.Fatalf("unexpected parse errors: %v", errs)
	}
	idx := symbols.Build(prog, src)
	result := Provide(idx)

	found := false
	for _, ds := range result {
		if ds.Name == "worker" {
			found = true
			if ds.Detail == nil || *ds.Detail != "thread" {
				t.Errorf("expected detail 'thread', got %v", ds.Detail)
			}
			break
		}
	}
	if !found {
		t.Error("thread worker not found in document symbols")
	}
}

func TestSkipsParams(t *testing.T) {
	src := `def foo(param1):
  return param1
end`
	prog, errs := parser.Parse(src)
	if len(errs) != 0 {
		t.Fatalf("unexpected parse errors: %v", errs)
	}
	idx := symbols.Build(prog, src)
	result := Provide(idx)

	for _, ds := range result {
		if ds.Name == "param1" {
			t.Error("params should not appear as top-level document symbols")
		}
	}
}
