package cionom

import (
	"errors"
	"fmt"
	"strconv"
)

type Parameter uint8

type Call struct {
	Identifier string
	Parameters []Parameter
}

type Routine struct {
	Identifier string
	External   bool
	Calls      []Call
}

type Program struct {
	Routines []Routine
}

type TokenPattern struct {
	Kind      TokenKind
	Instances TokenPosition
}

func ParseExpectSequence(Tokens []Token, Start TokenPosition, Sequence []TokenPattern) (TokenPosition, error) {
	var Pattern TokenPosition = 0
	var PatternRepetitions TokenPosition = 0
	var Position TokenPosition = TokenPosition(Start)
	for ; Pattern < TokenPosition(len(Sequence)); Position++ {
		if Position >= TokenPosition(len(Tokens)) {
			return 0, errors.New("unexpected EOF")
		}

		var Token Token = Tokens[Position]
		var Expected TokenPattern = Sequence[Pattern]

		if Token.Kind != Expected.Kind && Expected.Instances != 0 {
			return 0, fmt.Errorf("expected %s but got %s at offset %v", TokenKindString(Expected.Kind), TokenKindString(Token.Kind), Token.SourceSpan.Begin)
		} else if Token.Kind != Expected.Kind && Expected.Instances == 0 {
			PatternRepetitions = 0
			Pattern++
			Position--
			continue
		}

		PatternRepetitions++
		if Expected.Instances != 0 && PatternRepetitions >= Expected.Instances {
			PatternRepetitions = 0
			Pattern++
		}
	}

	return Position - Start, nil
}

func Parse(Tokens []Token, Source string) (Program, error) {
	var Program Program

	for Position := TokenPosition(0); Position < TokenPosition(len(Tokens)); {
		var Routine Routine

		Stride, Error := ParseExpectSequence(Tokens, Position, []TokenPattern{{Identifier, 1}, {Number, 1}})
		if Error != nil {
			return Program, Error
		}

		var Token Token = Tokens[Position]
		Routine.Identifier = Source[Token.SourceSpan.Begin:Token.SourceSpan.End]
		Routine.External = true

		Position += Stride

		if Tokens[Position].Kind == Block {
			Stride, Error = ParseExpectSequence(Tokens, Position, []TokenPattern{{Block, 1}})
			if Error != nil {
				return Program, Error
			}
			Position += Stride

			Routine.External = false

			for Tokens[Position].Kind != Block {
				var Call Call

				Stride, Error = ParseExpectSequence(Tokens, Position, []TokenPattern{{Identifier, 1}, {Number, 0}})
				if Error != nil {
					return Program, Error
				}

				Token = Tokens[Position]
				Call.Identifier = Source[Token.SourceSpan.Begin:Token.SourceSpan.End]

				for ParameterPosition := Position + 1; ParameterPosition < Position+Stride; ParameterPosition++ {
					Token = Tokens[ParameterPosition]
					Value, Error := strconv.Atoi(Source[Token.SourceSpan.Begin:Token.SourceSpan.End])
					if Error != nil {
						return Program, Error
					}
					Call.Parameters = append(Call.Parameters, Parameter(Value))
				}

				Position += Stride

				Routine.Calls = append(Routine.Calls, Call)
			}

			Stride, Error = ParseExpectSequence(Tokens, Position, []TokenPattern{{Block, 1}})
			if Error != nil {
				return Program, Error
			}
			Position += Stride
		}

		Program.Routines = append(Program.Routines, Routine)
	}

	return Program, nil
}
