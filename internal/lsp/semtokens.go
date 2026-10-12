package lsp

import (
	"strings"

	"go.lsp.dev/protocol"

	"github.com/ts91/urscript-ls/internal/builtins"
	"github.com/ts91/urscript-ls/internal/parser"
	"github.com/ts91/urscript-ls/internal/symbols"
)

// Token type indices — must match Legend() order.
const (
	ttKeyword   = 0
	ttFunction  = 1
	ttVariable  = 2
	ttParameter = 3
	ttComment   = 4
	ttString    = 5
	ttNumber    = 6
	ttOperator  = 7
	ttDecorator = 8
	ttType      = 9
)

// Token modifier bit flags — must match Legend() order.
const (
	tmDeclaration    = 1 << 0
	tmDefinition     = 1 << 1
	tmDefaultLib     = 1 << 2
	tmDocumentation  = 1 << 3
)

// SemanticTokensLegend returns the semantic tokens legend to register with the client.
func SemanticTokensLegend() protocol.SemanticTokensLegend {
	return protocol.SemanticTokensLegend{
		TokenTypes: []string{
			"keyword",
			"function",
			"variable",
			"parameter",
			"comment",
			"string",
			"number",
			"operator",
			"decorator",
			"type",
		},
		TokenModifiers: []string{
			"declaration",
			"definition",
			"defaultLibrary",
			"documentation",
		},
	}
}

// SemanticTokens returns delta-encoded semantic tokens for the given source.
func SemanticTokens(source string, registry *builtins.Registry, idx *symbols.Index) *protocol.SemanticTokens {
	tokens := parser.Lex(source)

	// Build sets for quick lookup.
	paramSet := make(map[string]bool)
	funcSet := make(map[string]bool)
	if idx != nil {
		for _, s := range idx.Symbols {
			switch s.Kind {
			case symbols.KindParam:
				paramSet[s.Name] = true
			case symbols.KindFunction, symbols.KindThread:
				funcSet[s.Name] = true
			}
		}
	}

	var data []uint32
	prevLine := 0
	prevCol := 0

	for _, tok := range tokens {
		tokenType := -1
		modifiers := uint32(0)

		switch tok.Type {
		case parser.TokenDef, parser.TokenEnd, parser.TokenIf, parser.TokenElif,
			parser.TokenElse, parser.TokenWhile, parser.TokenReturn,
			parser.TokenGlobal, parser.TokenLocal, parser.TokenTrue, parser.TokenFalse,
			parser.TokenNot, parser.TokenAnd, parser.TokenOr, parser.TokenThread:
			tokenType = ttKeyword

		case parser.TokenIdent:
			if _, ok := registry.Lookup(tok.Value); ok {
				tokenType = ttFunction
				modifiers = tmDefaultLib
			} else if funcSet[tok.Value] {
				tokenType = ttFunction
			} else if paramSet[tok.Value] {
				tokenType = ttParameter
			} else {
				tokenType = ttVariable
			}

		case parser.TokenComment:
			tokenType = ttComment

		case parser.TokenDocComment:
			// Emit granular tokens for doc block lines.
			emitDocTokens(tok, &data, &prevLine, &prevCol)
			continue

		case parser.TokenString:
			tokenType = ttString

		case parser.TokenNumber:
			tokenType = ttNumber

		case parser.TokenPlus, parser.TokenMinus, parser.TokenStar, parser.TokenSlash,
			parser.TokenPercent, parser.TokenEq, parser.TokenNeq, parser.TokenLt,
			parser.TokenGt, parser.TokenLte, parser.TokenGte, parser.TokenAssign:
			tokenType = ttOperator

		default:
			continue
		}

		if tokenType < 0 {
			continue
		}

		deltaLine := tok.Line - prevLine
		deltaCol := tok.Column
		if deltaLine == 0 {
			deltaCol = tok.Column - prevCol
		}

		data = append(data, uint32(deltaLine), uint32(deltaCol), uint32(len(tok.Value)), uint32(tokenType), modifiers)

		prevLine = tok.Line
		prevCol = tok.Column
	}

	return &protocol.SemanticTokens{
		Data: data,
	}
}

// emitDocTokens produces granular semantic tokens for a doc comment line.
// It parses @param and @returns tags into separate tokens.
func emitDocTokens(tok parser.Token, data *[]uint32, prevLine, prevCol *int) {
	trimmed := strings.TrimSpace(tok.Value)
	line := tok.Line
	col := tok.Column

	emit := func(l, c, length, tokenType int, modifiers uint32) {
		deltaLine := l - *prevLine
		deltaCol := c
		if deltaLine == 0 {
			deltaCol = c - *prevCol
		}
		*data = append(*data, uint32(deltaLine), uint32(deltaCol), uint32(length), uint32(tokenType), modifiers)
		*prevLine = l
		*prevCol = c
	}

	// ### delimiter lines
	if trimmed == "###" {
		emit(line, col, len(tok.Value), ttComment, tmDocumentation)
		return
	}

	// Find where the content starts (skip leading whitespace).
	contentStart := col
	for i := 0; i < len(tok.Value); i++ {
		if tok.Value[i] != ' ' && tok.Value[i] != '\t' {
			contentStart = col + i
			break
		}
	}
	content := strings.TrimSpace(tok.Value)

	if strings.HasPrefix(content, "@param ") || strings.HasPrefix(content, "@returns ") {
		fields := strings.Fields(content)
		pos := contentStart

		if fields[0] == "@param" {
			// @param
			emit(line, pos, len("@param"), ttDecorator, tmDocumentation)
			pos += len("@param")

			if len(fields) >= 2 {
				// parameter name
				nameIdx := strings.Index(content[len("@param"):], fields[1])
				pos = contentStart + len("@param") + nameIdx
				emit(line, pos, len(fields[1]), ttParameter, tmDocumentation)
				pos += len(fields[1])

				if len(fields) >= 3 {
					// type
					typeIdx := strings.Index(content[pos-contentStart:], fields[2])
					pos = contentStart + (pos - contentStart) + typeIdx
					emit(line, pos, len(fields[2]), ttType, tmDocumentation)
					pos += len(fields[2])

					// remaining description
					if len(fields) >= 4 {
						descStart := strings.Index(content, fields[3])
						if descStart >= 0 {
							descLen := len(content) - descStart
							emit(line, contentStart+descStart, descLen, ttComment, tmDocumentation)
						}
					}
				}
			}
		} else {
			// @returns
			emit(line, pos, len("@returns"), ttDecorator, tmDocumentation)
			pos += len("@returns")

			if len(fields) >= 2 {
				// type
				typeIdx := strings.Index(content[len("@returns"):], fields[1])
				pos = contentStart + len("@returns") + typeIdx
				emit(line, pos, len(fields[1]), ttType, tmDocumentation)
				pos += len(fields[1])

				// remaining description
				if len(fields) >= 3 {
					descStart := strings.Index(content, fields[2])
					if descStart >= 0 {
						descLen := len(content) - descStart
						emit(line, contentStart+descStart, descLen, ttComment, tmDocumentation)
					}
				}
			}
		}
	} else {
		// Plain doc text
		emit(line, col, len(tok.Value), ttComment, tmDocumentation)
	}
}
