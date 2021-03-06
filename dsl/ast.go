package dsl

import (
	"fmt"
	"strings"
)

type ASTNode interface {
	Execute(*Scope) (string, error)
}

type ASTField struct {
	Nodes []ASTNode
}

func (ast ASTField) Execute(scope *Scope) (string, error) {
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

func (ast ASTGroup) Execute(scope *Scope) (string, error) {
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

func (ast ASTString) Execute(_ *Scope) (string, error) {
	return ast.Value, nil
}

type ASTVariable struct {
	Name string
}

func (ast ASTVariable) Execute(scope *Scope) (string, error) {
	return scope.Variables[ast.Name], nil
}

// Due to the nature of the DSL requiring you to pass a Scope,
// which houses both functions and variales. We must use a "wrapper"
// that will provide the AST interface for functions.
// The disadvantage of this is that invalid function names or too many/little params
// cannot be detected at AST compile time.
type ASTFunctionWrapper struct {
	Name string
	Args []ASTNode
}

func (ast ASTFunctionWrapper) Execute(scope *Scope) (string, error) {
	v, ok := scope.Functions[ast.Name]
	if !ok {
		return "", fmt.Errorf("function \"%s\" does not exist", ast.Name)
	}
	ast.Args = append([]ASTNode{ASTString{Value: ast.Name}}, ast.Args...)

	if v.MaxParams() != -1 {
		if len(ast.Args)-1 > v.MaxParams() {
			return "", fmt.Errorf("too many arguments passed to function '%s'", ast.Name)
		}
		if len(ast.Args)-1 < v.MaxParams() {
			return "", fmt.Errorf("too few arguments passed to function '%s'", ast.Name)
		}
	}
	return v.Execute(ast.Args, scope)
}

type ASTFunction interface {
	Execute(args []ASTNode, scope *Scope) (string, error)
	MaxParams() int
}
