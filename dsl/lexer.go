package dsl

import (
	"errors"
	"io"
	"strings"
)

type Lexer struct {
	position int
	reader   *strings.Reader
}

func (lex *Lexer) ReadToken() (*Token, error) {
	ch, _, err := lex.reader.ReadRune()
	if err != nil {
		if err == io.EOF {
			return nil, nil
		}
		panic(err)
	}
	lex.position++

	var token = &Token{Value: string(ch), Position: lex.position}

	switch ch {
	case '{':
		token.Type = LCurlyBrace
	case '}':
		token.Type = RCurlyBrace
	case '|':
		token.Type = Separator
	case '%':
		token.Type = VarNotation
	case '<':
		token.Type = LArrow
	case '>':
		token.Type = RArrow
	case '(':
		token.Type = LParen
	case ')': 
		token.Type = RParen
	case '\\':
		ch, _, err = lex.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return nil, errors.New("unsatisfied escape sequence")
			}
			panic(err)
		}
		lex.position++
		fallthrough
	default:
		token.Value = string(ch)
		token.Position = lex.position
		token.Type = StringLiteral

	}

	return token, nil
}

func (lex *Lexer) Lex() (toks []*Token, err error) {
	for {
		currentToken, er := lex.ReadToken()
		if er != nil || currentToken == nil {
			err = er
			return
		}

		if currentToken.Type == StringLiteral && len(toks) >= 1 {
			if prevToken := toks[len(toks)-1]; prevToken.Type == StringLiteral {
				prevToken.Value = strings.Join([]string{prevToken.Value, currentToken.Value}, "")
				continue
			}
		}
		toks = append(toks, currentToken)
	}
}

func NewLexer(s string) *Lexer {
	return &Lexer{
		reader: strings.NewReader(s),
	}
}
