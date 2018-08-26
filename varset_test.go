package tnt

import (
	"reflect"
	"testing"
)

func TestVariableSet(t *testing.T) {
	t.Parallel()

	type testCase struct {
		v1, v2                  VariableSet
		union, complement       VariableSet
		intersection, symmetric VariableSet
	}
	for i, test := range []testCase{
		{
			v1:           NewVariableSetString("a"),
			v2:           NewVariableSetString(),
			union:        NewVariableSetString("a"),
			complement:   NewVariableSetString("a"),
			intersection: NewVariableSetString(),
			symmetric:    NewVariableSetString("a"),
		},
		{
			v1:           NewVariableSetString(),
			v2:           NewVariableSetString("a"),
			union:        NewVariableSetString("a"),
			complement:   NewVariableSetString(),
			intersection: NewVariableSetString(),
			symmetric:    NewVariableSetString("a"),
		},
		{
			v1:           NewVariableSetString("a"),
			v2:           NewVariableSetString("a"),
			union:        NewVariableSetString("a"),
			complement:   NewVariableSetString(),
			intersection: NewVariableSetString("a"),
			symmetric:    NewVariableSetString(),
		},
		{
			v1:           NewVariableSetString("a"),
			v2:           NewVariableSetString("b"),
			union:        NewVariableSetString("a", "b"),
			complement:   NewVariableSetString("a"),
			intersection: NewVariableSetString(),
			symmetric:    NewVariableSetString("a", "b"),
		},
		{
			v1:           NewVariableSetString("a", "b"),
			v2:           NewVariableSetString("b", "c"),
			union:        NewVariableSetString("a", "b", "c"),
			complement:   NewVariableSetString("a"),
			intersection: NewVariableSetString("b"),
			symmetric:    NewVariableSetString("a", "c"),
		},
	} {
		check := func(
			f func(VariableSet, VariableSet) VariableSet,
			fname string,
			exp VariableSet,
			invert bool,
		) {
			a, b := test.v1, test.v2
			if invert {
				a, b = b, a
			}
			got := f(a, b)
			if !reflect.DeepEqual(got, exp) {
				t.Errorf("%d: %v %s %v; expected %v, got %v",
					i, a, fname, b, exp, got)
			}
		}
		check(VariableSet.Union, "union", test.union, false)
		check(VariableSet.Union, "union", test.union, true)

		check(VariableSet.Complement, "complement", test.complement, false)

		check(VariableSet.Intersection, "intersection", test.intersection, false)
		check(VariableSet.Intersection, "intersection", test.intersection, true)

		check(VariableSet.SymmetricDifference, "symmetric difference",
			test.symmetric, false)
		check(VariableSet.SymmetricDifference, "symmetric difference",
			test.symmetric, true)
	}
}
