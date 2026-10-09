package semtokens

import (
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
)

// Token modifier bit flags — must match Legend() order.
const (
	tmDeclaration  = 1 << 0
	tmDefinition   = 1 << 1
	tmDefaultLib   = 1 << 2
)

// Legend returns the semantic tokens legend to register with the client.
func Legend() protocol.SemanticTokensLegend {
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
		},
		TokenModifiers: []string{
			"declaration",
			"definition",
			"defaultLibrary",
		},
	}
}

// Provide returns delta-encoded semantic tokens for the given source.
func Provide(source string, registry *builtins.Registry, idx *symbols.Index) *protocol.SemanticTokens {
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
			parser.TokenElse, parser.TokenWhile, parser.TokenFor, parser.TokenReturn,
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
