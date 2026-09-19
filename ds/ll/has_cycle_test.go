package ll

import (
	"slices"
	"testing"
	"time"
)

// buildCycle makes a list out of vals and points the tail at nodes[at].
// A negative at leaves the list open.
// It gives back the head and the nodes in list order.
func buildCycle[T any](vals []T, at int) (*node[T], []*node[T]) {
	first, nodes := buildList(vals)
	if at >= 0 && len(nodes) > 0 {
		nodes[len(nodes)-1].next = nodes[at]
	}
	return first, nodes
}

// snapshot gives back the next pointer of every node.
func snapshot[T any](nodes []*node[T]) []*node[T] {
	links := make([]*node[T], len(nodes))
	for i, n := range nodes {
		links[i] = n.next
	}
	return links
}

// callHasCycle runs has_cycle and fails the test if the call does not return.
// A solution that walks a cycle forever fails here instead of hanging the test binary.
// The stuck goroutine stays alive until the test binary exits.
func callHasCycle[T any](t *testing.T, first *node[T]) bool {
	t.Helper()

	done := make(chan bool, 1)
	go func() { done <- has_cycle(first) }()

	select {
	case got := <-done:
		return got
	case <-time.After(5 * time.Second):
		t.Fatal("has_cycle did not return after 5s, so it walks the list forever")
		return false
	}
}

func TestHasCycle(t *testing.T) {
	cases := []struct {
		name string
		n    int // number of nodes
		at   int // the tail points at this node, or -1 for an open list
		want bool
	}{
		{"Empty", 0, -1, false},
		{"OneNodeOpen", 1, -1, false},
		{"OneNodeSelfLoop", 1, 0, true},
		{"TwoNodesOpen", 2, -1, false},
		{"TwoNodesLoopToHead", 2, 0, true},
		{"TwoNodesLoopToTail", 2, 1, true},
		{"ThreeNodesOpen", 3, -1, false},
		{"ThreeNodesLoopToHead", 3, 0, true},
		{"ThreeNodesLoopToMiddle", 3, 1, true},
		{"ThreeNodesLoopToTail", 3, 2, true},
		{"LongOpen", 100, -1, false},
		{"LongLoopToHead", 100, 0, true},
		{"LongLoopToMiddle", 100, 50, true},
		{"LongLoopToTail", 100, 99, true},
		{"OddCycleLength", 9, 2, true},
		{"EvenCycleLength", 10, 2, true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			first, _ := buildCycle(countUp(c.n), c.at)

			got := callHasCycle(t, first)

			if got != c.want {
				t.Errorf("has_cycle(%d nodes, tail at %d) = %v, want %v", c.n, c.at, got, c.want)
			}
		})
	}
}

func TestHasCycleString(t *testing.T) {
	first, nodes := buildList([]string{"a", "b", "c"})

	if got := callHasCycle(t, first); got {
		t.Errorf("has_cycle on an open string list = %v, want false", got)
	}

	nodes[2].next = nodes[1]

	if got := callHasCycle(t, first); !got {
		t.Errorf("has_cycle on a closed string list = %v, want true", got)
	}
}

// TestHasCycleKeepsLinks checks that has_cycle reads the list and does not change it.
// A solution that cuts links or marks nodes fails here.
func TestHasCycleKeepsLinks(t *testing.T) {
	cases := []struct {
		name string
		at   int
	}{
		{"Open", -1},
		{"LoopToHead", 0},
		{"LoopToMiddle", 5},
		{"LoopToTail", 19},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			first, nodes := buildCycle(countUp(20), c.at)
			want := snapshot(nodes)

			callHasCycle(t, first)

			if got := snapshot(nodes); !slices.Equal(got, want) {
				t.Error("has_cycle changed the next pointers of the list")
			}
		})
	}
}

// TestHasCycleKeepsValues checks that has_cycle leaves the values alone.
func TestHasCycleKeepsValues(t *testing.T) {
	in := countUp(20)
	first, nodes := buildCycle(in, 7)

	callHasCycle(t, first)

	got := make([]int, len(nodes))
	for i, n := range nodes {
		got[i] = n.val
	}
	if !slices.Equal(got, in) {
		t.Errorf("has_cycle changed the values to %v, want %v", got, in)
	}
}

// TestHasCycleAllShapes walks every list length up to 30 and every cycle start.
// It catches an off by one error in the walk, because the pointers meet at a
// different place for each cycle length.
func TestHasCycleAllShapes(t *testing.T) {
	const maxLen = 30

	for n := 0; n <= maxLen; n++ {
		vals := countUp(n)

		first, _ := buildCycle(vals, -1)
		if got := callHasCycle(t, first); got {
			t.Fatalf("length %d, open list: has_cycle = %v, want false", n, got)
		}

		for at := 0; at < n; at++ {
			first, _ := buildCycle(vals, at)
			if got := callHasCycle(t, first); !got {
				t.Fatalf("length %d, tail at %d: has_cycle = %v, want true", n, at, got)
			}
		}
	}
}

func BenchmarkHasCycleOpen(b *testing.B) {
	first, _ := buildCycle(countUp(1000), -1)

	for b.Loop() {
		has_cycle(first)
	}
}

func BenchmarkHasCycleClosed(b *testing.B) {
	first, _ := buildCycle(countUp(1000), 500)

	for b.Loop() {
		has_cycle(first)
	}
}
