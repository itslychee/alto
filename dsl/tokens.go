package dsl

type TokenType int

const (
	VarNotation = iota + 1
	StringLiteral
	LParen
	RParen
	LCurlyBrace
	RCurlyBrace
	LArrow
	RArrow
	Separator
)

var tokens = map[rune]TokenType{
	'{': LCurlyBrace,
	'}': RCurlyBrace,
	'|': Separator,
	'(': LParen,
	')': RParen,
	'%': VarNotation,
	// '<': LArrow,
	// '>': RArrow,
}

type Token struct {
	Type     TokenType
	Value    string
	Position int
}
