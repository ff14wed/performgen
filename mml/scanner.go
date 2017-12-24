package mml

import (
	"bufio"
	"bytes"
	"io"
)

// TokenType is defined to specifically talk about token types rather than ints
type TokenType int

// These constants define the different possible token types
const (
	TNote TokenType = iota
	TRest
	TTempo
	TLength
	TOctave
	TOctaveUp
	TOctaveDown
	TVolume
	TExtend
	TDot
	TModifier
	TNumeric
	TEOF
	TIllegal
)

// Token defines the type for a lexical token
type Token struct {
	typ   TokenType
	ident string
	pos   Position
}

// Type returns the type of the token
func (t Token) Type() TokenType {
	return t.typ
}

// Ident returns the identifier of the token
func (t Token) Ident() string {
	return t.ident
}

// Position returns the position in the input string of the token
func (t Token) Position() Position {
	return t.pos
}

func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\n' }

func isNumeric(ch rune) bool { return (ch >= '0' && ch <= '9') }

var eof = rune(0)

// Scanner represents a lexical scanner.
type Scanner struct {
	r          *bufio.Reader
	lineNum    int
	colNum     int
	prevColNum int
	reachedEOF bool
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

// read reads the next rune from the bufferred reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) read() rune {
	if s.lineNum == 0 {
		s.lineNum = 1
	}
	ch, _, err := s.r.ReadRune()
	if err != nil {
		if !s.reachedEOF {
			s.colNum++
			s.reachedEOF = true
		}
		return eof
	}
	s.colNum++
	if ch == '\n' {
		s.prevColNum = s.colNum
		s.colNum = 0
		s.lineNum++
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() {
	_ = s.r.UnreadRune()
	s.colNum--
	if s.colNum < 0 {
		s.lineNum--
	}
}

// Scan returns the next token and literal value.
func (s *Scanner) Scan() Token {
	// Read the next rune.
	ch := s.read()

	// Eat all whitespace
	if isWhitespace(ch) {
		s.eatWhitespace()
		ch = s.read()
	}

	// If we see a numeric then consume as an ident.
	if isNumeric(ch) {
		s.unread()
		return s.scanNumeric(s.lineNum, s.colNum)
	}

	// Otherwise read the individual character.
	switch ch {
	case eof:
		return s.buildToken(TEOF, string(ch))
	case 'r', 'R':
		return s.buildToken(TRest, string(ch))
	case 't', 'T':
		return s.buildToken(TTempo, string(ch))
	case 'l', 'L':
		return s.buildToken(TLength, string(ch))
	case '#', '+', '-':
		return s.buildToken(TModifier, string(ch))
	case 'o', 'O':
		return s.buildToken(TOctave, string(ch))
	case '>':
		return s.buildToken(TOctaveUp, string(ch))
	case '<':
		return s.buildToken(TOctaveDown, string(ch))
	case 'V', 'v':
		return s.buildToken(TVolume, string(ch))
	case '&':
		return s.buildToken(TExtend, string(ch))
	case '.':
		return s.buildToken(TDot, string(ch))
	default:
		if (ch >= 'a' && ch <= 'g') || (ch >= 'A' && ch <= 'G') {
			return s.buildToken(TNote, string(ch))
		}
	}

	return s.buildToken(TIllegal, string(ch))
}

func (s *Scanner) buildToken(typ TokenType, ch string) Token {
	return Token{
		typ:   typ,
		ident: ch,
		pos:   Position{s.lineNum, s.colNum},
	}
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Scanner) eatWhitespace() {
	// Continuously read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		}
	}

	return
}

// scanNumeric consumes the current rune and all contiguous numeric runes.
func (s *Scanner) scanNumeric(lineNum, colNum int) Token {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	_, _ = buf.WriteRune(s.read())

	tok := s.buildToken(TNumeric, "")

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isNumeric(ch) {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	tok.ident = buf.String()
	// Otherwise return as a regular identifier.
	return tok
}
