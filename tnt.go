/*
Package tnt provides a representation of a Formula according to the
Typographical Number Theory presented in Chapter 8 of the book "Gödel, Escher,
Bach" by Douglas Hofstadter.

A TNT Formula follows roughly this format:

	term           := NUMERAL | VARIABLE | (term + term) | (term * term) | SUCCESSOR term
	atom           := term = term
	negation       := ~ formula
	compound       := < formula ( AND | OR ) formula >
	quantification := (THERE_EXISTS | FOR_ALL) VARIABLE : formula
	formula        := atom | negation | compound | quantification
*/
package tnt

// Term is a Numeral, Variable, Successor or CompoundTerm.
type Term interface {
	Variables() VariableSet
}

// Formula is an Atom, Negation, Compound or Quantification.
type Formula interface {
	Variables() VariableSet
	FreeVariables() VariableSet
	Open() bool
	WellFormed() bool
}

// Numeral is a Term of the form 0, S0, SS0, etc.
type Numeral int

// Variables returns an empty set, since a Numeral has no Variables.
func (n Numeral) Variables() VariableSet {
	return nil
}

// Variable is a Term of the form a, b, c, d, e, a', b', etc.
type Variable string

// Variables returns the set including only this Variable.
func (v Variable) Variables() VariableSet {
	return NewVariableSet(v)
}

// Successor is a Term of the form S*x where x is a Term.
type Successor struct {
	Quantity int
	Term     Term
}

// Variables returns the Variables of the Successor's Term.
func (s Successor) Variables() VariableSet {
	return s.Term.Variables()
}

// CompoundTermKind is either + or *.
type CompoundTermKind int

//go:generate stringer -type CompoundTermKind

const (
	PLUS CompoundTermKind = iota
	MULTIPLY
)

// CompoundTerm is a Term in the form (x+y) or (x*y) where x and y are Terms.
type CompoundTerm struct {
	Kind  CompoundTermKind
	Left  Term
	Right Term
}

// Variables returns the union of the Variables of the two
// contained Terms.
func (c CompoundTerm) Variables() VariableSet {
	return c.Left.Variables().Union(c.Right.Variables())
}

// Atom is a Formula in the form x=y, where x and y are Terms.
type Atom struct {
	Left  Term
	Right Term
}

// Variables returns the union of the Variables of the two
// contained Terms.
func (a Atom) Variables() VariableSet {
	return a.Left.Variables().Union(a.Right.Variables())
}

// FreeVariables returns the same as Variables since Atoms
// cannot contain quantifications.
func (a Atom) FreeVariables() VariableSet {
	return a.Variables()
}

// Open returns true if the atom contains any variables,
// since all variables in an atom are free by definition.
func (a Atom) Open() bool {
	return len(a.Variables()) != 0
}

// WellFormed always returns true.
func (a Atom) WellFormed() bool {
	return true
}

// Negation is a Formula that is the negation of another Formula.
type Negation struct {
	Formula Formula
}

// Variables returns the same value as the contained Formula.
func (n Negation) Variables() VariableSet {
	return n.Formula.Variables()
}

// FreeVariables returns the same value as the contained Formula.
func (n Negation) FreeVariables() VariableSet {
	return n.Formula.FreeVariables()
}

// Open returns the same value as the contained Formula.
func (n Negation) Open() bool {
	return n.Formula.Open()
}

// WellFormed returns the same value as the contained Formula.
func (n Negation) WellFormed() bool {
	return n.Formula.WellFormed()
}

// CompoundKind is "and", "or" or "if, then"
type CompoundKind int

//go:generate stringer -type CompoundKind

const (
	AND CompoundKind = iota
	OR
	IF_THEN
)

// Compound is a Formula of the form <x or y>, <x and y> or <if x then y> where
// x and y are well-formed Formulas.
type Compound struct {
	Kind  CompoundKind
	Left  Formula
	Right Formula
}

// Variables is the union of the Variables of the two contained
// Formulas.
func (c Compound) Variables() VariableSet {
	return c.Left.Variables().Union(c.Right.Variables())
}

// FreeVariables is the union of the FreeVariables of the two
// contained Formulas.  This only makes sense if this Compound
// is WellFormed.
func (c Compound) FreeVariables() VariableSet {
	return c.Left.FreeVariables().Union(c.Right.FreeVariables())
}

// Open returns true if either contained Formula is Open.
func (c Compound) Open() bool {
	return c.Left.Open() || c.Right.Open()
}

// WellFormed returns true if both contained Formulas are WellFormed,
// and no Variable with is free in one is quantified in the other.
func (c Compound) WellFormed() bool {
	if !c.Left.WellFormed() || !c.Right.WellFormed() {
		return false
	}

	lv := c.Left.Variables()
	lf := c.Left.FreeVariables()
	lq := lv.Complement(lf)

	rv := c.Right.Variables()
	rf := c.Right.FreeVariables()
	rq := rv.Complement(rv)

	lfrq := lf.Intersection(rq)
	rflq := rf.Intersection(lq)

	// This is rather complicated to follow.  It might be
	// good to implement an error message stating which
	// variables in which formula are in violation.
	return len(lfrq) == 0 && len(rflq) == 0
}

// QuantificationKind is "there exists" or "for all".
type QuantificationKind int

//go:generate stringer -type QuantificationKind

const (
	THERE_EXISTS QuantificationKind = iota
	FOR_ALL
)

// Quantification is a Formula that quantifies one Variable of another Formula.
type Quantification struct {
	Kind     QuantificationKind
	Variable Variable
	Formula  Formula
}

// Variables returns the same value as the contained Formula.
func (q Quantification) Variables() VariableSet {
	return q.Formula.Variables()
}

// FreeVariables returns the same value as the contained Formula,
// but without the Variable quantified by this Quantification.
func (q Quantification) FreeVariables() VariableSet {
	return q.Formula.Variables().Complement(NewVariableSet(q.Variable))
}

// Open returns true if this Quantification has no FreeVariables.
func (q Quantification) Open() bool {
	return len(q.FreeVariables()) == 0
}

// WellFormed returns true if the contained Formula is WellFormed
// and the Variable quantified is free in the contained Formula.
func (q Quantification) WellFormed() bool {
	if !q.Formula.WellFormed() {
		return false
	}
	_, ok := q.Formula.FreeVariables()[q.Variable]
	return ok
}
