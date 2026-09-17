package removeduplicates

import (
	"math/rand"
	"slices"
	"testing"
)

// brute gives the answer of the task in the slow and plain way: it walks the
// sorted numbers and keeps a number only if it is not the same as the number
// before it. The tests compare the solution with it.
func brute(nums []int) []int {
	out := make([]int, 0, len(nums))

	for i, v := range nums {
		if i == 0 || v != nums[i-1] {
			out = append(out, v)
		}
	}

	return out
}

// check calls removeDuplicates on a copy of nums and compares the first k
// numbers with want. The numbers after k - 1 are free, so check reads none of
// them.
func check(t *testing.T, nums, want []int) {
	t.Helper()

	got := slices.Clone(nums)

	k := removeDuplicates(got)

	if k != len(want) {
		t.Errorf("removeDuplicates(%v) = %d, want %d", nums, k, len(want))
		return
	}
	if !slices.Equal(got[:k], want) {
		t.Errorf("removeDuplicates(%v) left %v, want %v", nums, got[:k], want)
	}
}

// randomNums gives n sorted numbers between -max and max. A small step gives
// long runs of the same number. A large step gives few duplicates.
func randomNums(n, max, step int, rnd *rand.Rand) []int {
	nums := make([]int, n)

	v := -max
	for i := range nums {
		nums[i] = v
		v += rnd.Intn(step)
		if v > max {
			v = max
		}
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

func TestRemoveDuplicates(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		want []int
	}{
		{"Example1", []int{1, 1, 2}, []int{1, 2}},
		{"Example2", []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}, []int{0, 1, 2, 3, 4}},
		{"OneNumber", []int{1}, []int{1}},
		{"TwoNumbersSame", []int{1, 1}, []int{1}},
		{"TwoNumbersDifferent", []int{1, 2}, []int{1, 2}},
		{"NoDuplicates", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
		{"AllTheSame", []int{7, 7, 7, 7}, []int{7}},
		{"DuplicatesAtTheStart", []int{1, 1, 1, 2, 3}, []int{1, 2, 3}},
		{"DuplicatesInTheMiddle", []int{1, 2, 2, 2, 3}, []int{1, 2, 3}},
		{"DuplicatesAtTheEnd", []int{1, 2, 3, 3, 3}, []int{1, 2, 3}},
		{"EveryNumberTwice", []int{1, 1, 2, 2, 3, 3}, []int{1, 2, 3}},
		{"EveryNumberThreeTimes", []int{1, 1, 1, 2, 2, 2}, []int{1, 2}},
		{"LongRunThenOneNumber", []int{1, 1, 1, 1, 1, 2}, []int{1, 2}},
		{"OneNumberThenLongRun", []int{1, 2, 2, 2, 2, 2}, []int{1, 2}},
		{"NegativeNumbers", []int{-3, -3, -2, -1, -1}, []int{-3, -2, -1}},
		{"MixedSigns", []int{-2, -2, -1, 0, 0, 1, 1}, []int{-2, -1, 0, 1}},
		{"WithZero", []int{0, 0, 0}, []int{0}},
		{"SmallestValue", []int{-100, -100, -99}, []int{-100, -99}},
		{"LargestValue", []int{99, 100, 100}, []int{99, 100}},
		{"WholeRangeOfValues", []int{-100, -100, 0, 0, 100, 100}, []int{-100, 0, 100}},
		{"NeighbourValues", []int{1, 2, 3}, []int{1, 2, 3}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			check(t, c.nums, c.want)
		})
	}
}

// TestRemoveDuplicatesEveryShortArray walks every sorted array of up to 9
// numbers taken from 0, 1 and 2. It catches a write that starts or stops one
// step too late, and a solution that drops a unique number or keeps a copy.
func TestRemoveDuplicatesEveryShortArray(t *testing.T) {
	const maxLen = 9

	vals := []int{0, 1, 2}

	for n := 1; n <= maxLen; n++ {
		failed := false

		everyArray(n, vals, func(nums []int) {
			if failed || !slices.IsSorted(nums) {
				return
			}

			got := slices.Clone(nums)
			k := removeDuplicates(got)

			want := brute(nums)
			if k != len(want) {
				t.Errorf("removeDuplicates(%v) = %d, want %d", nums, k, len(want))
				failed = true
				return
			}
			if !slices.Equal(got[:k], want) {
				t.Errorf("removeDuplicates(%v) left %v, want %v", nums, got[:k], want)
				failed = true
			}
		})
	}
}

// TestRemoveDuplicatesRandom compares the solution with brute on random sorted
// arrays. A small step gives many duplicates. A large step gives few.
func TestRemoveDuplicatesRandom(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	steps := []int{2, 3, 10, 100}

	for _, step := range steps {
		for range 200 {
			nums := randomNums(1+rnd.Intn(200), 100, step, rnd)

			got := slices.Clone(nums)
			k := removeDuplicates(got)

			want := brute(nums)
			if k != len(want) {
				t.Fatalf("removeDuplicates(%v) = %d, want %d", nums, k, len(want))
			}
			if !slices.Equal(got[:k], want) {
				t.Fatalf("removeDuplicates(%v) left %v, want %v", nums, got[:k], want)
			}
		}
	}
}

// TestRemoveDuplicatesOneRunEverywhere puts one run of two of the same number
// at every position of an array of different numbers.
func TestRemoveDuplicatesOneRunEverywhere(t *testing.T) {
	const n = 10

	for at := 0; at < n; at++ {
		nums := make([]int, 0, n+1)
		for i := 0; i < n; i++ {
			nums = append(nums, i)
			if i == at {
				nums = append(nums, i)
			}
		}

		check(t, nums, brute(nums))
	}
}

// TestRemoveDuplicatesRunLengths gives one number many times, then one other
// number. The run grows with every step, so the write point moves further away
// from the read point.
func TestRemoveDuplicatesRunLengths(t *testing.T) {
	const maxRun = 20

	for run := 1; run <= maxRun; run++ {
		nums := make([]int, run, run+1)
		for i := range nums {
			nums[i] = 1
		}
		nums = append(nums, 2)

		check(t, nums, []int{1, 2})
	}
}

// TestRemoveDuplicatesInPlace checks that removeDuplicates works on the array
// of the caller. A solution that builds a new slice and drops it fails here.
func TestRemoveDuplicatesInPlace(t *testing.T) {
	var arr [3]int
	copy(arr[:], []int{1, 1, 2})

	k := removeDuplicates(arr[:])

	if k != 2 {
		t.Fatalf("removeDuplicates([1 1 2]) = %d, want 2", k)
	}
	if want := []int{1, 2}; !slices.Equal(arr[:k], want) {
		t.Errorf("removeDuplicates left the array as %v, want %v first", arr, want)
	}
}

// TestRemoveDuplicatesKeepsNeighbours gives a part of a larger array to
// removeDuplicates. A solution that writes outside of nums fails here.
func TestRemoveDuplicatesKeepsNeighbours(t *testing.T) {
	arr := []int{111, 222, 1, 1, 2, 2, 3, 333, 444}
	nums := arr[2:7]

	k := removeDuplicates(nums)

	if k != 3 {
		t.Fatalf("removeDuplicates(%v) = %d, want 3", []int{1, 1, 2, 2, 3}, k)
	}
	if want := []int{1, 2, 3}; !slices.Equal(nums[:k], want) {
		t.Errorf("removeDuplicates left %v, want %v", nums[:k], want)
	}
	if got, want := arr[:2], []int{111, 222}; !slices.Equal(got, want) {
		t.Errorf("removeDuplicates changed the numbers before nums to %v, want %v", got, want)
	}
	if got, want := arr[7:], []int{333, 444}; !slices.Equal(got, want) {
		t.Errorf("removeDuplicates changed the numbers after nums to %v, want %v", got, want)
	}
}

// TestRemoveDuplicatesLarge uses the largest array of the constraints.
func TestRemoveDuplicatesLarge(t *testing.T) {
	const n = 30000

	t.Run("AllTheSame", func(t *testing.T) {
		nums := make([]int, n)

		k := removeDuplicates(nums)

		if k != 1 {
			t.Fatalf("removeDuplicates of 30000 zeroes = %d, want 1", k)
		}
		if nums[0] != 0 {
			t.Errorf("removeDuplicates left nums[0] = %d, want 0", nums[0])
		}
	})

	t.Run("EveryValueOfTheRange", func(t *testing.T) {
		// The values run from -100 to 100, so every value comes back about 150
		// times.
		nums := make([]int, n)
		for i := range nums {
			nums[i] = -100 + i*201/n
		}
		want := brute(nums)

		k := removeDuplicates(nums)

		if k != len(want) {
			t.Fatalf("removeDuplicates = %d, want %d", k, len(want))
		}
		if !slices.Equal(nums[:k], want) {
			t.Error("removeDuplicates gave a wrong answer for an array of 30000 numbers")
		}
	})

	t.Run("NoDuplicates", func(t *testing.T) {
		nums := make([]int, n)
		for i := range nums {
			nums[i] = i
		}
		want := slices.Clone(nums)

		k := removeDuplicates(nums)

		if k != n {
			t.Fatalf("removeDuplicates of 30000 different numbers = %d, want %d", k, n)
		}
		if !slices.Equal(nums, want) {
			t.Error("removeDuplicates changed an array without a duplicate")
		}
	})
}

func BenchmarkRemoveDuplicatesSmall(b *testing.B) {
	nums := []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}
	work := make([]int, len(nums))

	for b.Loop() {
		copy(work, nums)
		removeDuplicates(work)
	}
}

// BenchmarkRemoveDuplicatesNoDuplicates gives an array without a duplicate. A
// solution that writes every number back does more work here than a solution
// that writes only the numbers that move.
func BenchmarkRemoveDuplicatesNoDuplicates(b *testing.B) {
	const n = 30000

	nums := make([]int, n)
	for i := range nums {
		nums[i] = i
	}
	work := make([]int, n)

	b.ResetTimer()
	for b.Loop() {
		copy(work, nums)
		removeDuplicates(work)
	}
}

// BenchmarkRemoveDuplicatesAllTheSame gives an array of one number, so the
// answer holds one number.
func BenchmarkRemoveDuplicatesAllTheSame(b *testing.B) {
	const n = 30000

	work := make([]int, n)

	b.ResetTimer()
	for b.Loop() {
		clear(work)
		removeDuplicates(work)
	}
}

// BenchmarkRemoveDuplicatesRandom gives random sorted numbers, so the runs
// have different lengths.
func BenchmarkRemoveDuplicatesRandom(b *testing.B) {
	const n = 30000

	nums := randomNums(n, 100, 3, rand.New(rand.NewSource(2)))
	work := make([]int, n)

	b.ResetTimer()
	for b.Loop() {
		copy(work, nums)
		removeDuplicates(work)
	}
}
