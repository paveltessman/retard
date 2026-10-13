package longestconsecutive

import (
	"testing"
)

func TestLongestConsecutive(t *testing.T) {

	cases := []struct {
		input []int
		want  int
	}{
		{[]int{2, 20, 4, 10, 3, 4, 5}, 4},
		{[]int{0, 3, 2, 5, 4, 6, 1, 1}, 7},
	}

	for _, c := range cases {
		got := longestConsecutive(c.input)
		if got != c.want {
			t.Fatalf("%v: want %d, got %d", c.input, c.want, got)
		}
	}
}
