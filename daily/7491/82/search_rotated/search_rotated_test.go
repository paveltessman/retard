package searchrotated

import (
	"math/rand"
	"slices"
	"testing"
	"time"
)

// brute gives the answer of the task in the slow and plain way: it looks at
// every number and gives back the index of the one that is target. The values
// are distinct, so at most one index fits. The tests compare the solution with
// it.
func brute(nums []int, target int) int {
	for i, v := range nums {
		if v == target {
			return i
		}
	}

	return -1
}

// rotate gives back a new array with the numbers of nums left rotated by k
// places, as the task describes it. rotate(nums, 0) is a copy of nums.
func rotate(nums []int, k int) []int {
	n := len(nums)
	k %= n

	out := make([]int, 0, n)
	out = append(out, nums[k:]...)
	out = append(out, nums[:k]...)

	return out
}

// check calls search on a copy of nums and compares the answer with want. It
// also reads the copy back, because search must leave the array as it was. It
// gives back false on the first error, so a loop can stop.
func check(t *testing.T, nums []int, target, want int) bool {
	t.Helper()

	got := slices.Clone(nums)
	ok := true

	if i := search(got, target); i != want {
		t.Errorf("search(%v, %d) = %d, want %d", nums, target, i, want)
		ok = false
	}
	if !slices.Equal(got, nums) {
		t.Errorf("search(%v, %d) changed the array to %v", nums, target, got)
		ok = false
	}

	return ok
}

// sortedDistinct gives n distinct numbers from low to high in ascending order.
// The gaps between them are random, so the misses fall between the numbers as
// well as outside of the array.
func sortedDistinct(n, low, high int, rnd *rand.Rand) []int {
	size := high - low + 1
	if n > size {
		panic("the range is too small for n distinct numbers")
	}

	// Take n numbers of the range without repetition, then sort them.
	picked := make(map[int]bool, n)
	nums := make([]int, 0, n)

	for len(nums) < n {
		v := low + rnd.Intn(size)
		if picked[v] {
			continue
		}

		picked[v] = true
		nums = append(nums, v)
	}

	slices.Sort(nums)

	return nums
}

func TestSearch(t *testing.T) {
	cases := []struct {
		name   string
		nums   []int
		target int
		want   int
	}{
		{"Example1", []int{4, 5, 6, 7, 0, 1, 2}, 0, 4},
		{"Example2", []int{4, 5, 6, 7, 0, 1, 2}, 3, -1},
		{"Example3", []int{1}, 0, -1},
		{"OneNumberHit", []int{1}, 1, 0},
		{"OneNumberBelow", []int{1}, -5, -1},
		{"OneNumberAbove", []int{1}, 5, -1},
		{"TwoNumbersSorted", []int{1, 3}, 3, 1},
		{"TwoNumbersRotated", []int{3, 1}, 3, 0},
		{"TwoNumbersRotatedSecond", []int{3, 1}, 1, 1},
		{"TwoNumbersMiss", []int{3, 1}, 2, -1},
		{"NotRotated", []int{0, 1, 2, 4, 5, 6, 7}, 4, 3},
		{"NotRotatedFirst", []int{0, 1, 2, 4, 5, 6, 7}, 0, 0},
		{"NotRotatedLast", []int{0, 1, 2, 4, 5, 6, 7}, 7, 6},
		{"NotRotatedMissInTheMiddle", []int{0, 1, 2, 4, 5, 6, 7}, 3, -1},
		{"FirstNumber", []int{4, 5, 6, 7, 0, 1, 2}, 4, 0},
		{"LastNumber", []int{4, 5, 6, 7, 0, 1, 2}, 2, 6},
		{"EndOfTheLeftPart", []int{4, 5, 6, 7, 0, 1, 2}, 7, 3},
		{"StartOfTheRightPart", []int{4, 5, 6, 7, 0, 1, 2}, 0, 4},
		{"MissBetweenTheParts", []int{4, 5, 6, 7, 0, 1, 2}, 3, -1},
		{"MissBelowEveryNumber", []int{4, 5, 6, 7, 0, 1, 2}, -1, -1},
		{"MissAboveEveryNumber", []int{4, 5, 6, 7, 0, 1, 2}, 8, -1},
		{"RotatedByOne", []int{7, 0, 1, 2, 4, 5, 6}, 7, 0},
		{"RotatedByAllButOne", []int{2, 4, 5, 6, 7, 0, 1}, 1, 6},
		{"NegativeNumbers", []int{3, 4, 5, -3, -2, -1, 0, 1, 2}, -3, 3},
		{"NegativeTarget", []int{3, 4, 5, -3, -2, -1, 0, 1, 2}, -1, 5},
		{"NegativeMiss", []int{3, 4, 5, -3, -2, -1, 0, 1, 2}, -4, -1},
		{"LargestValues", []int{10000, -10000}, 10000, 0},
		{"SmallestValue", []int{10000, -10000}, -10000, 1},
		{"MissBetweenTheLargestValues", []int{10000, -10000}, 0, -1},
		{"LeftPartOfOneNumber", []int{5, 1, 2, 3, 4}, 5, 0},
		{"RightPartOfOneNumber", []int{2, 3, 4, 5, 1}, 1, 4},
		{"MiddleOfTheLeftPart", []int{5, 6, 7, 8, 1, 2, 3, 4}, 6, 1},
		{"MiddleOfTheRightPart", []int{5, 6, 7, 8, 1, 2, 3, 4}, 3, 6},
		{"WideGaps", []int{100, 1000, -1000, -100}, -100, 3},
		{"WideGapsMiss", []int{100, 1000, -1000, -100}, 0, -1},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			check(t, c.nums, c.target, c.want)
		})
	}
}

// TestSearchEveryRotation walks every rotation of every short array and asks
// for every target that the array holds, for every target between two of the
// numbers, and for the targets outside of the array. It catches a solution
// that picks the wrong half, a solution that drops the number at an end of a
// half, and a solution that reads the not rotated array wrong.
func TestSearchEveryRotation(t *testing.T) {
	const maxLen = 9

	for n := 1; n <= maxLen; n++ {
		bases := [][]int{
			// The gaps of 2 give a miss between every two numbers.
			make([]int, n),
			// The same array moved down, so a part of the numbers is negative.
			make([]int, n),
		}

		for i := range n {
			bases[0][i] = 2 * i
			bases[1][i] = 2*i - n
		}

		for _, base := range bases {
			for k := range n {
				nums := rotate(base, k)

				for target := base[0] - 2; target <= base[n-1]+2; target++ {
					if !check(t, nums, target, brute(nums, target)) {
						return
					}
				}
			}
		}
	}
}

// TestSearchRandom compares the solution with brute on random arrays with
// random gaps and a random rotation. It asks for every number of the array and
// for random targets, most of which are misses.
func TestSearchRandom(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))

	for range 300 {
		n := 1 + rnd.Intn(60)

		base := sortedDistinct(n, -100, 100, rnd)
		nums := rotate(base, rnd.Intn(n))

		for _, target := range base {
			if !check(t, nums, target, brute(nums, target)) {
				return
			}
		}

		for range 20 {
			target := rnd.Intn(221) - 110

			if !check(t, nums, target, brute(nums, target)) {
				return
			}
		}
	}
}

// TestSearchKeepsNeighbours gives a part of a larger array to search. The
// numbers around the part break the order, so a solution that reads outside of
// nums gives a wrong answer here.
func TestSearchKeepsNeighbours(t *testing.T) {
	arr := []int{50, 60, 30, 40, 10, 20, 70, 80}
	nums := arr[2:6] // 30, 40, 10, 20

	for target, want := range map[int]int{30: 0, 40: 1, 10: 2, 20: 3, 50: -1, 60: -1, 70: -1, 80: -1, 35: -1} {
		check(t, nums, target, want)
	}
}

// TestSearchLarge uses the largest array of the constraints, with the largest
// and the smallest values. It asks for every number of the array and for a
// miss next to it, over several rotations.
func TestSearchLarge(t *testing.T) {
	const n = 5000

	// The values stay inside -10^4 and 10^4, and the gap of 4 leaves room for
	// a miss on both sides of every number.
	base := make([]int, n)
	for i := range base {
		base[i] = -10000 + 4*i
	}

	for _, k := range []int{0, 1, 2, n / 3, n / 2, n - 2, n - 1} {
		nums := rotate(base, k)

		for i, v := range base {
			want := (i - k + n) % n

			if !check(t, nums, v, want) {
				return
			}
			if !check(t, nums, v+1, -1) {
				return
			}
			if !check(t, nums, v-1, -1) {
				return
			}
		}
	}
}

// TestSearchIsLogarithmic makes many searches in a very large rotated array. A
// binary search needs some milliseconds. A linear scan needs minutes, so it
// runs into the limit. The array is larger than the constraints allow, because
// a linear scan over 5000 numbers is still fast.
func TestSearchIsLogarithmic(t *testing.T) {
	if testing.Short() {
		t.Skip("the test needs a large array")
	}

	const n = 1 << 22
	const queries = 100000
	const limit = 10 * time.Second

	// The numbers are even, so every odd target is a miss.
	base := make([]int, n)
	for i := range base {
		base[i] = 2 * i
	}

	const k = n / 3
	nums := rotate(base, k)

	rnd := rand.New(rand.NewSource(2))
	start := time.Now()

	for range queries {
		target := rnd.Intn(2 * n)

		want := -1
		if target%2 == 0 {
			want = (target/2 - k + n) % n
		}

		if got := search(nums, target); got != want {
			t.Fatalf("search(an array of %d numbers rotated by %d, %d) = %d, want %d", n, k, target, got, want)
		}
		if time.Since(start) > limit {
			t.Fatalf("search needed more than %v for the searches in an array of %d numbers, so it is not logarithmic", limit, n)
		}
	}
}

func BenchmarkSearch(b *testing.B) {
	const n = 100000

	base := make([]int, n)
	for i := range base {
		base[i] = 2 * i
	}

	nums := rotate(base, n/3)

	shapes := []struct {
		name   string
		target int
	}{
		{"HitInTheLeftPart", 2 * (n / 2)},
		{"HitInTheRightPart", 2 * (n / 6)},
		{"HitAtTheFirstIndex", 2 * (n / 3)},
		{"HitAtTheLastIndex", 2*(n/3) - 2},
		{"Miss", 2*(n/2) + 1},
		{"MissOutsideOfTheArray", 2 * n},
	}

	for _, s := range shapes {
		b.Run(s.name, func(b *testing.B) {
			for b.Loop() {
				search(nums, s.target)
			}
		})
	}
}

func TestStartIndex(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		want int
	}{
		{"Example0", []int{2, 4, 5, 6, 7, 0, 1}, 5},
		{"Example1", []int{4, 5, 6, 7, 0, 1, 2}, 4},
		{"Example2", []int{0, 1, 2, 4, 5, 6, 7}, 0},
		{"Example3", []int{5, 6, 7, 0, 1, 2, 4}, 3},
		{"Example4", []int{5, 5, 5, 5, 5, 5, 5}, 0},
		{"Example5", []int{1}, 0},
		{"Example6", []int{1, 1}, 0},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := startIndex(c.nums)
			if want := c.want; got != want {
				t.Errorf("Incorrect offset %v: got %d, want %d", c.nums, got, want)
			}
		})
	}
}
