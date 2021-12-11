package dsl

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type ConditionalType int

const (
	EqualTo ConditionalType = iota + 1
	GreaterThan
	LessThan
	NotEqualTo
	GreaterOrEqualTo
	LessOrEqualTo
)

var DefaultFunctions = map[string]ASTFunction{
	"trim": DefaultTrimSpace{},
	"must": DefaultMust{},
	"exit": DefaultExit{},
	"eq":   DefaultConditional{Type: EqualTo},
	"neq":  DefaultConditional{Type: NotEqualTo},
	"gt":   DefaultConditional{Type: GreaterThan},
	"lt":   DefaultConditional{Type: LessThan},
	"gte":  DefaultConditional{Type: GreaterOrEqualTo},
	"lte":  DefaultConditional{Type: LessOrEqualTo},
	"fset": DefaultSetVariable{force: true},
	"set":  DefaultSetVariable{},
}

type DefaultTrimSpace struct{}

func (t DefaultTrimSpace) Execute(args []ASTNode, scope *Scope) (string, error) {
	param, err := args[0].Execute(scope)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(param), nil
}

func (t DefaultTrimSpace) MaxParams() int {
	return 1
}

type DefaultMust struct{}

func (t DefaultMust) Execute(args []ASTNode, scope *Scope) (string, error) {
	var builder strings.Builder
	for i, v := range args {
		s, err := v.Execute(scope)
		if err != nil {
			return "", errors.Wrapf(err, "must: field at arg index '%d' returned an error", i)
		}
		if s == "" {
			return "", fmt.Errorf("must: field at arg index '%d' returned an empty response", i)
		}
		builder.WriteString(s)
	}
	return builder.String(), nil
}

func (t DefaultMust) MaxParams() int {
	return -1
}

type DefaultExit struct{}

func (t DefaultExit) Execute([]ASTNode, *Scope) (string, error) {
	return "", errors.New("user called exit()")
}

func (t DefaultExit) MaxParams() int {
	return 0
}

type DefaultConditional struct {
	Type ConditionalType
}

func (t DefaultConditional) Execute(nodes []ASTNode, scope *Scope) (string, error) {
	cond1, err := nodes[0].Execute(scope)
	if err != nil {
		return "", err
	}
	cond2, err := nodes[1].Execute(scope)
	if err != nil {
		return "", err
	}

	var passed bool
	switch t.Type {
	case EqualTo, NotEqualTo:
		if t.Type == EqualTo {
			passed = cond1 == cond2
		} else {
			passed = cond1 != cond2
		}
	default:
		cnd1, err := strconv.Atoi(cond1)
		if err != nil {
			return "", errors.Wrapf(err, "error while converting %s to integer internally", cond1)
		}
		cnd2, err := strconv.Atoi(cond2)
		if err != nil {
			return "", errors.Wrapf(err, "error while converting %s to integer internally", cond2)
		}

		switch t.Type {
		case GreaterThan:
			passed = cnd1 > cnd2
		case LessThan:
			passed = cnd1 < cnd2
		case GreaterOrEqualTo:
			passed = cnd1 >= cnd2
		case LessOrEqualTo:
			passed = cnd1 <= cnd2
		}
	}
	if passed {
		return nodes[2].Execute(scope)
	}
	return "", nil

}

func (t DefaultConditional) MaxParams() int {
	return 3
}

type DefaultSetVariable struct {
	force bool
}

func (t DefaultSetVariable) Execute(args []ASTNode, scope *Scope) (string, error) {
	key, err := args[0].Execute(scope)
	if err != nil {
		return "", err
	}

	if _, ok := scope.Variables[key]; ok && !t.force {
		return "", fmt.Errorf("variable with key '%s' already exists", key)
	}

	val, err := args[1].Execute(scope)
	if err != nil {
		return "", err
	}

	scope.Variables[key] = val
	return "", nil
}

func (t DefaultSetVariable) MaxParams() int {
	return 2
}
