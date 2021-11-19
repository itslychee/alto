package dsl

import (
	"strings"

	"github.com/pkg/errors"
)

var DefaultFunctions = map[string]ASTFunction{
	"trim": DefaultTrimSpace{},
	"must": DefaultMust{},
	"exit": DefaultExit{},
}

type DefaultTrimSpace struct {}
func (t DefaultTrimSpace) Execute(args []ASTField, scope *Scope) (string, error) {
	param, err := args[0].Execute(*scope)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(param), nil
}

func (t DefaultTrimSpace) MaxParams() int {
	return 1
}

type DefaultMust struct {}
func (t DefaultMust) Execute(args []ASTField, scope *Scope) (string, error) {
	var builder strings.Builder
	for i, v := range args {
		s, err := v.Execute(*scope)
		if err != nil {
			return "", errors.Wrapf(err, "must(): field at arg index '%d' returned an error", i)
		}
		builder.WriteString(s)
	}
	return builder.String(), nil
}

func (t DefaultMust) MaxParams() int {
	return -1
}

type DefaultExit struct {}

func (t DefaultExit) Execute([]ASTField, *Scope) (string, error) {
	return "", errors.New("exit()")
}

func (t DefaultExit) MaxParams() int {
	return 0
}