package main

import (
	"fmt"
	"io/ioutil"
	"os"

	cionom "github.com/Th3T3chn0G1t/CioGo/lib"
)

func main() {
	// TODO: Proper argparsing
	cionom.Assert("No source file specified", len(os.Args) < 2)

	RawSource, Error := ioutil.ReadFile(os.Args[1])
	cionom.AssertError(Error, "reading source file")
	var Source string = string(RawSource)

	var Tokens []cionom.Token = cionom.Tokenize(Source)

	Program, Error := cionom.Parse(Tokens, Source)
	cionom.AssertError(Error, "parsing program")

	fmt.Println(Program)
}
