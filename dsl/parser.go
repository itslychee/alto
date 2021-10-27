package dsl

import (
	"errors"
	"fmt"
	"regexp"
	// "strings"
)

var (
	ErrNoMoreTokens = errors.New("there are no more available tokens to parse, end of slice")
)

var IdentifierExpr = regexp.MustCompile(`\w+`)

type FunctionDecl func(args []ASTField, scope *Scope) (string, error)

type Scope struct {
	Variables map[string]string
	Functions map[string]FunctionDecl
	Parser    *Parser
}

type Parser struct {
	toks          []*Token
	position      int
	currentToken  *Token
	nextToken     *Token
	prevToken     *Token
	arrowDepth    int
	groupDepth    int
}

func (p *Parser) UpdateCursor() error {
	p.prevToken = new(Token)
	p.nextToken = new(Token)

	// set prevToken, if feasible
	if p.position > 0 && len(p.toks) > 0 {
		p.prevToken = p.toks[p.position-1]
	}

	// set currentToken, if feasible
	if len(p.toks)-1 >= p.position {
		p.currentToken = p.toks[p.position]
	} else {
		return ErrNoMoreTokens
	}

	// set nextToken, if feasible
	if len(p.toks) != p.position+1 {
		p.nextToken = p.toks[p.position+1]
	}

	p.position++
	return nil
}

func (p *Parser) ParseNode() (ASTNode, error) {
	err := p.UpdateCursor()
	if err != nil {
		return nil, err
	}

	switch p.currentToken.Type {
	case VarNotation:
		if p.arrowDepth > 0 || p.groupDepth > 0 {
			err := p.UpdateCursor()
			if err != nil {
				return nil, fmt.Errorf("unterminated variable at pos %d", p.prevToken.Position)
			}
			if p.nextToken.Type != VarNotation {
				return nil, fmt.Errorf("unterminated variable at pos %d", p.currentToken.Position)
			}

			if !IdentifierExpr.Match([]byte(p.currentToken.Value)) {
				return nil, fmt.Errorf("invalid variable name at pos %d", p.prevToken.Position)
			}

			p.UpdateCursor()

			return ASTVariable{name: p.prevToken.Value}, nil
		}
		fallthrough
	case StringLiteral:
		return ASTString{Value: p.currentToken.Value}, nil

	case LCurlyBrace:
		p.groupDepth++
		var group ASTGroup
		var field ASTField

		for {
			switch p.nextToken.Type {
			case RCurlyBrace:
				p.groupDepth--
				p.UpdateCursor()
				group.Fields = append(group.Fields, field)
				return group, nil

			case Separator:
				p.UpdateCursor()
				group.Fields = append(group.Fields, field)
				field.Nodes = []ASTNode{}
				continue
			}

			n, err := p.ParseNode()
			if err != nil {
				if err == ErrNoMoreTokens {
					return nil, fmt.Errorf("unterminated group")

				}
				return nil, err
			}
			field.Nodes = append(field.Nodes, n)
		}

	default:
		return nil, ErrNoMoreTokens
	}
}

func (p *Parser) Parse() ([]ASTNode, error) {
	var nodes []ASTNode
	for {
		n, err := p.ParseNode()
		if n != nil {
			nodes = append(nodes, n)
		}
		if err != nil {
			if err != ErrNoMoreTokens {
				return nil, err
			}
			return nodes, nil
		}
	}
}

func NewParser(toks []*Token) *Parser {
	return &Parser{toks: toks}
}
