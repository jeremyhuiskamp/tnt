package tnt

import (
	"reflect"
	"testing"
)

func TestParseFormula(t *testing.T) {
	type testCase struct {
		Input    string
		Expected Formula
	}

	for name, test := range map[string]testCase{
		"simplest atom": {
			Input: "0=0",
			Expected: Atom{
				Left:  Numeral(0),
				Right: Numeral(0),
			},
		},
		"variables": {
			Input: "a=b",
			Expected: Atom{
				Left:  Variable("a"),
				Right: Variable("b"),
			},
		},
		"compound terms": {
			Input: "(0+a)=(b*(c+d))",
			Expected: Atom{
				Left: CompoundTerm{
					Kind:  PLUS,
					Left:  Numeral(0),
					Right: Variable("a"),
				},
				Right: CompoundTerm{
					Kind: MULTIPLY,
					Left: Variable("b"),
					Right: CompoundTerm{
						Kind:  PLUS,
						Left:  Variable("c"),
						Right: Variable("d"),
					},
				},
			},
		},
		"successors": {
			Input: "Sa=SS(SSS0+b)",
			Expected: Atom{
				Left: Successor{
					Quantity: 1,
					Term:     Variable("a"),
				},
				Right: Successor{
					Quantity: 2,
					Term: CompoundTerm{
						Kind:  PLUS,
						Left:  Numeral(3),
						Right: Variable("b"),
					},
				},
			},
		},
		"negation": {
			Input: "~0=S0",
			Expected: Negation{
				Atom{
					Left:  Numeral(0),
					Right: Numeral(1),
				},
			},
		},
		"double negation": {
			Input: "~~0=S0",
			Expected: Negation{
				Negation{
					Atom{
						Left:  Numeral(0),
						Right: Numeral(1),
					},
				},
			},
		},
		"compound": {
			Input: "<<0=0 ∧ a=b> ∨ <a=b ⊃ a=b>>",
			Expected: Compound{
				Kind: OR,
				Left: Compound{
					Kind: AND,
					Left: Atom{
						Left:  Numeral(0),
						Right: Numeral(0),
					},
					Right: Atom{
						Left:  Variable("a"),
						Right: Variable("b"),
					},
				},
				Right: Compound{
					Kind: IF_THEN,
					Left: Atom{
						Left:  Variable("a"),
						Right: Variable("b"),
					},
					Right: Atom{
						Left:  Variable("a"),
						Right: Variable("b"),
					},
				},
			},
		},
		"quantification": {
			Input: "Aa:Eb:0=0",
			Expected: Quantification{
				Kind:     FOR_ALL,
				Variable: Variable("a"),
				Formula: Quantification{
					Kind:     THERE_EXISTS,
					Variable: Variable("b"),
					Formula: Atom{
						Left:  Numeral(0),
						Right: Numeral(0),
					},
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			formula, err := ParseFormula(test.Input)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(formula, test.Expected) {
				t.Fatalf("expected %+v, got %+v",
					test.Expected, formula)
			}
		})
	}
}

func TestParseInvalidFormula(t *testing.T) {
	// TODO: add assertions about the meaning of the error
	// eg, position, expected token, actual token,
	// and perhaps what we thought we were working on
	for _, badFormula := range []string{
		"",
		"a",
		"0=0 a",
		"a+b=c+d",
		"(a-b)=0",
		"0=(a-b)",
		"S(a-b)=0",
		"((a-b)+c)=0",
		"(a+(b-c))=0",
		"(a+b=0)",
		"<(a)^0=0>",
		"<0=0_0=0>",
		"<0=0^(a)>",
		"<0=0^0=0 b>",
		"~(a)=0",
		"A0:S0=0",
		"Aa_0=0",
		"Aa:(a)=b",
	} {
		_, err := ParseFormula(badFormula)
		if err == nil {
			t.Fatalf("expected error for formula: %q", badFormula)
		}
	}
}
