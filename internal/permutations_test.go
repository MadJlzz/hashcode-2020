package internal

import (
	"testing"
)

func TestPermutations(t *testing.T) {
	in := []int{1, 2, 3}
	out := Permutations(in)

	v := func(a []int) []interface{} {
		res := make([]interface{}, len(a))
		for i, v := range a {
			res[i] = v
		}
		return res
	}

	Assert(t, 6, len(out))
	AssertArrays(t, []interface{}{1, 2, 3}, v(out[0]))
	AssertArrays(t, []interface{}{2, 1, 3}, v(out[1]))
	AssertArrays(t, []interface{}{3, 2, 1}, v(out[2]))
	AssertArrays(t, []interface{}{2, 3, 1}, v(out[3]))
	AssertArrays(t, []interface{}{3, 1, 2}, v(out[4]))
	AssertArrays(t, []interface{}{1, 3, 2}, v(out[5]))
}
