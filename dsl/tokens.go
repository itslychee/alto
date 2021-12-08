package dsl

type TokenType int

const (
	VarNotation = iota + 1
	StringLiteral
	LCurlyBrace
	RCurlyBrace
	LParen
	RParen
	LArrow
	RArrow
	Separator
)

type Token struct {
	Type     TokenType
	Value    string
	Position int
}
