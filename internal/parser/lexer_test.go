package parser

import (
	"testing"
)

func TestLexKeywords(t *testing.T) {
	input := "def end if elif else while return global local True False not and or thread"
	tokens := Lex(input)

	expected := []TokenType{
		TokenDef, TokenEnd, TokenIf, TokenElif, TokenElse,
		TokenWhile, TokenReturn, TokenGlobal, TokenLocal,
		TokenTrue, TokenFalse, TokenNot, TokenAnd, TokenOr, TokenThread,
		TokenEOF,
	}

	if len(tokens) != len(expected) {
		t.Fatalf("expected %d tokens, got %d", len(expected), len(tokens))
	}
	for i, tt := range expected {
		if tokens[i].Type != tt {
			t.Errorf("token %d: expected type %v, got %v (%q)", i, tt, tokens[i].Type, tokens[i].Value)
		}
	}
}

func TestLexIdentifiers(t *testing.T) {
	input := "movej my_var _private x1"
	tokens := Lex(input)

	expected := []struct {
		tt  TokenType
		val string
	}{
		{TokenIdent, "movej"},
		{TokenIdent, "my_var"},
		{TokenIdent, "_private"},
		{TokenIdent, "x1"},
		{TokenEOF, ""},
	}

	if len(tokens) != len(expected) {
		t.Fatalf("expected %d tokens, got %d", len(expected), len(tokens))
	}
	for i, e := range expected {
		if tokens[i].Type != e.tt || tokens[i].Value != e.val {
			t.Errorf("token %d: expected (%v, %q), got (%v, %q)", i, e.tt, e.val, tokens[i].Type, tokens[i].Value)
		}
	}
}

func TestLexNumbers(t *testing.T) {
	input := "42 3.14 1e5 2.0E-3"
	tokens := Lex(input)

	expected := []string{"42", "3.14", "1e5", "2.0E-3"}
	idx := 0
	for _, tok := range tokens {
		if tok.Type == TokenNumber {
			if idx >= len(expected) {
				t.Fatalf("too many number tokens")
			}
			if tok.Value != expected[idx] {
				t.Errorf("number %d: expected %q, got %q", idx, expected[idx], tok.Value)
			}
			idx++
		}
	}
	if idx != len(expected) {
		t.Errorf("expected %d numbers, got %d", len(expected), idx)
	}
}

func TestLexStrings(t *testing.T) {
	input := `"hello" "world\"s"`
	tokens := Lex(input)

	var strs []string
	for _, tok := range tokens {
		if tok.Type == TokenString {
			strs = append(strs, tok.Value)
		}
	}

	if len(strs) != 2 {
		t.Fatalf("expected 2 strings, got %d", len(strs))
	}
	if strs[0] != `"hello"` {
		t.Errorf("string 0: expected %q, got %q", `"hello"`, strs[0])
	}
	if strs[1] != `"world\"s"` {
		t.Errorf("string 1: expected %q, got %q", `"world\"s"`, strs[1])
	}
}

func TestLexOperators(t *testing.T) {
	input := "= == != < > <= >= + - * / %"
	tokens := Lex(input)

	expected := []struct {
		tt  TokenType
		val string
	}{
		{TokenAssign, "="},
		{TokenEq, "=="},
		{TokenNeq, "!="},
		{TokenLt, "<"},
		{TokenGt, ">"},
		{TokenLte, "<="},
		{TokenGte, ">="},
		{TokenPlus, "+"},
		{TokenMinus, "-"},
		{TokenStar, "*"},
		{TokenSlash, "/"},
		{TokenPercent, "%"},
		{TokenEOF, ""},
	}

	if len(tokens) != len(expected) {
		t.Fatalf("expected %d tokens, got %d", len(expected), len(tokens))
	}
	for i, e := range expected {
		if tokens[i].Type != e.tt || tokens[i].Value != e.val {
			t.Errorf("token %d: expected (%v, %q), got (%v, %q)", i, e.tt, e.val, tokens[i].Type, tokens[i].Value)
		}
	}
}

func TestLexComment(t *testing.T) {
	input := "x = 1 # this is a comment\ny = 2"
	tokens := Lex(input)

	foundComment := false
	for _, tok := range tokens {
		if tok.Type == TokenComment {
			foundComment = true
			if tok.Value != "# this is a comment" {
				t.Errorf("comment: expected %q, got %q", "# this is a comment", tok.Value)
			}
		}
	}
	if !foundComment {
		t.Error("expected a comment token")
	}
}

func TestLexDelimiters(t *testing.T) {
	input := "()[],:."
	tokens := Lex(input)

	expected := []TokenType{
		TokenLParen, TokenRParen, TokenLBracket, TokenRBracket,
		TokenComma, TokenColon, TokenDot, TokenEOF,
	}

	if len(tokens) != len(expected) {
		t.Fatalf("expected %d tokens, got %d", len(expected), len(tokens))
	}
	for i, tt := range expected {
		if tokens[i].Type != tt {
			t.Errorf("token %d: expected type %v, got %v", i, tt, tokens[i].Type)
		}
	}
}

func TestLexPositions(t *testing.T) {
	input := "x = 1\ny = 2"
	tokens := Lex(input)

	// 'y' should be on line 1, column 0
	for _, tok := range tokens {
		if tok.Type == TokenIdent && tok.Value == "y" {
			if tok.Line != 1 || tok.Column != 0 {
				t.Errorf("'y' expected at (1,0), got (%d,%d)", tok.Line, tok.Column)
			}
			return
		}
	}
	t.Error("did not find identifier 'y'")
}

func TestLexFunctionDef(t *testing.T) {
	input := "def myFunc(a, b):\n  return a + b\nend"
	tokens := Lex(input)

	expectedTypes := []TokenType{
		TokenDef, TokenIdent, TokenLParen, TokenIdent, TokenComma,
		TokenIdent, TokenRParen, TokenColon, TokenNewline,
		TokenReturn, TokenIdent, TokenPlus, TokenIdent, TokenNewline,
		TokenEnd, TokenEOF,
	}

	if len(tokens) != len(expectedTypes) {
		t.Fatalf("expected %d tokens, got %d", len(expectedTypes), len(tokens))
	}
	for i, tt := range expectedTypes {
		if tokens[i].Type != tt {
			t.Errorf("token %d: expected type %v, got %v (%q)", i, tt, tokens[i].Type, tokens[i].Value)
		}
	}
}
