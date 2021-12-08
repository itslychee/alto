package main

import (
	"github.com/ItsLychee/alto/dsl"
)

func ParseFormatString(s string) ([]dsl.ASTNode, error) {
	toks, err := dsl.NewLexer(s).Lex()

	if err != nil {
		return nil, err
	}
	return dsl.NewParser(toks).Parse()
}