package dsl

import (
	"fmt"
	"strings"
)

type ASTNode interface {
	Execute(scope Scope) (string, error)
}

type ASTField struct {
	Nodes []ASTNode
}

func (ast ASTField) Execute(scope Scope) (string, error) {
	var builder strings.Builder
	for _, val := range ast.Nodes {
		s, err := val.Execute(scope)
		if err != nil {
			return "", err
		}
		if _, ok := val.(ASTVariable); ok && len(s) == 0 {
			return "", nil
		} else {
			builder.WriteString(s)
		}
	}
	return builder.String(), nil
}

type ASTGroup struct {
	Fields []ASTField
}

func (ast ASTGroup) Execute(scope Scope) (string, error) {
	for _, val := range ast.Fields {
		s, err := val.Execute(scope)
		if err != nil {
			return "", err
		}
		if len(s) != 0 {
			return s, nil
		}
	}
	return "", nil
}

type ASTString struct {
	Value string
}

func (ast ASTString) Execute(_ Scope) (string, error) {
	return ast.Value, nil
}

type ASTVariable struct {
	Name string
}

func (ast ASTVariable) Execute(scope Scope) (string, error) {
	return scope.Variables[ast.Name], nil
}

type ASTFunctionWrapper struct {
	Name     string
	Function ASTFunction
	Args     []ASTField
}

func (ast ASTFunctionWrapper) Execute(scope *Scope) (string, error) {
	if len(ast.Args) != -1 && len(ast.Args) != ast.Function.MaxParams() {
		return "", fmt.Errorf("too many arguments passed to function '%s'", ast.Name)
	}
	return ast.Function.Execute(ast.Args, scope)
}

type ASTFunction interface {
	Execute(args []ASTField, scope *Scope) (string, error)
	MaxParams() int
}
