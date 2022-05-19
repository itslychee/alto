package dsl

// Scope represents the node's execution scope
type Scope struct {
	// Variables are like functions, but passing parameters is forbidden.
	Variables map[string]map[string]ASTNode
	// Functions serve as a pipeline to Go from alto
	Functions map[string]map[string]ASTNode
}

// <strings:split {one two three} {four five 6}
type ASTNode struct {
	ASTExecBase
	// Lexemes contains the tokens that were used to build the node
	Lexemes []*Token
}

type ASTExecBase interface {
	// Validate can be used for validating the passed parameters and their types
	Validate(node ASTNode, scope *Scope, parameters []ASTNode) (err error)
	// Execute contains the implementation of the node
	Execute(node ASTNode, scope *Scope, parameters []ASTNode) (val any, err error)
}
