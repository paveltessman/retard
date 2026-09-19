package movezeroes

import (
	"math/rand"
	"slices"
	"testing"
)

// brute gives the answer of the task in the slow and plain way: it copies the
// numbers that are not 0 into a new slice, then adds the 0's at the end. The
// tests compare the solution with it.
func brute(nums []int) []int {
	out := make([]int, 0, len(nums))

	for _, v := range nums {
		if v != 0 {
			out = append(out, v)
		}
	}
	for len(out) < len(nums) {
		out = append(out, 0)
	}

	return out
}

// randomNums gives n numbers between -max and max. One number out of
// zeroIn is a 0, so the arrays hold long runs of 0's and long runs of other
// numbers.
func randomNums(n, max, zeroIn int, rnd *rand.Rand) []int {
	nums := make([]int, n)
	for i := range nums {
		if rnd.Intn(zeroIn) == 0 {
			nums[i] = 0
			continue
		}
		// The value 0 is left out here, because only the branch above makes
		// a 0.
		v := 1 + rnd.Intn(max)
		if rnd.Intn(2) == 0 {
			v = -v
		}
		nums[i] = v
	}
	return nums
}

// everyArray calls f once for every array of n numbers taken from vals.
func everyArray(n int, vals []int, f func([]int)) {
	nums := make([]int, n)

	var build func(i int)
	build = func(i int) {
		if i == n {
			f(nums)
			return
		}
		for _, v := range vals {
			nums[i] = v
			build(i + 1)
		}
	}

	build(0)
}

func TestMoveZeroes(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		want []int
	}{
		{"Example1", []int{0, 1, 0, 3, 12}, []int{1, 3, 12, 0, 0}},
		{"Example2", []int{0}, []int{0}},
		{"OneNumberNotZero", []int{7}, []int{7}},
		{"NoZeroes", []int{1, 2, 3}, []int{1, 2, 3}},
		{"AllZeroes", []int{0, 0, 0}, []int{0, 0, 0}},
		{"OneZeroAtTheStart", []int{0, 1, 2}, []int{1, 2, 0}},
		{"OneZeroInTheMiddle", []int{1, 0, 2}, []int{1, 2, 0}},
		{"OneZeroAtTheEnd", []int{1, 2, 0}, []int{1, 2, 0}},
		{"TwoNumbersSwap", []int{0, 1}, []int{1, 0}},
		{"TwoNumbersAlreadyDone", []int{1, 0}, []int{1, 0}},
		{"ZeroesAtTheStart", []int{0, 0, 0, 1, 2}, []int{1, 2, 0, 0, 0}},
		{"ZeroesAtTheEnd", []int{1, 2, 0, 0, 0}, []int{1, 2, 0, 0, 0}},
		{"ZeroesInTheMiddle", []int{1, 0, 0, 0, 2}, []int{1, 2, 0, 0, 0}},
		{"EveryOtherNumberIsZero", []int{0, 1, 0, 2, 0, 3}, []int{1, 2, 3, 0, 0, 0}},
		{"EveryOtherNumberIsZeroFromOne", []int{1, 0, 2, 0, 3, 0}, []int{1, 2, 3, 0, 0, 0}},
		{"SameNumberManyTimes", []int{5, 0, 5, 0, 5}, []int{5, 5, 5, 0, 0}},
		{"NegativeNumbers", []int{0, -1, 0, -3, -12}, []int{-1, -3, -12, 0, 0}},
		{"MixedSigns", []int{-1, 0, 2, 0, -3, 0, 4}, []int{-1, 2, -3, 4, 0, 0, 0}},
		{"OrderStaysTheSame", []int{3, 0, 1, 0, 2}, []int{3, 1, 2, 0, 0}},
		{"LongRunOfZeroesFirst", []int{0, 0, 0, 0, 0, 9}, []int{9, 0, 0, 0, 0, 0}},
		{"LongRunOfZeroesLast", []int{9, 0, 0, 0, 0, 0}, []int{9, 0, 0, 0, 0, 0}},
		{"LargestAndSmallestValue", []int{0, 2147483647, 0, -2147483648}, []int{2147483647, -2147483648, 0, 0}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			nums := slices.Clone(c.nums)

			moveZeroes(nums)

			if !slices.Equal(nums, c.want) {
				t.Errorf("moveZeroes(%v) left %v, want %v", c.nums, nums, c.want)
			}
		})
	}
}

// TestMoveZeroesEveryShortArray walks every array of up to 8 numbers taken
// from 0, 1 and 2. It catches a swap that starts or stops one step too late,
// and a solution that loses a number or keeps an old copy of it.
func TestMoveZeroesEveryShortArray(t *testing.T) {
	const maxLen = 8

	vals := []int{0, 1, 2}

	for n := 1; n <= maxLen; n++ {
		failed := false

		everyArray(n, vals, func(nums []int) {
			if failed {
				return
			}

			got := slices.Clone(nums)
			moveZeroes(got)

			if want := brute(nums); !slices.Equal(got, want) {
				t.Errorf("moveZeroes(%v) left %v, want %v", nums, got, want)
				failed = true
			}
		})
	}
}

// TestMoveZeroesRandom compares the solution with brute on random arrays.
// A small zeroIn gives many 0's. A large zeroIn gives few 0's.
func TestMoveZeroesRandom(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	zeroIns := []int{2, 3, 10, 100}

	for _, zeroIn := range zeroIns {
		for range 200 {
			nums := randomNums(1+rnd.Intn(200), 1000, zeroIn, rnd)

			got := slices.Clone(nums)
			moveZeroes(got)

			if want := brute(nums); !slices.Equal(got, want) {
				t.Fatalf("moveZeroes(%v) left %v, want %v", nums, got, want)
			}
		}
	}
}

// TestMoveZeroesOneZeroEverywhere puts one 0 at every position of an array of
// different numbers. The answer holds the other numbers in the old order.
func TestMoveZeroesOneZeroEverywhere(t *testing.T) {
	const n = 10

	for at := 0; at < n; at++ {
		nums := make([]int, n)
		for i := range nums {
			nums[i] = i + 1
		}
		nums[at] = 0

		got := slices.Clone(nums)
		moveZeroes(got)

		if want := brute(nums); !slices.Equal(got, want) {
			t.Errorf("moveZeroes(%v) left %v, want %v", nums, got, want)
		}
	}
}

// TestMoveZeroesOneNumberEverywhere puts one number that is not 0 at every
// position of an array of 0's. The answer holds that number first.
func TestMoveZeroesOneNumberEverywhere(t *testing.T) {
	const n = 10

	for at := 0; at < n; at++ {
		nums := make([]int, n)
		nums[at] = 7

		got := slices.Clone(nums)
		moveZeroes(got)

		want := make([]int, n)
		want[0] = 7

		if !slices.Equal(got, want) {
			t.Errorf("moveZeroes(%v) left %v, want %v", nums, got, want)
		}
	}
}

// TestMoveZeroesInPlace checks that moveZeroes works on the array of the
// caller. A solution that builds a new slice and drops it fails here.
func TestMoveZeroesInPlace(t *testing.T) {
	nums := []int{0, 1, 0, 3, 12}
	want := []int{1, 3, 12, 0, 0}

	// The test keeps the array behind the slice, so it can read the array
	// after the call.
	var arr [5]int
	copy(arr[:], nums)

	moveZeroes(arr[:])

	if !slices.Equal(arr[:], want) {
		t.Errorf("moveZeroes left the array as %v, want %v", arr, want)
	}
}

// TestMoveZeroesKeepsNeighbours gives a part of a larger array to moveZeroes.
// A solution that writes outside of nums fails here.
func TestMoveZeroesKeepsNeighbours(t *testing.T) {
	arr := []int{111, 222, 0, 1, 0, 3, 12, 333, 444}
	nums := arr[2:7]

	moveZeroes(nums)

	if want := []int{1, 3, 12, 0, 0}; !slices.Equal(nums, want) {
		t.Errorf("moveZeroes left %v, want %v", nums, want)
	}
	if got, want := arr[:2], []int{111, 222}; !slices.Equal(got, want) {
		t.Errorf("moveZeroes changed the numbers before nums to %v, want %v", got, want)
	}
	if got, want := arr[7:], []int{333, 444}; !slices.Equal(got, want) {
		t.Errorf("moveZeroes changed the numbers after nums to %v, want %v", got, want)
	}
}

// TestMoveZeroesLarge uses the largest array of the constraints.
func TestMoveZeroesLarge(t *testing.T) {
	const n = 10000

	t.Run("AllZeroesButTheLast", func(t *testing.T) {
		// The only number that is not 0 sits at the end, so the solution must
		// read the whole array and move that number to the front.
		nums := make([]int, n)
		nums[n-1] = 5

		moveZeroes(nums)

		want := make([]int, n)
		want[0] = 5

		if !slices.Equal(nums, want) {
			t.Errorf("moveZeroes left nums[0] = %d and nums[%d] = %d, want 5 and 0", nums[0], n-1, nums[n-1])
		}
	})

	t.Run("NoZeroes", func(t *testing.T) {
		nums := make([]int, n)
		for i := range nums {
			nums[i] = i + 1
		}
		want := slices.Clone(nums)

		moveZeroes(nums)

		if !slices.Equal(nums, want) {
			t.Error("moveZeroes changed an array without a 0")
		}
	})

	t.Run("EveryOtherNumberIsZero", func(t *testing.T) {
		nums := make([]int, n)
		for i := range nums {
			if i%2 == 1 {
				nums[i] = i
			}
		}
		want := brute(nums)

		moveZeroes(nums)

		if !slices.Equal(nums, want) {
			t.Error("moveZeroes gave a wrong answer for an array of 10000 numbers")
		}
	})
}

func BenchmarkMoveZeroesSmall(b *testing.B) {
	nums := []int{0, 1, 0, 3, 12}
	work := make([]int, len(nums))

	for b.Loop() {
		copy(work, nums)
		moveZeroes(work)
	}
}

// BenchmarkMoveZeroesNoZeroes gives an array without a 0. A solution that
// writes every number back does more work here than a solution that writes
// only the numbers that move.
func BenchmarkMoveZeroesNoZeroes(b *testing.B) {
	const n = 10000

	nums := make([]int, n)
	for i := range nums {
		nums[i] = i + 1
	}
	work := make([]int, n)

	b.ResetTimer()
	for b.Loop() {
		copy(work, nums)
		moveZeroes(work)
	}
}

// BenchmarkMoveZeroesAllZeroes gives an array of 0's, so no number moves.
func BenchmarkMoveZeroesAllZeroes(b *testing.B) {
	const n = 10000

	work := make([]int, n)

	b.ResetTimer()
	for b.Loop() {
		clear(work)
		moveZeroes(work)
	}
}

// BenchmarkMoveZeroesRandom gives random numbers, so the 0's sit in many
// different places.
func BenchmarkMoveZeroesRandom(b *testing.B) {
	const n = 10000

	nums := randomNums(n, 1000, 3, rand.New(rand.NewSource(2)))
	work := make([]int, n)

	b.ResetTimer()
	for b.Loop() {
		copy(work, nums)
		moveZeroes(work)
	}
}
