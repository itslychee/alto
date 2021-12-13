package main

import (
	"errors"

	"github.com/ItsLychee/alto/dsl"
)

var (
	ErrSkip       = errors.New("requested skip")
	AltoFunctions = map[string]dsl.ASTFunction{
	    "uniqueFp": dsl.WrapFunction(2, uniqueFilepath),
		// <uniqueFp {%title%|%filename%} {%result% %index%}>

	}
)



func uniqueFilepath(nodes []dsl.ASTNode, scope *dsl.Scope) (string, error) {
	// <uniqueFp {<fset fp {%artist%/%title%.%filetype%}>|%fp%} {}>
	
}




func ParseFormatString(s string) (*dsl.Scope, []dsl.ASTNode, error) {
	toks, err := dsl.NewLexer(s).Lex()
	if err != nil {
		return nil, nil, err
	}
	parser := dsl.NewParser(toks)
	nodes, err := parser.Parse()
	return &dsl.Scope{Parser: parser}, nodes, err
}
