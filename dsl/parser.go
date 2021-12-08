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

type Scope struct {
	Variables map[string]string
	Functions map[string]ASTFunction
	Parser    *Parser
}

type Parser struct {
	toks         []*Token
	position     int
	CurrentToken *Token
	NextToken    *Token
	PrevToken    *Token
	arrowDepth   int
	groupDepth   int
}

func (p *Parser) UpdateCursor() error {
	p.PrevToken = new(Token)
	p.NextToken = new(Token)

	// set prevToken, if feasible
	if p.position > 0 && len(p.toks) > 0 {
		p.PrevToken = p.toks[p.position-1]
	}

	// set currentToken, if feasible
	if len(p.toks)-1 >= p.position {
		p.CurrentToken = p.toks[p.position]
	} else {
		return ErrNoMoreTokens
	}

	// set nextToken, if feasible
	if len(p.toks) != p.position+1 {
		p.NextToken = p.toks[p.position+1]
	}

	p.position++
	return nil
}

func (p *Parser) ParseNode() (ASTNode, error) {
	err := p.UpdateCursor()
	if err != nil {
		return nil, err
	}

	switch p.CurrentToken.Type {
	case VarNotation:
		if p.arrowDepth > 0 || p.groupDepth > 0 {
			err := p.UpdateCursor()
			if err != nil {
				return nil, fmt.Errorf("unterminated variable at pos %d", p.PrevToken.Position)
			}
			if p.NextToken.Type != VarNotation {
				return nil, fmt.Errorf("unterminated variable at pos %d", p.CurrentToken.Position)
			}

			if !IdentifierExpr.Match([]byte(p.CurrentToken.Value)) {
				return nil, fmt.Errorf("invalid variable name at pos %d", p.PrevToken.Position)
			}

			p.UpdateCursor()

			return ASTVariable{Name: p.PrevToken.Value}, nil
		}
		fallthrough
	case LArrow:
		p.arrowDepth++
		var wrapper ASTFunctionWrapper
		if err := p.UpdateCursor(); err != nil {
			return nil, fmt.Errorf("unterminated function at pos %d", p.CurrentToken.Position)
		}
		if p.CurrentToken.Type != StringLiteral {
			return nil, fmt.Errorf("invalid type after function group")
		}
		wrapper.Name = p.CurrentToken.Value
		if err := p.UpdateCursor(); err != nil || p.CurrentToken.Type != LParen {
			return nil, fmt.Errorf("unterminated function at pos %d", p.CurrentToken.Position)
		}

		var field ASTField
		for {
			switch p.NextToken.Type {
			case RParen:
				if len(field.Nodes) > 0 {
					wrapper.Args = append(wrapper.Args, field)
				}
				if err := p.UpdateCursor(); err != nil {
					return nil, fmt.Errorf("unterminated function at pos %d", p.CurrentToken.Position)
				}
				if err := p.UpdateCursor(); err != nil || p.CurrentToken.Type != RArrow {
					return nil, fmt.Errorf("unterminated function at pos %d", p.CurrentToken.Position)
				}
				return wrapper, nil
			case Separator:
				wrapper.Args = append(wrapper.Args, field)
				field = ASTField{}
			case RArrow:
				return nil, fmt.Errorf("unterminated function at pos %d", p.NextToken.Position)
			}
			n, err := p.ParseNode()
			if err != nil {
				return nil, err
			}
			field.Nodes = append(field.Nodes, n)

		}

	case StringLiteral, LParen, RParen:
		return ASTString{Value: p.CurrentToken.Value}, nil
	case LCurlyBrace:
		p.groupDepth++
		var group ASTGroup
		var field ASTField

		for {
			switch p.NextToken.Type {
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
