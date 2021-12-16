package dsl

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type functionDecl func([]ASTNode, *Scope) (string, error)

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
	"trim": WrapFunction(1, TrimFunc),
	"exit": WrapFunction(0, ExitFunc),
	"must": WrapFunction(-1, MustFunc),
	"eq":   DefaultConditional{Type: EqualTo},
	"neq":  DefaultConditional{Type: NotEqualTo},
	"gt":   DefaultConditional{Type: GreaterThan},
	"lt":   DefaultConditional{Type: LessThan},
	"gte":  DefaultConditional{Type: GreaterOrEqualTo},
	"lte":  DefaultConditional{Type: LessOrEqualTo},
	"fset": DefaultSetVariable{force: true},
	"set":  DefaultSetVariable{},
}

type DefaultSetVariable struct {
	force bool
}

func (t DefaultSetVariable) Execute(args []ASTNode, scope *Scope) (string, error) {
	key, err := args[1].Execute(scope)
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

type DefaultConditional struct {
	Type ConditionalType
}

func (t DefaultConditional) Execute(nodes []ASTNode, scope *Scope) (string, error) {
	cond1, err := nodes[1].Execute(scope)
	if err != nil {
		return "", err
	}
	cond2, err := nodes[2].Execute(scope)
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
		return nodes[3].Execute(scope)
	}
	return "", nil

}

func (t DefaultConditional) MaxParams() int {
	return 3
}

// FuncWrapper is NOT the same like ASTFunctionWrapper, this just provides
// a convienent wrapper for third-party functions
type FuncWrapper struct {
	function   functionDecl
	paramCount int
}

func (wrapper FuncWrapper) Execute(args []ASTNode, scope *Scope) (string, error) {
	return wrapper.function(args, scope)
}

func (wrapper FuncWrapper) MaxParams() int {
	return wrapper.paramCount
}

func MustFunc(args []ASTNode, scope *Scope) (string, error) {
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

func TrimFunc(args []ASTNode, scope *Scope) (string, error) {
	param, err := args[1].Execute(scope)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(param), nil
}

func ExitFunc([]ASTNode, *Scope) (string, error) {
	return "", errors.New("user called exit func")
}

func WrapFunction(paramCount int, function functionDecl) *FuncWrapper {
	return &FuncWrapper{
		function:   function,
		paramCount: paramCount,
	}
}
