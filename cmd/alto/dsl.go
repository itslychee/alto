package main

import (
	"regexp"
	"runtime"

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



var AltoFunctions = map[string]dsl.ASTFunction{
	"fn_clean":  func() FnCleanFunction {
		var reservedKeywords *regexp.Regexp
		if runtime.GOOS == "windows" {
			reservedKeywords = regexp.MustCompile(`[\pC"*/:<>?\\|]+`)
		} else {
			reservedKeywords = regexp.MustCompile(`[/\x{0}]+`)
		}
		return FnCleanFunction{regex: reservedKeywords}
	}(),
	"fp_unique": nil,
	"print": PrintFunction{},
	"skip": SkipFunc{},
}
