package tnt_test

import (
	"fmt"

	"github.com/jeremyhuiskamp/tnt"
)

func Example() {
	formula, _ := tnt.ParseFormula("<∀b:b=b∧~c=c>")
	fmt.Println("variables:", formula.Variables())
	fmt.Println("free variables:", formula.FreeVariables())
	fmt.Println("open:", formula.Open())
	fmt.Println("well formed:", formula.WellFormed())
	// Output: variables: [b c]
	// free variables: [c]
	// open: true
	// well formed: true
}
