package main

import (
	"os"

	cionom "github.com/Th3T3chn0G1t/CioGo/lib"
)

func main() {
	// TODO: Proper argparsing
	cionom.Assert("No source file specified", len(os.Args) < 2)
	cionom.Assert("No output file specified", len(os.Args) < 3)

	RawSource, Error := os.ReadFile(os.Args[1])
	cionom.AssertError(Error, "reading source file")
	var Source string = string(RawSource)

	var Tokens []cionom.Token = cionom.Tokenize(Source)

	Program, Error := cionom.Parse(Tokens, Source)
	cionom.AssertError(Error, "parsing program")

	var Bytecode cionom.Bytecode = cionom.Emit(Program)

	os.WriteFile(os.Args[2], Bytecode, 0666)
}
