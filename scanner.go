package jsontoloyo

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"strings"
	"unicode/utf8"
)

var eof = rune(0)

type Scanner struct {
	r             *bufio.Reader
	buf           *bytes.Buffer
	lastReadToken string
	lastReadRune  rune
	lastReadItem  Item
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		r:   bufio.NewReader(r),
		buf: bytes.NewBufferString(""),
	}
}

// peek reads the next rune from the bufferred reader and then returns it.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) peek() rune {
	defer s.unread()
	return s.read()
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() {
	//log.Println("buf.UnreadRune (pre) is ", s.buf.String(), s.buf.Len())
	if err := s.buf.UnreadRune(); err == nil {
		//log.Println("buf.UnreadRune", string(s.lastReadRune))
		//log.Println("buf.UnreadRune is ", s.buf.String(), s.buf.Len())
		return
	} else if s.buf.Len() == 0 {
		//log.Println("buf.WriteRune", string(s.lastReadRune))
		s.buf.WriteRune(s.lastReadRune)
	} else {
		//log.Println("NewBuf", string(s.lastReadRune))
		// stuff in buffer and can't unread rune.
		newBuf := &bytes.Buffer{}
		newBuf.WriteRune(s.lastReadRune)
		newBuf.Write(s.buf.Bytes())
		s.buf = newBuf
	}
	// if err := s.r.UnreadRune(); err != nil {
	// 	log.Println("Error Unreading rune:", err)
	// 	panic(err)
	// }
}

// read reads the next rune from the bufferred reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) read() rune {
	if s.buf.Len() > 0 {
		r, _, err := s.buf.ReadRune()
		if err != nil {
			panic(err)
		}
		// log.Println("buf.ReadRune", string(r))
		s.lastReadRune = r
		return r
	}
	s.buf.Truncate(0) // need to get rid of back buffer
	r, _, err := s.r.ReadRune()
	//log.Println("s.r.ReadRune", string(r))
	s.lastReadRune = r
	if err == io.EOF {
		return eof
	}
	if err != nil {
		log.Println("Error reading rune:", err)
		panic(err)
	}
	return r
}

// read token attempts to read "s" from the buffer.
// If the next few bytes are not equal to s, all read bytes are unread and readToken returns false.
// otherwise readToken returns true and saves the value to lastReadToken.
func (s *Scanner) tryReadToken(token string) bool {
	//log.Println("tryReadToken", token)
	// if token == "Current_date" {
	// 	fmt.Println("Try read Current_date", token)
	// }
	readRunes := &bytes.Buffer{}
	for i := 0; i < len(token); i++ {
		r := s.read()

		if utf8.RuneLen(r) != 1 || strings.ToUpper(string(r)) != string(token[i]) {
			s.unread()
			if readRunes.Len() != 0 {
				// fmt.Println("UnreadStuff", readRunes.String())
				s.unreadString(readRunes.String())
			}
			return false
		}
		readRunes.WriteRune(r)
	}
	// check that we're at also a border character.
	r := s.peek()
	if len(token) == 1 {
		if token == readRunes.String() {
			s.lastReadToken = readRunes.String()
			return true
		}
	}
	if !isWhitespace(r) && !isParenthesis(r) && !(r == ',') && r != eof {
		s.unreadString(readRunes.String())
		return false
	}
	s.lastReadToken = readRunes.String()
	s.buf.Reset()
	return true
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}
func isParenthesis(ch rune) bool {
	return ch == '(' || ch == ')'
}

func isSquareBracket(ch rune) bool {
	return ch == '[' || ch == ']'
}
func isDot(ch rune) bool {
	return ch == '.'
}
func isDollar(ch rune) bool {
	return ch == '$'
}

func (s *Scanner) unreadString(str string) {
	//log.Println("unreadString", str)
	if s.buf.Len() == 0 {
		s.buf.WriteString(str)
	} else {
		// stuff in buffer and can't unread rune.
		newBuf := &bytes.Buffer{}
		newBuf.WriteString(str)
		newBuf.Write(s.buf.Bytes())
		s.buf = newBuf
	}
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanWhitespace() Item {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return Item{WHITESPACE, buf.String()}
}
func (s *Scanner) tryKeywords() bool {
	if s.tryReadToken("SUM") {
		s.lastReadItem = Item{SUM, s.lastReadToken}
		return true
	}
	if s.tryReadToken("AVG") {
		s.lastReadItem = Item{AVG, s.lastReadToken}
		return true
	}
	return false
}

// scanIdentifier consumes the current rune and all contiguous Identifier runes.
func (s *Scanner) scanIdentifier() Item {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent Identifier character into the buffer.
	// Non-Identifier characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof || ch == '.' {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '_' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return Item{IDENTIFIER, buf.String()}
}

// scanNumber consumes the current rune and all contiguous numeric runes.
func (s *Scanner) scanNumber() Item {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isDigit(ch) && ch != '.' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return Item{NUMBER, buf.String()}
}

func (s *Scanner) Scan() Item {
	ch := s.read()
	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	}

	if isDigit(ch) {
		s.unread()
		return s.scanNumber()
	}
	if isDot(ch) {
		//s.unread()
		return Item{DOT, string(ch)}
	}
	if ch == '(' {
		// s.unread()
		return Item{BRACKET_OPEN, string('(')}
	}
	if ch == ')' {
		// s.unread()
		return Item{BRACKET_OPEN, string('(')}
	}
	if isDollar(ch) {
		// s.unread()
		if s.tryKeywords() {
			return s.lastReadItem
		}
		// fmt.Println(string(s.buf.Bytes()))
		if s.peek() == '.' {
			s.buf.Reset()
			return Item{DOLLARSIGN, string(ch)}
		}
		if s.peek() == '[' {

			return Item{DOLLARSIGN, string(ch)}
		}

	}
	if isSquareBracket(ch) {

		if ch == '[' {
			// fmt.Println(s.buf.String())
			s.buf.Reset()
			return Item{SQUAREBRACKET_OPEN, string(ch)}
		} else {
			s.buf.Reset()
			return Item{SQUAREBRACKET_CLOSE, string(ch)}
		}

	}
	if isLetter(ch) {
		s.unread()

		return s.scanIdentifier()
	}
	switch ch {
	case eof:
		return Item{EOF, ""}
	}
	return Item{ILLEGAL, string(ch)}
}
