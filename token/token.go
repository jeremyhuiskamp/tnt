// Package token tokenizes a Typographical Number Theorem string.
//
// In addition to the mathematical symbols used in the book, several
// similar-looking characters that occur on a normal keyboard are accepted in
// order to make input easier.
package token

import "unicode"

type Token int

//go:generate stringer -type Token

const (
	ILLEGAL Token = iota
	EOF

	ZERO         // 0
	SUCCESSOR    // S*
	VARIABLE     // [a-e]'*
	OPEN_PAREN   // (
	CLOSE_PAREN  // )
	PLUS         // +
	MULTIPLY     // * . ·
	EQUALS       // =
	NEGATION     // ~
	OPEN_ANGLE   // <
	CLOSE_ANGLE  // >
	THERE_EXISTS // ∃ E
	FOR_ALL      // ∀ A
	COLON        // :
	AND          // ∧ ^
	OR           // ∨ V
	IF_THEN      // ⊃
)

type Scanner struct {
	src []rune
	pos int
}

func NewScanner(src string) *Scanner {
	return &Scanner{
		src: []rune(src),
		pos: 0,
	}
}

// Scan returns the next token in the expression.
//
// Both the token type and content are returned, but the content is
// only interesting for types VARIABLE and SUCCESSOR, which can have
// any number of different values.
//
// If an illegal token is encountered, no further progress is made and
// subsequent calls continue to return ILLEGAL.
//
// Once end of file is reached, EOF is returned for all subsequent calls.
//
// TODO: return position information for debugging.
func (s *Scanner) Scan() (Token, string) {
	for s.pos < len(s.src) && unicode.IsSpace(s.src[s.pos]) {
		s.pos++
	}

	if !(s.pos < len(s.src)) {
		return EOF, ""
	}

	ch := s.src[s.pos]
	switch ch {
	case '0':
		s.pos++
		return ZERO, "0"
	case '(':
		s.pos++
		return OPEN_PAREN, "("
	case ')':
		s.pos++
		return CLOSE_PAREN, ")"
	case '+':
		s.pos++
		return PLUS, "+"
	case '*', '.', '·':
		s.pos++
		return MULTIPLY, string(ch)
	case '=':
		s.pos++
		return EQUALS, "="
	case '~':
		s.pos++
		return NEGATION, "~"
	case '<':
		s.pos++
		return OPEN_ANGLE, "<"
	case '>':
		s.pos++
		return CLOSE_ANGLE, ">"
	case 'E', '∃':
		s.pos++
		return THERE_EXISTS, string(ch)
	case 'A', '∀':
		s.pos++
		return FOR_ALL, string(ch)
	case ':':
		s.pos++
		return COLON, ":"
	case '^', '∧':
		s.pos++
		return AND, string(ch)
	case 'V', '∨':
		s.pos++
		return OR, string(ch)
	case '⊃': // TODO: pick an easier-to-enter alternative character
		s.pos++
		return IF_THEN, string(ch)
	case 'a', 'b', 'c', 'd', 'e':
		variable := string(ch)
		s.pos++
		for s.pos < len(s.src) && s.src[s.pos] == '\'' {
			s.pos++
			variable += "'"
		}
		return VARIABLE, variable
	case 'S':
		successor := "S"
		s.pos++
		for s.pos < len(s.src) && s.src[s.pos] == 'S' {
			s.pos++
			successor += "S"
		}
		return SUCCESSOR, successor
	default:
		return ILLEGAL, string(ch)
	}
}
