package ll

import (
	"slices"
	"testing"
)

// buildList makes a list out of vals.
// It gives back the head and the nodes in list order.
func buildList[T any](vals []T) (*node[T], []*node[T]) {
	nodes := make([]*node[T], len(vals))
	for i := range vals {
		nodes[i] = &node[T]{val: vals[i]}
	}
	for i := 0; i+1 < len(nodes); i++ {
		nodes[i].next = nodes[i+1]
	}
	if len(nodes) == 0 {
		return nil, nil
	}
	return nodes[0], nodes
}

// collectList walks the list and gives back the values and the nodes.
// It stops after limit steps, because a broken reverse can make a cycle.
func collectList[T any](t *testing.T, first *node[T], limit int) ([]T, []*node[T]) {
	t.Helper()

	vals := make([]T, 0, limit)
	nodes := make([]*node[T], 0, limit)
	for n := first; n != nil; n = n.next {
		if len(nodes) == limit {
			t.Fatalf("the list holds more than %d nodes, so reverse made a cycle", limit)
		}
		vals = append(vals, n.val)
		nodes = append(nodes, n)
	}
	return vals, nodes
}

// countUp gives the values 1, 2, ... n.
func countUp(n int) []int {
	vals := make([]int, n)
	for i := range vals {
		vals[i] = i + 1
	}
	return vals
}

func TestReverseInt(t *testing.T) {
	cases := []struct {
		name string
		in   []int
		want []int
	}{
		{"Empty", nil, nil},
		{"OneNode", []int{7}, []int{7}},
		{"TwoNodes", []int{1, 2}, []int{2, 1}},
		{"ThreeNodes", []int{1, 2, 3}, []int{3, 2, 1}},
		{"FiveNodes", []int{1, 2, 3, 4, 5}, []int{5, 4, 3, 2, 1}},
		{"DuplicateValues", []int{4, 4, 9, 4}, []int{4, 9, 4, 4}},
		{"ZeroValues", []int{0, 0, 0}, []int{0, 0, 0}},
		{"NegativeValues", []int{-1, 0, 1}, []int{1, 0, -1}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			first, _ := buildList(c.in)

			got, _ := collectList(t, reverse(first), len(c.in)+1)

			if !slices.Equal(got, c.want) {
				t.Errorf("reverse(%v) = %v, want %v", c.in, got, c.want)
			}
		})
	}
}

func TestReverseString(t *testing.T) {
	in := []string{"a", "b", "", "c"}
	want := []string{"c", "", "b", "a"}
	first, _ := buildList(in)

	got, _ := collectList(t, reverse(first), len(in)+1)

	if !slices.Equal(got, want) {
		t.Errorf("reverse(%v) = %v, want %v", in, got, want)
	}
}

// TestReverseKeepsNodes checks that reverse turns the links around in place.
// A solution that allocates new nodes fails here. Delete this test if you
// decide that copying is allowed.
func TestReverseKeepsNodes(t *testing.T) {
	in := countUp(5)
	first, nodes := buildList(in)
	want := slices.Clone(nodes)
	slices.Reverse(want)

	_, got := collectList(t, reverse(first), len(in)+1)

	if !slices.Equal(got, want) {
		t.Errorf("reverse used %d of the original nodes in the wrong order or made new ones", len(got))
	}
}

// TestReverseOldHeadIsTail checks that the old head ends the list.
// A forgotten "prev = nil" start leaves the old head pointing back, which is a cycle.
func TestReverseOldHeadIsTail(t *testing.T) {
	in := countUp(4)
	first, nodes := buildList(in)
	oldHead := nodes[0]

	last := reverse(first)
	for last != nil && last.next != nil {
		last = last.next
	}

	if last != oldHead {
		t.Errorf("the last node is %v, want the old head %v", last, oldHead)
	}
	if oldHead.next != nil {
		t.Errorf("the old head still points at %v, want nil", oldHead.next)
	}
}

func TestReverseTwiceGivesOriginal(t *testing.T) {
	in := countUp(6)
	first, _ := buildList(in)

	got, _ := collectList(t, reverse(reverse(first)), len(in)+1)

	if !slices.Equal(got, in) {
		t.Errorf("reverse twice = %v, want %v", got, in)
	}
}

func TestReverseLengths(t *testing.T) {
	for n := 0; n <= 50; n++ {
		in := countUp(n)
		want := slices.Clone(in)
		slices.Reverse(want)
		first, _ := buildList(in)

		got, _ := collectList(t, reverse(first), n+1)

		if !slices.Equal(got, want) {
			t.Fatalf("length %d: reverse = %v, want %v", n, got, want)
		}
	}
}

func BenchmarkReverse(b *testing.B) {
	first, _ := buildList(countUp(1000))

	for b.Loop() {
		first = reverse(first)
	}
}
