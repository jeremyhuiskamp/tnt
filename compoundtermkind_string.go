// Code generated by "stringer -type CompoundTermKind"; DO NOT EDIT.

package tnt

import "strconv"

const _CompoundTermKind_name = "PLUSMULTIPLY"

var _CompoundTermKind_index = [...]uint8{0, 4, 12}

func (i CompoundTermKind) String() string {
	if i < 0 || i >= CompoundTermKind(len(_CompoundTermKind_index)-1) {
		return "CompoundTermKind(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _CompoundTermKind_name[_CompoundTermKind_index[i]:_CompoundTermKind_index[i+1]]
}