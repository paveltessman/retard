package twosum

import (
	"slices"
	"testing"
)

// sortedPair checks that got holds two indices and gives them back in
// increasing order, because the task accepts the pair in any order.
func sortedPair(t *testing.T, got []int) []int {
	t.Helper()

	if len(got) != 2 {
		t.Fatalf("twosum gave %d indices %v, want 2", len(got), got)
	}
	pair := slices.Clone(got)
	slices.Sort(pair)
	return pair
}

// checkPair checks that the two indices are a real answer for nums and target.
// It reports an index out of range, a repeated index, and a wrong sum.
func checkPair(t *testing.T, nums []int, target int, got []int) {
	t.Helper()

	pair := sortedPair(t, got)
	for _, i := range pair {
		if i < 0 || i >= len(nums) {
			t.Fatalf("twosum gave index %d, want an index between 0 and %d", i, len(nums)-1)
		}
	}
	if pair[0] == pair[1] {
		t.Fatalf("twosum used index %d twice", pair[0])
	}
	if sum := nums[pair[0]] + nums[pair[1]]; sum != target {
		t.Fatalf("nums[%d] + nums[%d] = %d, want %d", pair[0], pair[1], sum, target)
	}
}

// powersFrom gives n values whose pair sums are all different.
// The values are 2^i - 1000, so the test also covers negative numbers.
// One sum belongs to one pair, so the answer of the task stays unique.
func powersFrom(n int) []int {
	vals := make([]int, n)
	for i := range vals {
		vals[i] = 1<<i - 1000
	}
	return vals
}

func TestTwoSum(t *testing.T) {
	cases := []struct {
		name   string
		nums   []int
		target int
		want   []int // the indices in increasing order
	}{
		{"Example1", []int{2, 7, 11, 15}, 9, []int{0, 1}},
		{"Example2", []int{3, 2, 4}, 6, []int{1, 2}},
		{"Example3", []int{3, 3}, 6, []int{0, 1}},
		{"TwoNumbers", []int{1, 2}, 3, []int{0, 1}},
		{"FirstAndLast", []int{5, 1, 2, 3, 6}, 11, []int{0, 4}},
		{"LastTwo", []int{1, 2, 3, 40, 50}, 90, []int{3, 4}},
		{"MiddleTwo", []int{9, 4, 6, 8}, 10, []int{1, 2}},
		{"Zeros", []int{0, 0}, 0, []int{0, 1}},
		{"ZeroAndValue", []int{7, 0, 3}, 7, []int{0, 1}},
		{"Negatives", []int{-3, -1, -4, -2}, -7, []int{0, 2}},
		{"MixedSigns", []int{-10, 20, 5, -5}, 10, []int{0, 1}},
		{"NegativeTarget", []int{4, -9, 8, 1}, -8, []int{1, 3}},
		{"DuplicateValues", []int{3, 2, 3, 5}, 6, []int{0, 2}},
		{"DuplicatesNotTheAnswer", []int{5, 5, 1, 2}, 3, []int{2, 3}},
		{"HalfOfTargetOnce", []int{1, 5, 4, 6}, 10, []int{2, 3}},
		{"Limits", []int{1000000000, -1000000000, 1}, 0, []int{0, 1}},
		{"BigSum", []int{999999999, 1, 1000000000}, 1999999999, []int{0, 2}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := sortedPair(t, twosum(c.nums, c.target))

			if !slices.Equal(got, c.want) {
				t.Errorf("twosum(%v, %d) = %v, want %v", c.nums, c.target, got, c.want)
			}
		})
	}
}

// TestTwoSumAllPairs walks every array length up to 12 and every pair inside it.
// The values make every pair sum different, so one pair is the only answer.
// The test catches a loop that starts or stops one step too late.
func TestTwoSumAllPairs(t *testing.T) {
	const maxLen = 12

	for n := 2; n <= maxLen; n++ {
		nums := powersFrom(n)
		for a := 0; a < n; a++ {
			for b := a + 1; b < n; b++ {
				target := nums[a] + nums[b]
				want := []int{a, b}

				got := sortedPair(t, twosum(nums, target))

				if !slices.Equal(got, want) {
					t.Fatalf("length %d: twosum(%v, %d) = %v, want %v", n, nums, target, got, want)
				}
			}
		}
	}
}

// TestTwoSumRepeatedValue puts the same value in every place but two.
// It moves the answer through every pair of positions, and the array around
// the answer holds one value many times.
func TestTwoSumRepeatedValue(t *testing.T) {
	const n = 9

	for a := 0; a < n; a++ {
		for b := a + 1; b < n; b++ {
			nums := make([]int, n)
			for i := range nums {
				nums[i] = 100
			}
			nums[a] = 1
			nums[b] = 2
			want := []int{a, b}

			got := sortedPair(t, twosum(nums, 3))

			if !slices.Equal(got, want) {
				t.Fatalf("twosum(%v, 3) = %v, want %v", nums, got, want)
			}
		}
	}
}

// TestTwoSumSameValueTwice makes the answer out of two equal numbers.
// A solution that removes duplicate values fails here.
func TestTwoSumSameValueTwice(t *testing.T) {
	nums := []int{1, 7, 3, 7, 9}
	want := []int{1, 3}

	got := sortedPair(t, twosum(nums, 14))

	if !slices.Equal(got, want) {
		t.Errorf("twosum(%v, 14) = %v, want %v", nums, got, want)
	}
}

// TestTwoSumKeepsInput checks that twosum leaves the array alone.
// A solution that sorts nums in place gives back indices of the sorted array,
// so it fails this test and the tests above.
func TestTwoSumKeepsInput(t *testing.T) {
	nums := []int{8, 3, 11, 2, 7, 5}
	want := slices.Clone(nums)

	twosum(nums, 19)

	if !slices.Equal(nums, want) {
		t.Errorf("twosum changed the array to %v, want %v", nums, want)
	}
}

// TestTwoSumLarge uses the largest array of the constraints and puts the
// answer at the end, so the solution must read the whole array.
func TestTwoSumLarge(t *testing.T) {
	const n = 10000

	nums := make([]int, n)
	for i := range nums {
		nums[i] = i
	}
	// Only the last two values add up to this target.
	target := (n - 1) + (n - 2)
	want := []int{n - 2, n - 1}

	got := sortedPair(t, twosum(nums, target))

	if !slices.Equal(got, want) {
		t.Errorf("twosum(0..%d, %d) = %v, want %v", n-1, target, got, want)
	}
	checkPair(t, nums, target, got)
}

func BenchmarkTwoSumSmall(b *testing.B) {
	nums := []int{2, 7, 11, 15}

	for b.Loop() {
		twosum(nums, 26)
	}
}

// BenchmarkTwoSumLarge shows the cost of the follow-up question.
// A map solution stays near the array length. A double loop grows much faster.
func BenchmarkTwoSumLarge(b *testing.B) {
	const n = 10000

	nums := make([]int, n)
	for i := range nums {
		nums[i] = i
	}
	target := (n - 1) + (n - 2)

	for b.Loop() {
		twosum(nums, target)
	}
}
