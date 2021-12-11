package main

import (
	"errors"
	"regexp"
	"strings"

	"github.com/ItsLychee/alto/dsl"
)

var errSkip = errors.New("requested skip")

type FnCleanFunction struct {
	regex *regexp.Regexp
}

func (f FnCleanFunction) Execute(args []dsl.ASTNode, scope *dsl.Scope) (string, error) {
	if len(args) == 1 {
		s, err := args[0].Execute(scope)
		return f.regex.ReplaceAllString(s, "-"), err
	}

	joiner, err := args[len(args)-1].Execute(scope)
	if err != nil {
		return "", err
	}

	var executedNodes []string
	for _, v := range args[:len(args)-1] {
		s, err := v.Execute(scope)
		if err != nil {
			return "", err
		}
		if s == "" {
			continue
		}
		executedNodes = append(executedNodes, f.regex.ReplaceAllString(s, "-"))
	}
	return strings.Join(executedNodes, joiner), nil
}

func (f FnCleanFunction) MaxParams() int {
	return -1
}
