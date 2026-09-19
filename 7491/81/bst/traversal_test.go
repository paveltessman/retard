package bst

import (
	"math"
	"math/rand"
	"slices"
	"testing"
)

// missing marks an absent child in a level order list. See buildTree.
const missing = math.MinInt

// order names the three traversal orders. The test uses it to build the
// answer that it wants from a traversal.
type order int

const (
	// preOrder visits the node, then the left subtree, then the right one.
	preOrder order = iota

	// inOrder visits the left subtree, then the node, then the right one. On
	// a BST this gives the values from small to large.
	inOrder

	// postOrder visits the left subtree, then the right one, then the node.
	postOrder
)

// traversalCase names one traversal of the package and gives the function.
//
// The contract of the function: walk(root) gives back the values of the tree
// in the order that the third field names. It gives back an empty list for a
// nil root, and it must not change the tree.
//
// An iterative traversal fits the same shape. Add one line for it:
//
//	{"PreorderIterative", preorderIterative, preOrder}
type traversalCase struct {
	name  string
	walk  func(root *node) []int
	order order
}

// traversals holds every traversal under test. Add one line for each new
// function, then run go test ./7491/81/bst.
var traversals = []traversalCase{
	{"PreorderRecursive", preorderRecursive, preOrder},
	{"InorderRecursive", inorderRecursive, inOrder},
	{"PostorderRecursive", postorderRecursive, postOrder},
}

// buildTree makes a tree out of a level order list. The list holds the value
// of each node, level by level and from left to right. A missing entry marks
// an absent child, and an absent child has no entries of its own.
//
// The list 1, 2, 3, missing, 4 makes this tree:
//
//	  1
//	 / \
//	2   3
//	 \
//	  4
func buildTree(level []int) *node {
	if len(level) == 0 || level[0] == missing {
		return nil
	}

	root := &node{val: level[0]}
	queue := []*node{root}
	i := 1

	for len(queue) > 0 && i < len(level) {
		n := queue[0]
		queue = queue[1:]

		if i < len(level) {
			if level[i] != missing {
				n.left = &node{val: level[i]}
				queue = append(queue, n.left)
			}
			i++
		}
		if i < len(level) {
			if level[i] != missing {
				n.right = &node{val: level[i]}
				queue = append(queue, n.right)
			}
			i++
		}
	}

	return root
}

// buildBST makes a BST out of vals by insertion, one value after the other.
// A value that is already in the tree goes to the left, so the tree can hold
// equal values.
func buildBST(vals []int) *node {
	var root *node

	for _, v := range vals {
		n := &node{val: v}

		if root == nil {
			root = n
			continue
		}

		cur := root
		for {
			if v <= cur.val {
				if cur.left == nil {
					cur.left = n
					break
				}
				cur = cur.left
				continue
			}
			if cur.right == nil {
				cur.right = n
				break
			}
			cur = cur.right
		}
	}

	return root
}

// buildChain makes a chain of n nodes with the values 1 to n. The chain goes
// down to the left if left is true, and down to the right if not. A chain is
// the deepest tree of n nodes, so it is the hard shape for a recursion and for
// a stack.
func buildChain(n int, left bool) *node {
	var root *node

	for i := n; i >= 1; i-- {
		cur := &node{val: i}
		if left {
			cur.left = root
		} else {
			cur.right = root
		}
		root = cur
	}

	return root
}

// buildComb makes a comb of n spine nodes and n leaves. The spine goes down to
// the left, and each spine node carries one leaf on the right. The values make
// the comb a BST of the numbers 2 to 2n+1.
//
// The comb is the hard shape for an iterative walk with an explicit stack: the
// stack holds the whole spine, because every spine node still owes its right
// leaf. A chain does not ask that much of a preorder stack.
//
//	    6
//	   / \
//	  4   7
//	 / \
//	2   5
//	 \
//	  3
func buildComb(n int) *node {
	var root *node

	for k := n; k >= 1; k-- {
		val := 2 * (n - k + 1)
		root = &node{val: val, left: root, right: &node{val: val + 1}}
	}

	return root
}

// refPreorder, refInorder and refPostorder are the reference traversals. They
// are the plain recursive walks, and the test compares the answer of the
// package with them.
func refPreorder(root *node, vals []int) []int {
	if root == nil {
		return vals
	}
	vals = append(vals, root.val)
	vals = refPreorder(root.left, vals)
	return refPreorder(root.right, vals)
}

func refInorder(root *node, vals []int) []int {
	if root == nil {
		return vals
	}
	vals = refInorder(root.left, vals)
	vals = append(vals, root.val)
	return refInorder(root.right, vals)
}

func refPostorder(root *node, vals []int) []int {
	if root == nil {
		return vals
	}
	vals = refPostorder(root.left, vals)
	vals = refPostorder(root.right, vals)
	return append(vals, root.val)
}

// reference gives the values of the tree in the order that o names.
func reference(root *node, o order) []int {
	switch o {
	case preOrder:
		return refPreorder(root, nil)
	case inOrder:
		return refInorder(root, nil)
	default:
		return refPostorder(root, nil)
	}
}

// shape gives a list that holds the whole tree: the value of each node in
// preorder, with a missing entry for each absent child. Two trees give the
// same list only if they hold the same values in the same places, so the test
// uses the list to see whether a traversal changed the tree.
func shape(root *node, out []int) []int {
	if root == nil {
		return append(out, missing)
	}
	out = append(out, root.val)
	out = shape(root.left, out)
	return shape(root.right, out)
}

// checkTraversal walks the tree and compares the answer with the reference
// walk. It also checks that the traversal leaves the tree as it is. It gives
// back false on the first wrong answer, so a loop over many trees can stop
// early.
func checkTraversal(t *testing.T, c traversalCase, name string, root *node) bool {
	t.Helper()

	want := reference(root, c.order)
	before := shape(root, nil)

	got := c.walk(root)

	if !slices.Equal(got, want) {
		t.Errorf("%s of the %s tree = %v, want %v", c.name, name, got, want)
		return false
	}
	if after := shape(root, nil); !slices.Equal(after, before) {
		t.Errorf("%s changed the %s tree", c.name, name)
		return false
	}
	return true
}

// everyTree calls f once for every shape of a binary tree of n nodes. The
// values go up from 1 in inorder, so every tree is a BST of the numbers 1 to
// n. Each tree owns its nodes.
func everyTree(n int, f func(root *node)) {
	var build func(lo, hi int) []*node

	build = func(lo, hi int) []*node {
		if lo > hi {
			return []*node{nil}
		}

		var trees []*node
		for val := lo; val <= hi; val++ {
			for _, left := range build(lo, val-1) {
				for _, right := range build(val+1, hi) {
					trees = append(trees, &node{val: val, left: cloneTree(left), right: cloneTree(right)})
				}
			}
		}
		return trees
	}

	for _, root := range build(1, n) {
		f(root)
	}
}

// cloneTree makes a copy of the tree with new nodes.
func cloneTree(root *node) *node {
	if root == nil {
		return nil
	}
	return &node{val: root.val, left: cloneTree(root.left), right: cloneTree(root.right)}
}

// randomTree makes a tree of n nodes with a random shape. The values lie
// between -max and max, so a small max gives many equal values.
func randomTree(n, max int, rnd *rand.Rand) *node {
	if n == 0 {
		return nil
	}

	left := rnd.Intn(n)

	return &node{
		val:   rnd.Intn(2*max+1) - max,
		left:  randomTree(left, max, rnd),
		right: randomTree(n-1-left, max, rnd),
	}
}

// randomVals gives n numbers between -max and max.
func randomVals(n, max int, rnd *rand.Rand) []int {
	vals := make([]int, n)
	for i := range vals {
		vals[i] = rnd.Intn(2*max+1) - max
	}
	return vals
}

// countUp gives the values 1, 2, ... n.
func countUp(n int) []int {
	vals := make([]int, n)
	for i := range vals {
		vals[i] = i + 1
	}
	return vals
}

func TestTraversalsAreRegistered(t *testing.T) {
	if len(traversals) == 0 {
		t.Fatal("no traversal is registered in traversals")
	}
}

// TestTraversal walks small trees and compares the answer with the answer
// written out by hand. The trees are not symmetric, so a traversal that takes
// the two children in the wrong order fails here.
func TestTraversal(t *testing.T) {
	cases := []struct {
		name  string
		level []int
		pre   []int
		in    []int
		post  []int
	}{
		{
			name:  "Nil",
			level: nil,
		},
		{
			name:  "OneNode",
			level: []int{7},
			pre:   []int{7},
			in:    []int{7},
			post:  []int{7},
		},
		{
			// 1 - 2 - 3, down to the left.
			name:  "LeftChain",
			level: []int{1, 2, missing, 3},
			pre:   []int{1, 2, 3},
			in:    []int{3, 2, 1},
			post:  []int{3, 2, 1},
		},
		{
			// 1 - 2 - 3, down to the right.
			name:  "RightChain",
			level: []int{1, missing, 2, missing, 3},
			pre:   []int{1, 2, 3},
			in:    []int{1, 2, 3},
			post:  []int{3, 2, 1},
		},
		{
			// The root, its left child, and the right child of that child.
			name:  "LeftThenRight",
			level: []int{1, 2, missing, missing, 3},
			pre:   []int{1, 2, 3},
			in:    []int{2, 3, 1},
			post:  []int{3, 2, 1},
		},
		{
			// The root, its right child, and the left child of that child.
			name:  "RightThenLeft",
			level: []int{1, missing, 2, 3},
			pre:   []int{1, 2, 3},
			in:    []int{1, 3, 2},
			post:  []int{3, 2, 1},
		},
		{
			name:  "ThreeNodesFull",
			level: []int{2, 1, 3},
			pre:   []int{2, 1, 3},
			in:    []int{1, 2, 3},
			post:  []int{1, 3, 2},
		},
		{
			name:  "OnlyALeftChild",
			level: []int{2, 1},
			pre:   []int{2, 1},
			in:    []int{1, 2},
			post:  []int{1, 2},
		},
		{
			name:  "OnlyARightChild",
			level: []int{1, missing, 2},
			pre:   []int{1, 2},
			in:    []int{1, 2},
			post:  []int{2, 1},
		},
		{
			name:  "PerfectSeven",
			level: []int{4, 2, 6, 1, 3, 5, 7},
			pre:   []int{4, 2, 1, 3, 6, 5, 7},
			in:    []int{1, 2, 3, 4, 5, 6, 7},
			post:  []int{1, 3, 2, 5, 7, 6, 4},
		},
		{
			// 1 with the left child 2 and the right child 3. 2 has the right
			// child 6, and 3 has the right child 5.
			name:  "Lopsided",
			level: []int{1, 2, 3, 4, missing, missing, 5, missing, 6},
			pre:   []int{1, 2, 4, 6, 3, 5},
			in:    []int{4, 6, 2, 1, 3, 5},
			post:  []int{6, 4, 2, 5, 3, 1},
		},
		{
			name:  "NegativeAndZeroValues",
			level: []int{0, -1, 1},
			pre:   []int{0, -1, 1},
			in:    []int{-1, 0, 1},
			post:  []int{-1, 1, 0},
		},
		{
			name:  "EqualValues",
			level: []int{1, 1, 2, 1},
			pre:   []int{1, 1, 1, 2},
			in:    []int{1, 1, 1, 2},
			post:  []int{1, 1, 2, 1},
		},
		{
			name:  "ExtremeValues",
			level: []int{0, math.MaxInt, math.MinInt + 1},
			pre:   []int{0, math.MaxInt, math.MinInt + 1},
			in:    []int{math.MaxInt, 0, math.MinInt + 1},
			post:  []int{math.MaxInt, math.MinInt + 1, 0},
		},
	}

	for _, c := range traversals {
		t.Run(c.name, func(t *testing.T) {
			for _, tc := range cases {
				t.Run(tc.name, func(t *testing.T) {
					want := tc.pre
					switch c.order {
					case inOrder:
						want = tc.in
					case postOrder:
						want = tc.post
					}

					root := buildTree(tc.level)
					before := shape(root, nil)

					got := c.walk(root)

					if !slices.Equal(got, want) {
						t.Errorf("%s of the %s tree = %v, want %v", c.name, tc.name, got, want)
					}
					if after := shape(root, nil); !slices.Equal(after, before) {
						t.Errorf("%s changed the %s tree", c.name, tc.name)
					}
				})
			}
		})
	}
}

// TestTraversalEveryShape walks every shape of a tree of up to 8 nodes. The
// values go up from 1 in inorder, so inorder must give 1, 2, ... n. This is
// the check that does not need the reference walk.
func TestTraversalEveryShape(t *testing.T) {
	const maxNodes = 8

	for _, c := range traversals {
		t.Run(c.name, func(t *testing.T) {
			for n := 0; n <= maxNodes; n++ {
				failed := false

				everyTree(n, func(root *node) {
					if failed {
						return
					}
					if !checkTraversal(t, c, "shape", root) {
						failed = true
						return
					}
					if c.order == inOrder {
						if want := countUp(n); !slices.Equal(c.walk(root), want) {
							t.Errorf("%s of a BST of %d nodes = %v, want %v", c.name, n, c.walk(root), want)
							failed = true
						}
					}
				})

				if failed {
					return
				}
			}
		})
	}
}

// TestTraversalRandom walks random trees of many sizes. A small max gives many
// equal values.
func TestTraversalRandom(t *testing.T) {
	maxes := []int{1, 3, 1000000}

	for _, c := range traversals {
		t.Run(c.name, func(t *testing.T) {
			rnd := rand.New(rand.NewSource(1))

			for _, max := range maxes {
				for range 200 {
					root := randomTree(rnd.Intn(200), max, rnd)

					if !checkTraversal(t, c, "random", root) {
						return
					}
				}
			}
		})
	}
}

// TestInorderOfBSTIsSorted builds a BST by insertion and checks that inorder
// gives the values from small to large. The answer comes from the standard
// library, not from the reference walk.
func TestInorderOfBSTIsSorted(t *testing.T) {
	maxes := []int{3, 20, 1000000}

	rnd := rand.New(rand.NewSource(2))

	for _, c := range traversals {
		if c.order != inOrder {
			continue
		}

		t.Run(c.name, func(t *testing.T) {
			for _, max := range maxes {
				for range 200 {
					vals := randomVals(rnd.Intn(200), max, rnd)

					want := slices.Clone(vals)
					slices.Sort(want)

					got := c.walk(buildBST(vals))

					if !slices.Equal(got, want) {
						t.Fatalf("%s of the BST of %v = %v, want %v", c.name, vals, got, want)
					}
				}
			}
		})
	}
}

// TestTraversalGivesANewSlice walks two trees and checks that the first answer
// stays as it is. A traversal that gives back one buffer again and again fails
// here.
func TestTraversalGivesANewSlice(t *testing.T) {
	first := buildTree([]int{4, 2, 6, 1, 3, 5, 7})
	second := buildTree([]int{40, 20, 60, 10, 30, 50, 70})

	for _, c := range traversals {
		t.Run(c.name, func(t *testing.T) {
			got := c.walk(first)
			want := slices.Clone(got)

			c.walk(second)

			if !slices.Equal(got, want) {
				t.Errorf("the answer of the first %s changed to %v, want %v", c.name, got, want)
			}
		})
	}
}

// TestTraversalDeepTrees walks the two chains and the comb of many nodes. A
// recursion needs one frame per node here, and an iterative walk needs one
// stack entry per node. A traversal that keeps the path in a slice of a fixed
// size fails here. The answers come from the shape of the tree, not from the
// reference walk.
func TestTraversalDeepTrees(t *testing.T) {
	const n = 20000
	const spine = 10000

	up := countUp(n)
	down := slices.Clone(up)
	slices.Reverse(down)

	// The comb holds the numbers 2 to 2*spine+1. Preorder goes down the spine
	// first, so it gives the even numbers from large to small, then the odd
	// numbers from small to large. Inorder gives every number in order.
	// Postorder takes each leaf just before its spine node.
	combPre := make([]int, 0, 2*spine)
	for val := 2 * spine; val >= 2; val -= 2 {
		combPre = append(combPre, val)
	}
	for val := 3; val <= 2*spine+1; val += 2 {
		combPre = append(combPre, val)
	}

	combIn := make([]int, 0, 2*spine)
	for val := 2; val <= 2*spine+1; val++ {
		combIn = append(combIn, val)
	}

	combPost := make([]int, 0, 2*spine)
	combPost = append(combPost, 3, 2)
	for val := 4; val <= 2*spine; val += 2 {
		combPost = append(combPost, val+1, val)
	}

	chains := []struct {
		name string
		root *node
		pre  []int
		in   []int
		post []int
	}{
		{"LeftChain", buildChain(n, true), up, down, down},
		{"RightChain", buildChain(n, false), up, up, down},
		{"Comb", buildComb(spine), combPre, combIn, combPost},
	}

	for _, c := range traversals {
		t.Run(c.name, func(t *testing.T) {
			for _, ch := range chains {
				t.Run(ch.name, func(t *testing.T) {
					want := ch.pre
					switch c.order {
					case inOrder:
						want = ch.in
					case postOrder:
						want = ch.post
					}

					if got := c.walk(ch.root); !slices.Equal(got, want) {
						t.Errorf("%s of the %s tree of %d nodes is wrong", c.name, ch.name, len(want))
					}
				})
			}
		})
	}
}

// TestTraversalLarge walks large trees. A traversal that copies the answer
// again and again is very slow here.
func TestTraversalLarge(t *testing.T) {
	const n = 200000

	rnd := rand.New(rand.NewSource(3))

	shapes := []struct {
		name string
		root *node
	}{
		{"Random", randomTree(n, 1000000, rnd)},
		{"RandomBST", buildBST(randomVals(n, 1000000, rnd))},
		{"EqualValues", randomTree(n, 0, rnd)},
	}

	for _, c := range traversals {
		t.Run(c.name, func(t *testing.T) {
			for _, s := range shapes {
				t.Run(s.name, func(t *testing.T) {
					checkTraversal(t, c, s.name, s.root)
				})
			}
		})
	}
}

func BenchmarkTraversal(b *testing.B) {
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
	}

	for _, c := range traversals {
		b.Run(c.name, func(b *testing.B) {
			for _, s := range shapes {
				b.Run(s.name, func(b *testing.B) {
					for b.Loop() {
						c.walk(s.root)
					}
				})
			}
		})
	}
}
