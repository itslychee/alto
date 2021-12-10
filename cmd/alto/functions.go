package main

import (
	"errors"
	"log"
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
	// <func one two>
	// <func one two three>

	joiner := args[len(args)-1]
	args = args[:len(args)-1]
	join, err := joiner.Execute(scope)
	if err != nil {
		return "", err
	}

	var executedNodes []string
	for _, v := range args {
		s, err := v.Execute(scope)
		if err != nil {
			return "", err
		}
		if s == "" {
			continue
		}
		executedNodes = append(executedNodes, s)
	}
	return strings.Join(executedNodes, join), nil
}

func (f FnCleanFunction) MaxParams() int {
	return -1
}

type PrintFunction struct {}

func (p PrintFunction) Execute(args []dsl.ASTNode, scope *dsl.Scope) (string, error) {
	var builder strings.Builder
	for _, v := range args {
		s, err := v.Execute(scope)
		if err != nil {
			return "", err
		}
		builder.WriteString(s)
	}
	log.Println(builder.String())
	return "", nil
}

func (p PrintFunction) MaxParams() int {
	return -1
}

type SkipFunc struct {}

func (s SkipFunc) Execute([]dsl.ASTNode, *dsl.Scope) (string, error) {
	return "", errSkip
}
func (s SkipFunc) MaxParams() int {
	return 0
}