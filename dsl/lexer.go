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
	token := new(Token)

	if t, ok := tokens[ch]; !ok {
		switch ch {
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
			token.Type = StringLiteral
		}
	} else {
		token.Type = t
	}
	token.Position = lex.position
	token.Value = string(ch)
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
