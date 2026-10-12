package lsp

import (
	"strings"
	"testing"

	"go.lsp.dev/protocol"

	"github.com/ts91/urscript-ls/internal/builtins"
	"github.com/ts91/urscript-ls/internal/parser"
	"github.com/ts91/urscript-ls/internal/symbols"
)

func analyze(src string) []protocol.Diagnostic {
	prog, _ := parser.Parse(src)
	idx := symbols.Build(prog, src)
	reg := builtins.NewRegistry()
	return Analyze(prog, idx, reg)
}

func msgText(d protocol.Diagnostic) string {
	if s, ok := d.Message.(protocol.String); ok {
		return string(s)
	}
	return ""
}

func TestUndefinedFunction(t *testing.T) {
	diags := analyze("not_real()")
	found := false
	for _, d := range diags {
		if strings.Contains(msgText(d), "not defined") {
			found = true
		}
	}
	if !found {
		t.Error("expected 'not defined' warning for unknown function call")
	}
}

func TestDefinedFunctionNoWarning(t *testing.T) {
	src := `def my_func():
  return 1
end
my_func()`
	diags := analyze(src)
	for _, d := range diags {
		if strings.Contains(msgText(d), "not defined") {
			t.Errorf("unexpected 'not defined' warning: %s", msgText(d))
		}
	}
}

func TestBuiltinNoWarning(t *testing.T) {
	diags := analyze("movej(q1)")
	for _, d := range diags {
		if strings.Contains(msgText(d), "not defined") {
			t.Errorf("unexpected 'not defined' warning for builtin: %s", msgText(d))
		}
	}
}

func TestTooManyArgsBuiltin(t *testing.T) {
	// sleep takes 1 param.
	diags := analyze("sleep(1, 2, 3)")
	found := false
	for _, d := range diags {
		if strings.Contains(msgText(d), "expects at most") {
			found = true
		}
	}
	if !found {
		t.Error("expected arg count warning for sleep(1, 2, 3)")
	}
}

func TestTooManyArgsUser(t *testing.T) {
	src := `def foo(a):
  return a
end
foo(1, 2, 3)`
	diags := analyze(src)
	found := false
	for _, d := range diags {
		if strings.Contains(msgText(d), "expects at most") {
			found = true
		}
	}
	if !found {
		t.Error("expected arg count warning for foo(1, 2, 3)")
	}
}

func TestDuplicateFunction(t *testing.T) {
	src := `def foo():
  return 1
end
def foo():
  return 2
end`
	diags := analyze(src)
	found := false
	for _, d := range diags {
		if strings.Contains(msgText(d), "already defined") {
			found = true
		}
	}
	if !found {
		t.Error("expected 'already defined' warning for duplicate function")
	}
}
