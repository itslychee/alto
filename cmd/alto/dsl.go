package main

import (
	"github.com/ItsLychee/alto/dsl"
)

func ParseFormatString(s string) (*dsl.Scope, []dsl.ASTNode, error) {
	toks, err := dsl.NewLexer(s).Lex()
	if err != nil {
		return nil, nil, err
	}
	parser := dsl.NewParser(toks)
	nodes, err := parser.Parse()
	return &dsl.Scope{Parser: parser}, nodes, err
}
