package tnt

import (
	"fmt"

	"github.com/jeremyhuiskamp/tnt/token"
)

// ParseFormula parses a complete TNT Formula.
func ParseFormula(src string) (Formula, error) {
	s := token.NewScanner(src)
	formula, err := parseFormula(s)
	if err != nil {
		return nil, err
	}

	tok, _ := s.Scan()
	if tok != token.EOF {
		return nil, fmt.Errorf("expected EOF but got %s", tok)
	}

	return formula, nil
}

// parseFormula parses a TNT Formula from a Scanner, but does not
// check for more Tokens after what it consumes.
func parseFormula(s *token.Scanner) (Formula, error) {
	tok, val := s.Scan()
	switch tok {
	case token.NEGATION:
		return parseNegation(s)
	case token.OPEN_ANGLE:
		return parseCompound(s)
	case token.FOR_ALL:
		return parseQuantification(FOR_ALL, s)
	case token.THERE_EXISTS:
		return parseQuantification(THERE_EXISTS, s)
	default:
		return parseAtom(tok, val, s)
	}
}

func parseAtom(tok token.Token, val string, s *token.Scanner) (Formula, error) {
	left, err := parseTermToken(tok, val, s)
	if err != nil {
		return nil, err
	}

	tok, _ = s.Scan()
	if tok != token.EQUALS {
		return nil, fmt.Errorf("expected = but got %s", tok)
	}

	right, err := parseTerm(s)
	if err != nil {
		return nil, err
	}

	return Atom{
		Left:  left,
		Right: right,
	}, nil
}

// parseTermToken parses a Term from the Scanner, assuming the first Token
// has already been scanned.
func parseTermToken(tok token.Token, val string, s *token.Scanner) (Term, error) {
	switch tok {
	case token.ZERO:
		return Numeral(0), nil
	case token.SUCCESSOR:
		term, err := parseTerm(s)
		if err != nil {
			return nil, err
		}
		switch term := term.(type) {
		case Numeral:
			return Numeral(len(val)) + term, nil
		default:
			return Successor{
				Quantity: len(val),
				Term:     term,
			}, nil
		}
	case token.VARIABLE:
		return Variable(val), nil
	case token.OPEN_PAREN:
		left, err := parseTerm(s)
		if err != nil {
			return nil, err
		}

		var kind CompoundTermKind
		tok, _ := s.Scan()
		switch tok {
		case token.PLUS:
			kind = PLUS
		case token.MULTIPLY:
			kind = MULTIPLY
		default:
			return nil, fmt.Errorf("expected + or * but got %s", tok)
		}

		right, err := parseTerm(s)
		if err != nil {
			return nil, err
		}

		tok, _ = s.Scan()
		if tok != token.CLOSE_PAREN {
			return nil, fmt.Errorf("expected ) but got %s", tok)
		}

		return CompoundTerm{
			Kind:  kind,
			Left:  left,
			Right: right,
		}, nil
	}
	return nil, fmt.Errorf("unexpected token in term: %s", tok)
}

// parseTerm parses a Term from the Scanner assuming none of the
// Tokens of the Term have yet been scanned.
func parseTerm(s *token.Scanner) (Term, error) {
	tok, val := s.Scan()
	return parseTermToken(tok, val, s)
}

func parseCompound(s *token.Scanner) (Formula, error) {
	left, err := parseFormula(s)
	if err != nil {
		return nil, err
	}

	var kind CompoundKind
	tok, _ := s.Scan()
	switch tok {
	case token.AND:
		kind = AND
	case token.OR:
		kind = OR
	case token.IF_THEN:
		kind = IF_THEN
	default:
		return nil, fmt.Errorf("expected AND, OR or IF_THEN in compound formula, "+
			"but got %s", tok)
	}

	right, err := parseFormula(s)
	if err != nil {
		return nil, err
	}

	tok, _ = s.Scan()
	if tok != token.CLOSE_ANGLE {
		return nil, fmt.Errorf("expected > but got %s", tok)
	}

	return Compound{
		Kind:  kind,
		Left:  left,
		Right: right,
	}, nil
}

func parseNegation(s *token.Scanner) (Formula, error) {
	formula, err := parseFormula(s)
	if err != nil {
		return nil, err
	}
	return Negation{formula}, nil
}

func parseQuantification(kind QuantificationKind, s *token.Scanner) (Formula, error) {
	tok, varName := s.Scan()
	if tok != token.VARIABLE {
		return nil, fmt.Errorf("expected VARIABLE but got %s", tok)
	}

	tok, _ = s.Scan()
	if tok != token.COLON {
		return nil, fmt.Errorf("expected : but got %s", tok)
	}

	formula, err := parseFormula(s)
	if err != nil {
		return nil, err
	}

	return Quantification{
		Kind:     kind,
		Variable: Variable(varName),
		Formula:  formula,
	}, nil
}
