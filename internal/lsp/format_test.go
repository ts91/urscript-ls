package lsp

import (
	"testing"

	"go.lsp.dev/protocol"
)

var defaultOpts = protocol.FormattingOptions{TabSize: 2, InsertSpaces: true}

func TestFormatBasicIndent(t *testing.T) {
	input := "def foo():\nx = 1\nend\n"
	expected := "def foo():\n  x = 1\nend\n"

	edits := Format(input, defaultOpts)
	if len(edits) == 0 {
		t.Fatal("expected edits")
	}
	if edits[0].NewText != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, edits[0].NewText)
	}
}

func TestFormatIfElse(t *testing.T) {
	input := "if x > 0:\ny = 1\nelif x < 0:\ny = -1\nelse:\ny = 0\nend\n"
	expected := "if x > 0:\n  y = 1\nelif x < 0:\n  y = -1\nelse:\n  y = 0\nend\n"

	edits := Format(input, defaultOpts)
	if len(edits) == 0 {
		t.Fatal("expected edits")
	}
	if edits[0].NewText != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, edits[0].NewText)
	}
}

func TestFormatNested(t *testing.T) {
	input := "def foo():\nif True:\nx = 1\nend\nend\n"
	expected := "def foo():\n  if True:\n    x = 1\n  end\nend\n"

	edits := Format(input, defaultOpts)
	if len(edits) == 0 {
		t.Fatal("expected edits")
	}
	if edits[0].NewText != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, edits[0].NewText)
	}
}

func TestFormatAlreadyFormatted(t *testing.T) {
	input := "def foo():\n  x = 1\nend\n"
	edits := Format(input, defaultOpts)
	if len(edits) != 0 {
		t.Errorf("expected no edits for already formatted code, got %d", len(edits))
	}
}

func TestFormatTabs(t *testing.T) {
	opts := protocol.FormattingOptions{TabSize: 4, InsertSpaces: false}
	input := "def foo():\nx = 1\nend\n"
	expected := "def foo():\n\tx = 1\nend\n"

	edits := Format(input, opts)
	if len(edits) == 0 {
		t.Fatal("expected edits")
	}
	if edits[0].NewText != expected {
		t.Errorf("expected:\n%q\ngot:\n%q", expected, edits[0].NewText)
	}
}

func TestFormatTrimsTrailingWhitespace(t *testing.T) {
	input := "def foo():   \n  x = 1   \nend   \n"
	expected := "def foo():\n  x = 1\nend\n"

	edits := Format(input, defaultOpts)
	if len(edits) == 0 {
		t.Fatal("expected edits")
	}
	if edits[0].NewText != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, edits[0].NewText)
	}
}

func TestFormatPreservesEmptyLines(t *testing.T) {
	input := "x = 1\n\ny = 2\n"

	edits := Format(input, defaultOpts)
	if len(edits) != 0 {
		t.Errorf("expected no edits, got:\n%s", edits[0].NewText)
	}
}

func TestFormatThread(t *testing.T) {
	input := "thread worker():\ndo_work()\nend\n"
	expected := "thread worker():\n  do_work()\nend\n"

	edits := Format(input, defaultOpts)
	if len(edits) == 0 {
		t.Fatal("expected edits")
	}
	if edits[0].NewText != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, edits[0].NewText)
	}
}

func TestFormatWhile(t *testing.T) {
	input := "while True:\nx = 1\nend\n"
	expected := "while True:\n  x = 1\nend\n"

	edits := Format(input, defaultOpts)
	if len(edits) == 0 {
		t.Fatal("expected edits")
	}
	if edits[0].NewText != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, edits[0].NewText)
	}
}
