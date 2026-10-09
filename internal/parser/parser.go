package parser

import "fmt"

// Parse tokenizes and parses URScript source code, returning the AST and any
// parse errors encountered.
func Parse(input string) (*Program, []ParseError) {
	tokens := Lex(input)
	p := &parser{tokens: tokens}
	prog := p.parseProgram()
	return prog, p.errors
}

type parser struct {
	tokens []Token
	pos    int
	errors []ParseError
}

func (p *parser) cur() Token {
	if p.pos < len(p.tokens) {
		return p.tokens[p.pos]
	}
	return Token{Type: TokenEOF}
}

func (p *parser) peek() Token {
	next := p.pos + 1
	if next < len(p.tokens) {
		return p.tokens[next]
	}
	return Token{Type: TokenEOF}
}

func (p *parser) advance() Token {
	tok := p.cur()
	if p.pos < len(p.tokens) {
		p.pos++
	}
	return tok
}

func (p *parser) expect(tt TokenType) Token {
	tok := p.cur()
	if tok.Type != tt {
		p.errorf(tok, "expected %v, got %q", tokenName(tt), tok.Value)
		return tok
	}
	return p.advance()
}

func (p *parser) errorf(tok Token, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	p.errors = append(p.errors, ParseError{
		Line:    tok.Line,
		Column:  tok.Column,
		Message: msg,
	})
}

func (p *parser) skipNewlines() {
	for p.cur().Type == TokenNewline || p.cur().Type == TokenComment {
		p.advance()
	}
}

// parseProgram parses the top-level list of statements.
func (p *parser) parseProgram() *Program {
	prog := &Program{}
	p.skipNewlines()
	for p.cur().Type != TokenEOF {
		stmt := p.parseStatement()
		if stmt != nil {
			prog.Stmts = append(prog.Stmts, stmt)
		}
		p.skipNewlines()
	}
	return prog
}

func (p *parser) parseStatement() Node {
	switch p.cur().Type {
	case TokenDef:
		return p.parseFunctionDef()
	case TokenThread:
		return p.parseThreadDef()
	case TokenIf:
		return p.parseIfStmt()
	case TokenWhile:
		return p.parseWhileStmt()
	case TokenFor:
		return p.parseForStmt()
	case TokenReturn:
		return p.parseReturnStmt()
	case TokenGlobal:
		return p.parseGlobalStmt()
	case TokenLocal:
		return p.parseLocalStmt()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *parser) parseFunctionDef() Node {
	tok := p.advance() // consume 'def'
	pos := Pos{Line: tok.Line, Column: tok.Column}

	nameTok := p.expect(TokenIdent)
	name := nameTok.Value

	p.expect(TokenLParen)
	params := p.parseParamList()
	p.expect(TokenRParen)
	p.expect(TokenColon)
	p.skipNewlines()

	body := p.parseBlock(TokenEnd)
	p.expect(TokenEnd)
	p.skipNewlines()

	return &FunctionDef{Pos: pos, Name: name, Params: params, Body: body}
}

func (p *parser) parseThreadDef() Node {
	tok := p.advance() // consume 'thread'
	pos := Pos{Line: tok.Line, Column: tok.Column}

	nameTok := p.expect(TokenIdent)
	name := nameTok.Value

	p.expect(TokenLParen)
	p.expect(TokenRParen)
	p.expect(TokenColon)
	p.skipNewlines()

	body := p.parseBlock(TokenEnd)
	p.expect(TokenEnd)
	p.skipNewlines()

	return &ThreadDef{Pos: pos, Name: name, Body: body}
}

func (p *parser) parseParamList() []string {
	var params []string
	if p.cur().Type == TokenRParen {
		return params
	}
	nameTok := p.expect(TokenIdent)
	params = append(params, nameTok.Value)
	for p.cur().Type == TokenComma {
		p.advance()
		nameTok = p.expect(TokenIdent)
		params = append(params, nameTok.Value)
	}
	return params
}

func (p *parser) parseIfStmt() Node {
	tok := p.advance() // consume 'if'
	pos := Pos{Line: tok.Line, Column: tok.Column}

	cond := p.parseExpression()
	p.expect(TokenColon)
	p.skipNewlines()

	body := p.parseBlock(TokenElif, TokenElse, TokenEnd)

	var elseIfs []ElseIfClause
	for p.cur().Type == TokenElif {
		elifTok := p.advance()
		elifCond := p.parseExpression()
		p.expect(TokenColon)
		p.skipNewlines()
		elifBody := p.parseBlock(TokenElif, TokenElse, TokenEnd)
		elseIfs = append(elseIfs, ElseIfClause{
			Pos:       Pos{Line: elifTok.Line, Column: elifTok.Column},
			Condition: elifCond,
			Body:      elifBody,
		})
	}

	var elseBody []Node
	if p.cur().Type == TokenElse {
		p.advance()
		p.expect(TokenColon)
		p.skipNewlines()
		elseBody = p.parseBlock(TokenEnd)
	}

	p.expect(TokenEnd)
	p.skipNewlines()

	return &IfStmt{Pos: pos, Condition: cond, Body: body, ElseIfs: elseIfs, Else: elseBody}
}

func (p *parser) parseWhileStmt() Node {
	tok := p.advance() // consume 'while'
	pos := Pos{Line: tok.Line, Column: tok.Column}

	cond := p.parseExpression()
	p.expect(TokenColon)
	p.skipNewlines()

	body := p.parseBlock(TokenEnd)
	p.expect(TokenEnd)
	p.skipNewlines()

	return &WhileStmt{Pos: pos, Condition: cond, Body: body}
}

func (p *parser) parseForStmt() Node {
	tok := p.advance() // consume 'for'
	pos := Pos{Line: tok.Line, Column: tok.Column}

	varTok := p.expect(TokenIdent)
	// Expect 'in' as an identifier
	inTok := p.cur()
	if inTok.Type == TokenIdent && inTok.Value == "in" {
		p.advance()
	} else {
		p.errorf(inTok, "expected 'in', got %q", inTok.Value)
	}

	iter := p.parseExpression()
	p.expect(TokenColon)
	p.skipNewlines()

	body := p.parseBlock(TokenEnd)
	p.expect(TokenEnd)
	p.skipNewlines()

	return &ForStmt{Pos: pos, Var: varTok.Value, Iter: iter, Body: body}
}

func (p *parser) parseReturnStmt() Node {
	tok := p.advance() // consume 'return'
	pos := Pos{Line: tok.Line, Column: tok.Column}

	var value Node
	if p.cur().Type != TokenNewline && p.cur().Type != TokenEOF && p.cur().Type != TokenComment {
		value = p.parseExpression()
	}

	return &ReturnStmt{Pos: pos, Value: value}
}

func (p *parser) parseGlobalStmt() Node {
	tok := p.advance() // consume 'global'
	pos := Pos{Line: tok.Line, Column: tok.Column}

	nameTok := p.expect(TokenIdent)
	var init Node
	if p.cur().Type == TokenAssign {
		p.advance()
		init = p.parseExpression()
	}

	return &GlobalStmt{Pos: pos, Name: nameTok.Value, Init: init}
}

func (p *parser) parseLocalStmt() Node {
	tok := p.advance() // consume 'local'
	pos := Pos{Line: tok.Line, Column: tok.Column}

	nameTok := p.expect(TokenIdent)
	var init Node
	if p.cur().Type == TokenAssign {
		p.advance()
		init = p.parseExpression()
	}

	return &LocalStmt{Pos: pos, Name: nameTok.Value, Init: init}
}

func (p *parser) parseExpressionStatement() Node {
	expr := p.parseExpression()
	if expr == nil {
		// Recovery: skip unknown token
		tok := p.advance()
		p.errorf(tok, "unexpected token %q", tok.Value)
		return nil
	}

	// Check for assignment: ident = expr or expr[idx] = expr
	if p.cur().Type == TokenAssign {
		p.advance()
		value := p.parseExpression()
		return &Assignment{
			Pos:    expr.nodePos(),
			Target: expr,
			Value:  value,
		}
	}

	return &ExprStmt{Expr: expr}
}

// parseBlock parses statements until one of the terminators is found.
func (p *parser) parseBlock(terminators ...TokenType) []Node {
	var stmts []Node
	for {
		p.skipNewlines()
		cur := p.cur()
		if cur.Type == TokenEOF {
			break
		}
		for _, t := range terminators {
			if cur.Type == t {
				return stmts
			}
		}
		stmt := p.parseStatement()
		if stmt != nil {
			stmts = append(stmts, stmt)
		}
	}
	return stmts
}

// Expression parsing with precedence climbing

func (p *parser) parseExpression() Node {
	return p.parseOr()
}

func (p *parser) parseOr() Node {
	left := p.parseAnd()
	for p.cur().Type == TokenOr {
		tok := p.advance()
		right := p.parseAnd()
		left = &BinaryExpr{
			Pos:   Pos{Line: tok.Line, Column: tok.Column},
			Op:    "or",
			Left:  left,
			Right: right,
		}
	}
	return left
}

func (p *parser) parseAnd() Node {
	left := p.parseNot()
	for p.cur().Type == TokenAnd {
		tok := p.advance()
		right := p.parseNot()
		left = &BinaryExpr{
			Pos:   Pos{Line: tok.Line, Column: tok.Column},
			Op:    "and",
			Left:  left,
			Right: right,
		}
	}
	return left
}

func (p *parser) parseNot() Node {
	if p.cur().Type == TokenNot {
		tok := p.advance()
		operand := p.parseNot()
		return &UnaryExpr{
			Pos:     Pos{Line: tok.Line, Column: tok.Column},
			Op:      "not",
			Operand: operand,
		}
	}
	return p.parseComparison()
}

func (p *parser) parseComparison() Node {
	left := p.parseAddSub()
	for p.cur().Type == TokenEq || p.cur().Type == TokenNeq ||
		p.cur().Type == TokenLt || p.cur().Type == TokenGt ||
		p.cur().Type == TokenLte || p.cur().Type == TokenGte {
		tok := p.advance()
		right := p.parseAddSub()
		left = &BinaryExpr{
			Pos:   Pos{Line: tok.Line, Column: tok.Column},
			Op:    tok.Value,
			Left:  left,
			Right: right,
		}
	}
	return left
}

func (p *parser) parseAddSub() Node {
	left := p.parseMulDiv()
	for p.cur().Type == TokenPlus || p.cur().Type == TokenMinus {
		tok := p.advance()
		right := p.parseMulDiv()
		left = &BinaryExpr{
			Pos:   Pos{Line: tok.Line, Column: tok.Column},
			Op:    tok.Value,
			Left:  left,
			Right: right,
		}
	}
	return left
}

func (p *parser) parseMulDiv() Node {
	left := p.parseUnary()
	for p.cur().Type == TokenStar || p.cur().Type == TokenSlash || p.cur().Type == TokenPercent {
		tok := p.advance()
		right := p.parseUnary()
		left = &BinaryExpr{
			Pos:   Pos{Line: tok.Line, Column: tok.Column},
			Op:    tok.Value,
			Left:  left,
			Right: right,
		}
	}
	return left
}

func (p *parser) parseUnary() Node {
	if p.cur().Type == TokenMinus {
		tok := p.advance()
		operand := p.parseUnary()
		return &UnaryExpr{
			Pos:     Pos{Line: tok.Line, Column: tok.Column},
			Op:      "-",
			Operand: operand,
		}
	}
	return p.parsePostfix()
}

func (p *parser) parsePostfix() Node {
	expr := p.parsePrimary()
	for {
		switch p.cur().Type {
		case TokenLParen:
			p.advance()
			args := p.parseArgList()
			p.expect(TokenRParen)
			expr = &CallExpr{
				Pos:  expr.nodePos(),
				Func: expr,
				Args: args,
			}
		case TokenLBracket:
			bracketTok := p.advance()
			first := p.parseExpression()
			if p.cur().Type == TokenComma {
				// Pose/list literal: expr[a, b, c, ...]
				elems := []Node{first}
				for p.cur().Type == TokenComma {
					p.advance()
					elems = append(elems, p.parseExpression())
				}
				p.expect(TokenRBracket)
				expr = &IndexExpr{
					Pos:  expr.nodePos(),
					Expr: expr,
					Index: &ListLit{
						Pos:      Pos{Line: bracketTok.Line, Column: bracketTok.Column},
						Elements: elems,
					},
				}
			} else {
				p.expect(TokenRBracket)
				expr = &IndexExpr{
					Pos:   expr.nodePos(),
					Expr:  expr,
					Index: first,
				}
			}
		case TokenDot:
			p.advance()
			field := p.expect(TokenIdent)
			expr = &BinaryExpr{
				Pos:   expr.nodePos(),
				Op:    ".",
				Left:  expr,
				Right: &Identifier{Pos: Pos{Line: field.Line, Column: field.Column}, Name: field.Value},
			}
		default:
			return expr
		}
	}
}

func (p *parser) parsePrimary() Node {
	tok := p.cur()
	switch tok.Type {
	case TokenNumber:
		p.advance()
		return &NumberLit{Pos: Pos{Line: tok.Line, Column: tok.Column}, Value: tok.Value}

	case TokenString:
		p.advance()
		return &StringLit{Pos: Pos{Line: tok.Line, Column: tok.Column}, Value: tok.Value}

	case TokenTrue:
		p.advance()
		return &BoolLit{Pos: Pos{Line: tok.Line, Column: tok.Column}, Value: true}

	case TokenFalse:
		p.advance()
		return &BoolLit{Pos: Pos{Line: tok.Line, Column: tok.Column}, Value: false}

	case TokenIdent:
		p.advance()
		return &Identifier{Pos: Pos{Line: tok.Line, Column: tok.Column}, Name: tok.Value}

	case TokenLParen:
		p.advance()
		expr := p.parseExpression()
		p.expect(TokenRParen)
		return expr

	case TokenLBracket:
		return p.parseListLiteral()

	default:
		return nil
	}
}

func (p *parser) parseListLiteral() Node {
	tok := p.advance() // consume [
	pos := Pos{Line: tok.Line, Column: tok.Column}
	var elems []Node
	if p.cur().Type != TokenRBracket {
		elems = append(elems, p.parseExpression())
		for p.cur().Type == TokenComma {
			p.advance()
			elems = append(elems, p.parseExpression())
		}
	}
	p.expect(TokenRBracket)
	return &ListLit{Pos: pos, Elements: elems}
}

func (p *parser) parseArgList() []Node {
	var args []Node
	if p.cur().Type == TokenRParen {
		return args
	}
	args = append(args, p.parseArgument())
	for p.cur().Type == TokenComma {
		p.advance()
		args = append(args, p.parseArgument())
	}
	return args
}

// parseArgument handles both positional and keyword arguments (name=value).
func (p *parser) parseArgument() Node {
	if p.cur().Type == TokenIdent && p.peek().Type == TokenAssign {
		nameTok := p.advance()
		p.advance() // consume '='
		value := p.parseExpression()
		return &Assignment{
			Pos:    Pos{Line: nameTok.Line, Column: nameTok.Column},
			Target: &Identifier{Pos: Pos{Line: nameTok.Line, Column: nameTok.Column}, Name: nameTok.Value},
			Value:  value,
		}
	}
	return p.parseExpression()
}

func tokenName(tt TokenType) string {
	switch tt {
	case TokenEOF:
		return "EOF"
	case TokenNewline:
		return "newline"
	case TokenIdent:
		return "identifier"
	case TokenNumber:
		return "number"
	case TokenString:
		return "string"
	case TokenDef:
		return "def"
	case TokenEnd:
		return "end"
	case TokenIf:
		return "if"
	case TokenElif:
		return "elif"
	case TokenElse:
		return "else"
	case TokenWhile:
		return "while"
	case TokenFor:
		return "for"
	case TokenReturn:
		return "return"
	case TokenGlobal:
		return "global"
	case TokenLocal:
		return "local"
	case TokenLParen:
		return "("
	case TokenRParen:
		return ")"
	case TokenLBracket:
		return "["
	case TokenRBracket:
		return "]"
	case TokenComma:
		return ","
	case TokenColon:
		return ":"
	case TokenAssign:
		return "="
	default:
		return "token"
	}
}
