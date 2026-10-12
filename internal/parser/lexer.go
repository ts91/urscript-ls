package parser

import (
	"strings"
	"unicode"
)

// TokenType identifies the kind of token.
type TokenType int

const (
	// Special tokens
	TokenEOF TokenType = iota
	TokenIllegal
	TokenNewline

	// Literals
	TokenIdent
	TokenNumber
	TokenString

	// Keywords
	TokenDef
	TokenEnd
	TokenIf
	TokenElif
	TokenElse
	TokenWhile
	TokenReturn
	TokenGlobal
	TokenLocal
	TokenTrue
	TokenFalse
	TokenNot
	TokenAnd
	TokenOr
	TokenThread

	// Operators
	TokenAssign    // =
	TokenPlus      // +
	TokenMinus     // -
	TokenStar      // *
	TokenSlash     // /
	TokenPercent   // %
	TokenEq        // ==
	TokenNeq       // !=
	TokenLt        // <
	TokenGt        // >
	TokenLte       // <=
	TokenGte       // >=

	// Delimiters
	TokenLParen   // (
	TokenRParen   // )
	TokenLBracket // [
	TokenRBracket // ]
	TokenComma    // ,
	TokenColon    // :
	TokenDot      // .

	// Comments
	TokenComment
	TokenDocComment // line inside a ### block
)

var keywords = map[string]TokenType{
	"def":    TokenDef,
	"end":    TokenEnd,
	"if":     TokenIf,
	"elif":   TokenElif,
	"else":   TokenElse,
	"while":  TokenWhile,
	"return": TokenReturn,
	"global": TokenGlobal,
	"local":  TokenLocal,
	"True":   TokenTrue,
	"False":  TokenFalse,
	"not":    TokenNot,
	"and":    TokenAnd,
	"or":     TokenOr,
	"thread": TokenThread,
}

// Token represents a single lexical token.
type Token struct {
	Type   TokenType
	Value  string
	Line   int // 0-based
	Column int // 0-based
}

// Lexer tokenizes URScript source code.
type Lexer struct {
	input  string
	pos    int
	line   int
	col    int
	tokens []Token
}

// Lex tokenizes the entire input and returns the token slice.
func Lex(input string) []Token {
	l := &Lexer{input: input}
	l.tokenize()
	return l.tokens
}

func (l *Lexer) tokenize() {
	for l.pos < len(l.input) {
		ch := l.input[l.pos]

		switch {
		case ch == '\n':
			l.emit(TokenNewline, "\n")
			l.pos++
			l.line++
			l.col = 0

		case ch == '\r':
			l.pos++
			l.col++

		case ch == ' ' || ch == '\t':
			l.pos++
			l.col++

		case ch == '#':
			l.lexComment()

		case ch == '"':
			l.lexString()

		case isDigit(ch) || (ch == '.' && l.pos+1 < len(l.input) && isDigit(l.input[l.pos+1])):
			l.lexNumber()

		case isIdentStart(ch):
			l.lexIdent()

		case ch == '=':
			if l.peek(1) == '=' {
				l.emit(TokenEq, "==")
				l.advance(2)
			} else {
				l.emit(TokenAssign, "=")
				l.advance(1)
			}

		case ch == '!':
			if l.peek(1) == '=' {
				l.emit(TokenNeq, "!=")
				l.advance(2)
			} else {
				l.emit(TokenIllegal, string(ch))
				l.advance(1)
			}

		case ch == '<':
			if l.peek(1) == '=' {
				l.emit(TokenLte, "<=")
				l.advance(2)
			} else {
				l.emit(TokenLt, "<")
				l.advance(1)
			}

		case ch == '>':
			if l.peek(1) == '=' {
				l.emit(TokenGte, ">=")
				l.advance(2)
			} else {
				l.emit(TokenGt, ">")
				l.advance(1)
			}

		case ch == '+':
			l.emit(TokenPlus, "+")
			l.advance(1)
		case ch == '-':
			l.emit(TokenMinus, "-")
			l.advance(1)
		case ch == '*':
			l.emit(TokenStar, "*")
			l.advance(1)
		case ch == '/':
			l.emit(TokenSlash, "/")
			l.advance(1)
		case ch == '%':
			l.emit(TokenPercent, "%")
			l.advance(1)

		case ch == '(':
			l.emit(TokenLParen, "(")
			l.advance(1)
		case ch == ')':
			l.emit(TokenRParen, ")")
			l.advance(1)
		case ch == '[':
			l.emit(TokenLBracket, "[")
			l.advance(1)
		case ch == ']':
			l.emit(TokenRBracket, "]")
			l.advance(1)
		case ch == ',':
			l.emit(TokenComma, ",")
			l.advance(1)
		case ch == ':':
			l.emit(TokenColon, ":")
			l.advance(1)
		case ch == '.':
			l.emit(TokenDot, ".")
			l.advance(1)

		default:
			l.emit(TokenIllegal, string(ch))
			l.advance(1)
		}
	}

	l.emit(TokenEOF, "")
}

func (l *Lexer) emit(tt TokenType, val string) {
	l.tokens = append(l.tokens, Token{
		Type:   tt,
		Value:  val,
		Line:   l.line,
		Column: l.col,
	})
}

func (l *Lexer) advance(n int) {
	l.pos += n
	l.col += n
}

func (l *Lexer) peek(offset int) byte {
	i := l.pos + offset
	if i < len(l.input) {
		return l.input[i]
	}
	return 0
}

func (l *Lexer) lexComment() {
	start := l.pos
	startCol := l.col
	for l.pos < len(l.input) && l.input[l.pos] != '\n' {
		l.pos++
		l.col++
	}
	text := l.input[start:l.pos]

	// Check if this is a ### doc block delimiter.
	if strings.TrimSpace(text) == "###" {
		l.tokens = append(l.tokens, Token{
			Type:   TokenDocComment,
			Value:  text,
			Line:   l.line,
			Column: startCol,
		})
		l.lexDocBlock()
		return
	}

	l.tokens = append(l.tokens, Token{
		Type:   TokenComment,
		Value:  text,
		Line:   l.line,
		Column: startCol,
	})
}

// lexDocBlock consumes lines after an opening ### until a closing ### is found,
// emitting each line as a TokenDocComment.
func (l *Lexer) lexDocBlock() {
	for l.pos < len(l.input) {
		// Skip the newline after the previous line.
		if l.input[l.pos] == '\n' {
			l.pos++
			l.line++
			l.col = 0
		}

		// Capture the current line.
		lineStart := l.pos
		lineCol := l.col
		lineLine := l.line
		for l.pos < len(l.input) && l.input[l.pos] != '\n' {
			l.pos++
			l.col++
		}
		text := l.input[lineStart:l.pos]

		l.tokens = append(l.tokens, Token{
			Type:   TokenDocComment,
			Value:  text,
			Line:   lineLine,
			Column: lineCol,
		})

		// If this line is the closing ###, stop.
		if strings.TrimSpace(text) == "###" {
			return
		}
	}
}

func (l *Lexer) lexString() {
	startCol := l.col
	l.pos++ // skip opening "
	l.col++
	var sb strings.Builder
	sb.WriteByte('"')

	for l.pos < len(l.input) && l.input[l.pos] != '"' && l.input[l.pos] != '\n' {
		if l.input[l.pos] == '\\' && l.pos+1 < len(l.input) {
			sb.WriteByte(l.input[l.pos])
			l.pos++
			l.col++
			sb.WriteByte(l.input[l.pos])
			l.pos++
			l.col++
			continue
		}
		sb.WriteByte(l.input[l.pos])
		l.pos++
		l.col++
	}

	if l.pos < len(l.input) && l.input[l.pos] == '"' {
		sb.WriteByte('"')
		l.pos++
		l.col++
		l.tokens = append(l.tokens, Token{
			Type:   TokenString,
			Value:  sb.String(),
			Line:   l.line,
			Column: startCol,
		})
	} else {
		l.tokens = append(l.tokens, Token{
			Type:   TokenIllegal,
			Value:  sb.String(),
			Line:   l.line,
			Column: startCol,
		})
	}
}

func (l *Lexer) lexNumber() {
	start := l.pos
	startCol := l.col
	hasDot := false
	hasE := false

	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if isDigit(ch) {
			l.pos++
			l.col++
		} else if ch == '.' && !hasDot && !hasE {
			hasDot = true
			l.pos++
			l.col++
		} else if (ch == 'e' || ch == 'E') && !hasE {
			hasE = true
			l.pos++
			l.col++
			if l.pos < len(l.input) && (l.input[l.pos] == '+' || l.input[l.pos] == '-') {
				l.pos++
				l.col++
			}
		} else {
			break
		}
	}

	l.tokens = append(l.tokens, Token{
		Type:   TokenNumber,
		Value:  l.input[start:l.pos],
		Line:   l.line,
		Column: startCol,
	})
}

func (l *Lexer) lexIdent() {
	start := l.pos
	startCol := l.col
	for l.pos < len(l.input) && isIdentPart(l.input[l.pos]) {
		l.pos++
		l.col++
	}
	word := l.input[start:l.pos]
	tt, ok := keywords[word]
	if !ok {
		tt = TokenIdent
	}
	l.tokens = append(l.tokens, Token{
		Type:   tt,
		Value:  word,
		Line:   l.line,
		Column: startCol,
	})
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isIdentStart(ch byte) bool {
	return ch == '_' || unicode.IsLetter(rune(ch))
}

func isIdentPart(ch byte) bool {
	return ch == '_' || unicode.IsLetter(rune(ch)) || isDigit(ch)
}
