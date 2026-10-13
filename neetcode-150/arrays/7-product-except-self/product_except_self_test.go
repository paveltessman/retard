package productexceptself

import (
	"slices"
	"testing"
)

func TestProductExceptSelf(t *testing.T) {
	cases := map[string]struct {
		input []int
		want  []int
	}{
		"case1": {[]int{1, 2, 4, 6}, []int{48, 24, 12, 8}},
	}

	for desc, c := range cases {

		got := productExceptSelf(c.input)

		if !slices.Equal(got, c.want) {
			t.Fatalf("%s: want %v, got %v", desc, c.want, got)
		}

	}
}
