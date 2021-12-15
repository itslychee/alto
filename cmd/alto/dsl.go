package main

import (
	"errors"
	"os"
	"strconv"

	"github.com/ItsLychee/alto/dsl"
)

var (
	ErrSkip       = errors.New("requested skip")
	AltoFunctions = map[string]dsl.ASTFunction{
		"uniqueFp": dsl.WrapFunction(2, uniqueFilepath),
		"exists":   dsl.WrapFunction(1, exists),
	}
)

func uniqueFilepath(nodes []dsl.ASTNode, scope *dsl.Scope) (string, error) {
	// To prevent any unwanted behavior such as overwriting variables,
	// uniqueFp will just copy this scope.
	copiedScope := *scope
	for i := 1; ; i++ {
		s, err := nodes[1].Execute(&copiedScope)
		if err != nil {
			return s, err
		}
		// We only care if the actual filepath is available, no need for the
		// system to see the contents of what a possible symlink is pointing to
		_, err = os.Lstat(s)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return s, nil
			}
			return "", err
		}
		copiedScope.Variables["index"] = strconv.Itoa(i)
	}
}

func exists(nodes []dsl.ASTNode, scope *dsl.Scope) (string, error) {
	path, err := nodes[1].Execute(scope)
	if err != nil {
		return path, err
	}
	_, err = os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return path, err
	}
	return "", err
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
