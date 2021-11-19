package dsl

import (
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
	name string
}

func (ast ASTVariable) Execute(scope Scope) (string, error) {
	return scope.Variables[ast.name], nil
}

type ASTFunction interface {
	Execute(args []ASTField, scope *Scope) (string, error)
	MaxParams() int
}