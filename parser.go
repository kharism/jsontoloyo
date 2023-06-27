package jsontoloyo

import (
	"errors"
	"io"
)

// parser parse the path traversing

type Token int

const (
	IDENTIFIER Token = iota
	NUMBER
	ILLEGAL
	DOLLARSIGN
	DOT
	SQUAREBRACKET_OPEN
	SQUAREBRACKET_CLOSE
	BRACKET_OPEN
	BRACKET_CLOSE
	WHITESPACE
	EOF

	// operation
	SUM
	AVG
)

var literals = []Token{IDENTIFIER, NUMBER}
var aggregate = []Token{SUM, AVG}

type Parser struct {
	s        *Scanner
	itemBuf  []Item
	lastItem Item
}

func isAggregate(t Token) bool {
	for _, v := range aggregate {
		if v == t {
			return true
		}
	}
	return false
}
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() Item {
	// If we have a token on the buffer, then return it.
	if len(p.itemBuf) > 0 {
		item := p.itemBuf[0]
		p.itemBuf = p.itemBuf[1:]
		return item
	}

	// Otherwise read the next token from the scanner.
	item := p.s.Scan()

	// Save it to the buffer in case we unscan later.
	p.lastItem = item

	return item
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() {
	p.itemBuf = append([]Item{p.lastItem}, p.itemBuf...)
}

// nextItem scans the next non-whitespace token.
func (p *Parser) nextItem() Item {
	item := p.scan()
	if item.Token == WHITESPACE {
		item = p.scan()
	}
	return item
}

// parse parameters for aggregate functions
// results are stored on result
func (p *Parser) ParseAggregateParam(result *[]*Selector) error {
	bracketCount := 0
	for {
		item := p.nextItem()
		switch item.Token {
		case BRACKET_OPEN:
			bracketCount++
		case BRACKET_CLOSE:
			bracketCount--
			if bracketCount == 0 {
				return nil
			}
		case DOLLARSIGN:
			// p.unscan()
			subStatement := []*Selector{}
			e := p.Parse(&subStatement)
			if e != nil {
				return e
			}
			*result = append(*result, subStatement...)
		default:
			return nil
		}
	}
	// return nil
}
func (p *Parser) Parse(result *[]*Selector) (errRet error) {
	// squareBrackedCount := 0
	// isASubQuery := false
	for {
		item := p.nextItem()
		// fmt.Println(item.Val)
		switch item.Token {
		case DOLLARSIGN:
			newSelector := &Selector{FieldName: "$", Token: item.Token, FunctionName: ""}
			*result = append(*result, newSelector)
			// fmt.Println("DOLLARDETECTED")
		case DOT:
			continue
		case IDENTIFIER:
			newSelector := &Selector{FieldName: item.Val, Token: item.Token, FunctionName: ""}
			*result = append(*result, newSelector)
		case SQUAREBRACKET_OPEN:
			newSelector := &Selector{FieldName: item.Val, Token: item.Token, FunctionName: ""}
			*result = append(*result, newSelector)
			// fmt.Println("SquareOpen")
			// isASubQuery = true
			// squareBrackedCount++
			k := p.nextItem()
			if k.Token == NUMBER {
				newSelector = &Selector{Number: k.Val, Token: k.Token}
				*result = append(*result, newSelector)
				k = p.nextItem()
				newSelector = &Selector{FieldName: k.Val, Token: k.Token}
				*result = append(*result, newSelector)

			} else if k.Token == SQUAREBRACKET_CLOSE {
				newSelector := &Selector{FieldName: k.Val, FunctionName: "", Token: k.Token}
				*result = append(*result, newSelector)
			} else {
				p.unscan()
			}
		case SQUAREBRACKET_CLOSE:
			return errors.New("Parsing error, found ] where it's not supposed to")
			// newSelector := &Selector{FieldName: item.Val, FunctionName: ""}
			// *result = append(*result, newSelector)
			// squareBrackedCount--
			// if squareBrackedCount == 0 {

			// }
		case AVG, SUM:
			newSelector := &Selector{FunctionName: "", Token: item.Token}
			switch item.Token {
			case AVG:
				newSelector.FunctionName = "AVG"
			case SUM:
				newSelector.FunctionName = "SUM"
			}
			//scan parameters for this aggregate
			parameters := []*Selector{}
			p.ParseAggregateParam(&parameters)
			newSelector.FunctionParam = parameters
			*result = append(*result, newSelector)
		case NUMBER:
			newSelector := &Selector{Number: item.Val, FunctionName: ""}
			*result = append(*result, newSelector)
		case EOF:
			return nil
		}
	}
}
