package dsl

type Type uint

const (
    _ Type = iota
    Comment
    GroupBegin
    GroupEnd
    Separator
    VariableBegin
    VariableEnd
    VariableNamespace
    FunctionBegin
    FunctionEnd
    String
)

func (t Type) String() string {
	switch t {
	case Comment:
		return "Comment"
	case GroupBegin:
		return "GroupBegin"
	case GroupEnd:
		return "GroupEnd"
	case Separator:
		return "Separator"
	case VariableBegin:
		return "VariableBegin"
	case VariableEnd:
		return "VariableEnd"
	case VariableNamespace:
		return "VariableNamespace"
	case FunctionBegin:
		return "FunctionBegin"
	case FunctionEnd:
		return "FunctionEnd"
	default:
		return "?unknown?"
	}
}

type Lexer struct {
    lexemes []*Token
    position int
    dsl     string
}

func (l Lexer) CalcLocation(position int) (location [2]uint) {
    for k := 0; k <= position; k++ {
        if l.dsl[k] == '\n' {
            location[0]++
            location[1] = 0
        } else {
            location[1]++
        }
    }
    return
}

func (l Lexer) DSL() string {
    return l.dsl
}

func (l *Lexer) NextToken() {
    token := &Token{
        Value: string(l.dsl[l.position]),
        Position: l.position,
        Location: [2][2]uint{
            l.CalcLocation(l.position),
            {0,0},
        },
    }
    switch l.dsl[l.position] {
    case '{':
        token.Type = GroupBegin
    case '}':
        token.Type = GroupEnd
    case '|':
        token.Type = Separator
    case '[':
        token.Type = VariableBegin
    case ':':
        token.Type = VariableNamespace
    case ']':
        token.Type = VariableEnd
    case '<':
        token.Type = FunctionBegin
    case '>':
        token.Type = FunctionEnd
    default:
        if l.position == 0 {
            token.Type = String
        } else {
            prevToken := l.lexemes[l.position-1]
            // Append to previous token of same type
            // instead of polluting the lexeme array with individual unicode characters
            if prevToken.Type == String {
                prevToken.Value += token.Value
                prevToken.Location[1] = token.Location[0]
                token = nil
            } else {
                token.Type = String
            }
        }
    }

}

type Token struct {
    Type     Type
    Value    string
    Position int
    Location [2][2]uint
}
