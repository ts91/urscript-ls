package lsp

import (
	"strings"
	"testing"

	"go.lsp.dev/protocol"

	"github.com/ts91/urscript-ls/internal/builtins"
	"github.com/ts91/urscript-ls/internal/parser"
	"github.com/ts91/urscript-ls/internal/symbols"
)

func TestSignatureHelpBuiltin(t *testing.T) {
	src := "movej(q1, "
	reg := builtins.NewRegistry()
	pos := protocol.Position{Line: 0, Character: uint32(len(src))}

	result := SignatureHelp(src, pos, reg, nil)
	if result == nil {
		t.Fatal("expected signature help for builtin")
	}
	if len(result.Signatures) != 1 {
		t.Fatalf("expected 1 signature, got %d", len(result.Signatures))
	}
	sig := result.Signatures[0]
	if !strings.Contains(sig.Label, "movej") {
		t.Errorf("expected movej in label, got: %s", sig.Label)
	}
	if sig.Documentation == nil {
		t.Error("expected documentation on signature")
	}
	activeParam, ok := result.ActiveParameter.Get()
	if !ok {
		t.Fatal("expected active parameter to be set")
	}
	if activeParam != 1 {
		t.Errorf("expected active param 1, got %d", activeParam)
	}
}

func TestSignatureHelpParamDocs(t *testing.T) {
	src := "conveyor_pulse_decode(type1, "
	reg := builtins.NewRegistry()
	pos := protocol.Position{Line: 0, Character: uint32(len(src))}

	result := SignatureHelp(src, pos, reg, nil)
	if result == nil {
		t.Fatal("expected signature help")
	}
	sig := result.Signatures[0]
	if len(sig.Parameters) < 3 {
		t.Fatalf("expected at least 3 params, got %d", len(sig.Parameters))
	}

	// First param "type" should have documentation from @param tag.
	p0 := sig.Parameters[0]
	if p0.Documentation == nil {
		t.Fatal("expected documentation on first parameter")
	}
	doc0 := p0.Documentation.(*protocol.MarkupContent).Value
	if !strings.Contains(doc0, "int") {
		t.Errorf("expected type info in param doc, got: %s", doc0)
	}
	if !strings.Contains(doc0, "Encoder type") {
		t.Errorf("expected description in param doc, got: %s", doc0)
	}
}

func TestSignatureHelpUserFunction(t *testing.T) {
	parseSrc := `def add(a, b):
  return a + b
end`
	prog, errs := parser.Parse(parseSrc)
	if len(errs) != 0 {
		t.Fatalf("unexpected parse errors: %v", errs)
	}
	idx := symbols.Build(prog, parseSrc)
	reg := builtins.NewRegistry()

	src := parseSrc + "\nadd(1, "
	pos := protocol.Position{Line: 3, Character: 7}
	result := SignatureHelp(src, pos, reg, idx)
	if result == nil {
		t.Fatal("expected signature help for user function")
	}
	sig := result.Signatures[0]
	if sig.Label != "add(a, b)" {
		t.Errorf("expected 'add(a, b)', got %q", sig.Label)
	}
	if len(sig.Parameters) != 2 {
		t.Fatalf("expected 2 params, got %d", len(sig.Parameters))
	}
	// User functions have no @param docs.
	if sig.Parameters[0].Documentation != nil {
		t.Error("expected no documentation on user function params")
	}
	activeParam, ok := result.ActiveParameter.Get()
	if !ok {
		t.Fatal("expected active parameter to be set")
	}
	if activeParam != 1 {
		t.Errorf("expected active param 1, got %d", activeParam)
	}
}

func TestSignatureHelpNoContext(t *testing.T) {
	src := "x = 42"
	reg := builtins.NewRegistry()
	pos := protocol.Position{Line: 0, Character: 6}

	result := SignatureHelp(src, pos, reg, nil)
	if result != nil {
		t.Errorf("expected nil when not inside a call, got: %v", result)
	}
}

func TestFindCallContextCommaCount(t *testing.T) {
	tests := []struct {
		text      string
		wantName  string
		wantParam int
	}{
		{"foo(", "foo", 0},
		{"foo(a, ", "foo", 1},
		{"foo(a, b, ", "foo", 2},
		{"foo(bar(x), ", "foo", 1},
		{`foo("a,b", `, "foo", 1},
	}
	for _, tc := range tests {
		name, idx := findCallContext(tc.text)
		if name != tc.wantName {
			t.Errorf("findCallContext(%q): name = %q, want %q", tc.text, name, tc.wantName)
		}
		if idx != tc.wantParam {
			t.Errorf("findCallContext(%q): param = %d, want %d", tc.text, idx, tc.wantParam)
		}
	}
}

func TestExtractParamDocs(t *testing.T) {
	doc := "Do something.\n@param a int First arg\n@param b float Second arg"
	m := extractParamDocs(doc)
	if len(m) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(m))
	}
	if !strings.Contains(m["a"], "`int`") {
		t.Errorf("expected type in param doc for a, got: %s", m["a"])
	}
	if !strings.Contains(m["a"], "First arg") {
		t.Errorf("expected description in param doc for a, got: %s", m["a"])
	}
	if !strings.Contains(m["b"], "`float`") {
		t.Errorf("expected type in param doc for b, got: %s", m["b"])
	}
}

func TestExtractParamDocsEmpty(t *testing.T) {
	m := extractParamDocs("")
	if m != nil {
		t.Errorf("expected nil for empty doc, got: %v", m)
	}
}
