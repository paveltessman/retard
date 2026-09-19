package bst

import (
	"math"
	"math/rand"
	"slices"
	"testing"
)

// The contract of levelOrder: levelOrder(root) gives back one list per level
// of the tree, from the root down. Each list holds the values of that level
// from left to right. A nil root gives back no level at all, and no level is
// empty. The function must not change the tree, and each call must give back
// new slices.

// refLevelOrder is the reference answer. It walks the tree in preorder and
// puts each value in the list of its depth, so it uses no queue. The test
// compares the answer of the package with it.
func refLevelOrder(root *node) [][]int {
	var levels [][]int

	var visit func(n *node, depth int)
	visit = func(n *node, depth int) {
		if n == nil {
			return
		}
		if depth == len(levels) {
			levels = append(levels, nil)
		}
		levels[depth] = append(levels[depth], n.val)
		visit(n.left, depth+1)
		visit(n.right, depth+1)
	}
	visit(root, 0)

	return levels
}

// equalLevels compares two answers. A nil list and an empty list count as the
// same, on the outer level and on a single level.
func equalLevels(a, b [][]int) bool {
	return slices.EqualFunc(a, b, slices.Equal)
}

// levelProblem reports whether the levels of the answer sit in memory of
// their own. It gives back an empty string when they do, and a short note on
// the problem when they do not.
//
// The check runs in two steps. First it writes the number of the level in
// every entry of that level, then it reads the entries back: a level that
// shares memory with another level loses its own numbers. Second it writes a
// marker in the spare room of each level, because a level that can grow into
// the next one wipes that next level out as soon as the caller appends to it.
//
// The function puts the values back before it gives an answer.
func levelProblem(got [][]int) string {
	old := make([][]int, len(got))
	for i, level := range got {
		old[i] = slices.Clone(level)
	}

	problem := ""

	for i, level := range got {
		for j := range level {
			level[j] = i
		}
	}
	for i, level := range got {
		for j := range level {
			if level[j] != i {
				problem = "two levels share the same memory"
			}
		}
	}

	if problem == "" {
		for _, level := range got {
			spare := level[len(level):cap(level)]
			for j := range spare {
				spare[j] = math.MinInt
			}
		}
		for i, level := range got {
			for j := range level {
				if level[j] != i {
					problem = "the spare room of one level reaches into the next level"
				}
			}
		}
	}

	for i, level := range got {
		copy(level, old[i])
	}

	return problem
}

// checkLevelOrder walks the tree and compares the answer with the reference
// answer. It also checks that no level is empty, that the levels sit in
// memory of their own, and that the walk leaves the tree as it is. It gives
// back false on the first wrong answer, so a loop over many trees can stop
// early.
func checkLevelOrder(t *testing.T, name string, root *node) bool {
	t.Helper()

	want := refLevelOrder(root)
	before := shape(root, nil)

	got := levelOrder(root)

	if !equalLevels(got, want) {
		t.Errorf("levelOrder of the %s tree = %v, want %v", name, got, want)
		return false
	}
	for i, level := range got {
		if len(level) == 0 {
			t.Errorf("level %d of the %s tree is empty", i, name)
			return false
		}
	}
	if problem := levelProblem(got); problem != "" {
		t.Errorf("the levels of the %s tree are not separate: %s", name, problem)
		return false
	}
	if after := shape(root, nil); !slices.Equal(after, before) {
		t.Errorf("levelOrder changed the %s tree", name)
		return false
	}
	return true
}

// TestLevelOrder walks small trees and compares the answer with the answer
// written out by hand. The trees are not symmetric, so an implementation that
// takes the two children in the wrong order fails here.
func TestLevelOrder(t *testing.T) {
	cases := []struct {
		name  string
		level []int
		want  [][]int
	}{
		{
			name:  "Nil",
			level: nil,
			want:  nil,
		},
		{
			name:  "OneNode",
			level: []int{7},
			want:  [][]int{{7}},
		},
		{
			// 1 - 2 - 3, down to the left.
			name:  "LeftChain",
			level: []int{1, 2, missing, 3},
			want:  [][]int{{1}, {2}, {3}},
		},
		{
			// 1 - 2 - 3, down to the right.
			name:  "RightChain",
			level: []int{1, missing, 2, missing, 3},
			want:  [][]int{{1}, {2}, {3}},
		},
		{
			// The root, its left child, and the right child of that child.
			name:  "LeftThenRight",
			level: []int{1, 2, missing, missing, 3},
			want:  [][]int{{1}, {2}, {3}},
		},
		{
			// The root, its right child, and the left child of that child.
			name:  "RightThenLeft",
			level: []int{1, missing, 2, 3},
			want:  [][]int{{1}, {2}, {3}},
		},
		{
			name:  "ThreeNodesFull",
			level: []int{2, 1, 3},
			want:  [][]int{{2}, {1, 3}},
		},
		{
			name:  "OnlyALeftChild",
			level: []int{2, 1},
			want:  [][]int{{2}, {1}},
		},
		{
			name:  "OnlyARightChild",
			level: []int{1, missing, 2},
			want:  [][]int{{1}, {2}},
		},
		{
			name:  "PerfectSeven",
			level: []int{4, 2, 6, 1, 3, 5, 7},
			want:  [][]int{{4}, {2, 6}, {1, 3, 5, 7}},
		},
		{
			// The left child of the root has no node under it on the last
			// level, so the last level holds one value only. An
			// implementation that counts the nodes of a level from the width
			// of the level above fails here.
			name:  "Lopsided",
			level: []int{1, 2, 3, 4, missing, missing, 5, missing, 6},
			want:  [][]int{{1}, {2, 3}, {4, 5}, {6}},
		},
		{
			// The two deep nodes sit under different parents and far apart.
			// Their level holds them side by side.
			name:  "TwoFarLeaves",
			level: []int{1, 2, 3, 4, missing, missing, 5},
			want:  [][]int{{1}, {2, 3}, {4, 5}},
		},
		{
			name:  "NegativeAndZeroValues",
			level: []int{0, -1, 1},
			want:  [][]int{{0}, {-1, 1}},
		},
		{
			name:  "EqualValues",
			level: []int{1, 1, 2, 1},
			want:  [][]int{{1}, {1, 2}, {1}},
		},
		{
			name:  "ExtremeValues",
			level: []int{0, math.MaxInt, math.MinInt + 1},
			want:  [][]int{{0}, {math.MaxInt, math.MinInt + 1}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := buildTree(tc.level)
			before := shape(root, nil)

			got := levelOrder(root)

			if !equalLevels(got, tc.want) {
				t.Errorf("levelOrder of the %s tree = %v, want %v", tc.name, got, tc.want)
			}
			if after := shape(root, nil); !slices.Equal(after, before) {
				t.Errorf("levelOrder changed the %s tree", tc.name)
			}
		})
	}
}

// TestLevelOrderOfCompleteTree builds a complete tree out of a list of values
// and checks that the answer is that same list, cut into levels of 1, 2, 4, 8
// and so on. The answer comes from the list, not from the reference walk,
// because buildTree reads the list level by level.
func TestLevelOrderOfCompleteTree(t *testing.T) {
	rnd := rand.New(rand.NewSource(5))

	for n := 1; n <= 300; n++ {
		vals := randomVals(n, 1000000, rnd)

		var want [][]int
		for i, width := 0, 1; i < n; i, width = i+width, 2*width {
			want = append(want, vals[i:min(i+width, n)])
		}

		if got := levelOrder(buildTree(vals)); !equalLevels(got, want) {
			t.Fatalf("levelOrder of the complete tree of %d nodes = %v, want %v", n, got, want)
		}
	}
}

// TestLevelOrderEveryShape walks every shape of a tree of up to 8 nodes.
func TestLevelOrderEveryShape(t *testing.T) {
	const maxNodes = 8

	for n := 0; n <= maxNodes; n++ {
		failed := false

		everyTree(n, func(root *node) {
			if failed {
				return
			}
			if !checkLevelOrder(t, "shape", root) {
				failed = true
			}
		})

		if failed {
			return
		}
	}
}

// TestLevelOrderRandom walks random trees of many sizes. A small max gives
// many equal values.
func TestLevelOrderRandom(t *testing.T) {
	maxes := []int{1, 3, 1000000}

	rnd := rand.New(rand.NewSource(1))

	for _, max := range maxes {
		for range 200 {
			root := randomTree(rnd.Intn(200), max, rnd)

			if !checkLevelOrder(t, "random", root) {
				return
			}
		}
	}
}

// TestLevelOrderGivesNewSlices walks two trees and checks that the first
// answer stays as it is. An implementation that gives back one buffer again
// and again fails here.
func TestLevelOrderGivesNewSlices(t *testing.T) {
	first := buildTree([]int{4, 2, 6, 1, 3, 5, 7})
	second := buildTree([]int{40, 20, 60, 10, 30, 50, 70})

	got := levelOrder(first)

	want := make([][]int, len(got))
	for i, level := range got {
		want[i] = slices.Clone(level)
	}

	levelOrder(second)

	if !equalLevels(got, want) {
		t.Errorf("the answer of the first call changed to %v, want %v", got, want)
	}
}

// TestLevelOrderDeepTrees walks the two chains and the comb of many nodes.
// The chain gives one node per level, and the comb gives two. An
// implementation that keeps a level in a slice of a fixed size fails here.
// The answers come from the shape of the tree, not from the reference walk.
func TestLevelOrderDeepTrees(t *testing.T) {
	const n = 20000
	const spine = 10000

	// Each chain of n nodes holds the values 1 to n, one value per level.
	chainLevels := make([][]int, n)
	for i := range chainLevels {
		chainLevels[i] = []int{i + 1}
	}

	// The comb holds the numbers 2 to 2*spine+1. The top of the spine is
	// 2*spine, and each step down the spine takes 2 off that value. A spine
	// node with the value v carries the leaf v+1 on its right. The last level
	// holds the leaf of the smallest spine node alone.
	combLevels := make([][]int, 0, spine+1)
	combLevels = append(combLevels, []int{2 * spine})
	for i := 1; i < spine; i++ {
		combLevels = append(combLevels, []int{2*spine - 2*i, 2*spine - 2*i + 3})
	}
	combLevels = append(combLevels, []int{3})

	cases := []struct {
		name string
		root *node
		want [][]int
	}{
		{"LeftChain", buildChain(n, true), chainLevels},
		{"RightChain", buildChain(n, false), chainLevels},
		{"Comb", buildComb(spine), combLevels},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := levelOrder(tc.root)

			if len(got) != len(tc.want) {
				t.Fatalf("levelOrder of the %s tree gave %d levels, want %d", tc.name, len(got), len(tc.want))
			}
			for i := range got {
				if !slices.Equal(got[i], tc.want[i]) {
					t.Fatalf("level %d of the %s tree = %v, want %v", i, tc.name, got[i], tc.want[i])
				}
			}
		})
	}
}

// TestLevelOrderLarge walks large trees. An implementation that copies the
// answer again and again is very slow here.
func TestLevelOrderLarge(t *testing.T) {
	const n = 200000

	rnd := rand.New(rand.NewSource(3))

	shapes := []struct {
		name string
		root *node
	}{
		{"Random", randomTree(n, 1000000, rnd)},
		{"RandomBST", buildBST(randomVals(n, 1000000, rnd))},
		{"EqualValues", randomTree(n, 0, rnd)},
		{"Perfect", buildTree(countUp(1<<18 - 1))},
	}

	for _, s := range shapes {
		t.Run(s.name, func(t *testing.T) {
			checkLevelOrder(t, s.name, s.root)
		})
	}
}

func BenchmarkLevelOrder(b *testing.B) {
	const n = 10000

	rnd := rand.New(rand.NewSource(4))

	shapes := []struct {
		name string
		root *node
	}{
		{"Random", randomTree(n, 1000000, rnd)},
		{"RandomBST", buildBST(randomVals(n, 1000000, rnd))},
		{"LeftChain", buildChain(n, true)},
		{"RightChain", buildChain(n, false)},
		{"Comb", buildComb(n / 2)},
	}

	for _, s := range shapes {
		b.Run(s.name, func(b *testing.B) {
			for b.Loop() {
				levelOrder(s.root)
			}
		})
	}
}
