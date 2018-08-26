package tnt

import (
	"fmt"
	"sort"
)

// VariableSet is a set of Variables.
type VariableSet map[Variable]struct{}

func NewVariableSetString(strs ...string) VariableSet {
	v := make(VariableSet)
	for _, str := range strs {
		v[Variable(str)] = struct{}{}
	}
	return v
}

func NewVariableSet(vars ...Variable) VariableSet {
	vs := make(VariableSet)
	for _, v := range vars {
		vs[v] = struct{}{}
	}
	return vs
}

// Union returns a set including elements from v and v2
func (v VariableSet) Union(v2 VariableSet) VariableSet {
	union := make(VariableSet)
	union.add(v)
	union.add(v2)
	return union
}

// Complement returns a set of elements that are in v but not in v2
func (v VariableSet) Complement(v2 VariableSet) VariableSet {
	complement := make(VariableSet)
	complement.add(v)
	complement.subtract(v2)
	return complement
}

// Intersection returns a set of elements that are in both v and v2
func (v VariableSet) Intersection(v2 VariableSet) VariableSet {
	intersection := make(VariableSet)
	for item := range v {
		if _, ok := v2[item]; ok {
			intersection[item] = struct{}{}
		}
	}
	return intersection
}

// SymmetricDifference returns a set of elements that are in v or v2,
// but not both.
func (v VariableSet) SymmetricDifference(v2 VariableSet) VariableSet {
	diff := v.Complement(v2)
	diff.add(v2.Complement(v))
	return diff
}

func (v VariableSet) String() string {
	slice := make([]string, 0, len(v))
	for k, _ := range v {
		slice = append(slice, string(k))
	}
	sort.Strings(slice)
	return fmt.Sprintf("%s", slice)
}

func (v VariableSet) add(v2 VariableSet) {
	for item := range v2 {
		v[item] = struct{}{}
	}
}

func (v VariableSet) subtract(v2 VariableSet) {
	for item := range v2 {
		delete(v, item)
	}
}
