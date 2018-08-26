package token

import (
	"reflect"
	"testing"
)

func TestScanner(t *testing.T) {
	type token struct {
		Token Token
		Value string
	}

	type testCase struct {
		Input    string
		Expected []token
	}

	for name, test := range map[string]testCase{
		"empty string": testCase{
			Input:    "",
			Expected: nil,
		},
		"whitespace string": testCase{
			Input:    " \t\n\r",
			Expected: nil,
		},
		"illegal token": testCase{
			Input: "_",
			Expected: []token{
				{ILLEGAL, "_"},
			},
		},
		"illegal token surrounded by whitespace": testCase{
			Input: " _ ",
			Expected: []token{
				{ILLEGAL, "_"},
			},
		},
		"zero": testCase{
			Input: " 0 ",
			Expected: []token{
				{ZERO, "0"},
			},
		},
		"zero then illegal": testCase{
			Input: "0_000",
			Expected: []token{
				{ZERO, "0"},
				{ILLEGAL, "_"},
			},
		},
		"all single characters": testCase{
			Input: "0 ( ) + * . · = ~ < > E ∃ A ∀ : ^ ∧ V ∨ ⊃",
			Expected: []token{
				{ZERO, "0"},
				{OPEN_PAREN, "("},
				{CLOSE_PAREN, ")"},
				{PLUS, "+"},
				{MULTIPLY, "*"},
				{MULTIPLY, "."},
				{MULTIPLY, "·"},
				{EQUALS, "="},
				{NEGATION, "~"},
				{OPEN_ANGLE, "<"},
				{CLOSE_ANGLE, ">"},
				{THERE_EXISTS, "E"},
				{THERE_EXISTS, "∃"},
				{FOR_ALL, "A"},
				{FOR_ALL, "∀"},
				{COLON, ":"},
				{AND, "^"},
				{AND, "∧"},
				{OR, "V"},
				{OR, "∨"},
				{IF_THEN, "⊃"},
			},
		},
		"variables": testCase{
			Input: "a bc d e a' b'e' a''''''",
			Expected: []token{
				{VARIABLE, "a"},
				{VARIABLE, "b"},
				{VARIABLE, "c"},
				{VARIABLE, "d"},
				{VARIABLE, "e"},
				{VARIABLE, "a'"},
				{VARIABLE, "b'"},
				{VARIABLE, "e'"},
				{VARIABLE, "a''''''"},
			},
		},
		"successors": testCase{
			Input: "S0SSS0Se",
			Expected: []token{
				{SUCCESSOR, "S"},
				{ZERO, "0"},
				{SUCCESSOR, "SSS"},
				{ZERO, "0"},
				{SUCCESSOR, "S"},
				{VARIABLE, "e"},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			s := NewScanner(test.Input)
			var got []token

			for {
				tok, value := s.Scan()
				if tok == EOF {
					break
				}
				got = append(got, token{
					Token: tok,
					Value: value,
				})
				if tok == ILLEGAL {
					break
				}
			}

			if !reflect.DeepEqual(got, test.Expected) {
				t.Fatalf("for input %q, expected %+v but got %+v",
					test.Input, test.Expected, got)
			}

			tok, value := s.Scan()
			illegal := len(got) > 0 && got[len(got)-1].Token == ILLEGAL
			if illegal {
				if tok != ILLEGAL {
					t.Fatalf("expected scanner to get stuck on ILLEGAL, "+
						"but got %s", tok)
				}
				prevValue := got[len(got)-1].Value
				if value != prevValue {
					t.Fatalf("expected scanner stuck on ILLEGAL to "+
						"return the same value %q but got %q",
						prevValue, value)
				}
			} else {
				if tok != EOF {
					t.Fatalf("expected scanner to get stuck on EOF, "+
						"but got %s", tok)
				}
			}
		})
	}
}
