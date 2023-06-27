package jsontoloyo

type Item struct {
	Token Token
	Val   string
}

func (i Item) String() string {
	return i.Inspect()
}

func (i Item) Inspect() string {
	switch i.Token {
	case IDENTIFIER:
		return "Identifier"
	case NUMBER:
		return "Number"
	case DOLLARSIGN:
		return "Dollarsign"
	case SQUAREBRACKET_OPEN:
		return "Squarebracket_open"
	case SQUAREBRACKET_CLOSE:
		return "Squarebracket_close"
	case ILLEGAL:
		return "Illegal"
	case DOT:
		return "Dot"
	default:
		return "UNKNOWN"
	}
}
