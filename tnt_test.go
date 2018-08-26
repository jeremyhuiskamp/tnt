package tnt

import (
	"testing"
)

func TestWellFormed(t *testing.T) {
	for _, wellFormed := range []string{
		"0=0",
		"~0=0",
		"<0=0^0=0>",
		"<a=b^b=c>",
		"~<a=b^b=c>",
		"Aa:Eb:<a=b^b=c>",
	} {
		formula, err := ParseFormula(wellFormed)
		if err != nil {
			t.Errorf("error parsing %q: %s", wellFormed, err)
		} else if !formula.WellFormed() {
			t.Errorf("expected %q to be well formed", wellFormed)
		}
	}

	for _, notWellFormed := range []string{
		"<Aa:0=0^a=a>",
		"~<Aa:0=0^a=a>",
		"Eb:<Aa:0=0^a=b>",
		"Eb:Ab:a=b", // b is not free
	} {
		formula, err := ParseFormula(notWellFormed)
		if err != nil {
			t.Errorf("error parsing %q: %s", notWellFormed, err)
		} else if formula.WellFormed() {
			t.Errorf("expected %q to be not well formed", notWellFormed)
		}
	}
}

func TestOpen(t *testing.T) {
	for _, open := range []string{
		"a=b",
		"(a+b)=0",
		"~(a+b)=0",
		"S(a+b)=0",
	} {
		formula, err := ParseFormula(open)
		if err != nil {
			t.Errorf("error parsing %q: %s", open, err)
		} else if !formula.Open() {
			t.Errorf("expected %q to be open", open)
		}
	}

	for _, closed := range []string{
		"S0=0",
		"Aa:Eb:(a+b)=0",
		"~Aa:Eb:(a+b)=0",
	} {
		formula, err := ParseFormula(closed)
		if err != nil {
			t.Errorf("error parsing %q: %s", closed, err)
		} else if formula.Open() {
			t.Errorf("expected %q to be closed", closed)
		}
	}
}
