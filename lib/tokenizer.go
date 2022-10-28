package cionom

import "unicode"

type SourcePosition uint32
type TokenPosition uint32

type Span struct {
	Begin SourcePosition
	End   SourcePosition
}

type TokenKind uint8

const (
	Block      TokenKind = 0
	Identifier TokenKind = 1
	Number     TokenKind = 2
)

type Token struct {
	Kind       TokenKind
	SourceSpan Span
}

func TokenKindString(Kind TokenKind) string {
	switch Kind {
	case Block:
		return "Block"
	case Identifier:
		return "Identifier"
	case Number:
		return "Number"
	}
	return ""
}

func TokenSubstringFromCondition(String string, Kind TokenKind, Start SourcePosition, Condition func(rune) bool) Token {
	var Token Token
	Token.Kind = Kind
	Token.SourceSpan.Begin = Start
	Token.SourceSpan.End = Start + 1

	for Condition(rune(String[Token.SourceSpan.End])) && Token.SourceSpan.End < SourcePosition(len(String)) {
		Token.SourceSpan.End++
	}

	return Token
}

func Tokenize(Source string) []Token {
	var Tokens []Token
	for Position := SourcePosition(0); Position < SourcePosition(len(Source)); Position++ {
		if unicode.IsSpace(rune(Source[Position])) {
			continue
		}

		if Source[Position] == ':' {
			Tokens = append(Tokens, Token{Block, Span{Position, Position + 1}})
		} else if unicode.IsDigit(rune(Source[Position])) {
			Tokens = append(Tokens, TokenSubstringFromCondition(Source, Number, Position, unicode.IsDigit))
		} else {
			var IsNotSpace = func(Rune rune) bool {
				return !unicode.IsSpace(Rune)
			}
			Tokens = append(Tokens, TokenSubstringFromCondition(Source, Identifier, Position, IsNotSpace))
		}

		Position = Tokens[len(Tokens)-1].SourceSpan.End - 1
	}

	return Tokens
}
