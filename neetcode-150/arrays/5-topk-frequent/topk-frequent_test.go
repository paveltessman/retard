package topkfrequent

import (
	"math/rand"
	"slices"
	"testing"
)

// brute gives the answer of the task in the slow and plain way. It counts
// every distinct value with a full pass over nums. Then it takes the value
// with the highest count k times, and skips the values that it took before.
// The tests compare the solution with it.
func brute(nums []int, k int) []int {
	var values, counts []int
	for _, v := range nums {
		if slices.Contains(values, v) {
			continue
		}
		c := 0
		for _, w := range nums {
			if w == v {
				c++
			}
		}
		values = append(values, v)
		counts = append(counts, c)
	}

	taken := make([]bool, len(values))
	var res []int
	for range k {
		best := -1
		for i := range values {
			if !taken[i] && (best == -1 || counts[i] > counts[best]) {
				best = i
			}
		}
		taken[best] = true
		res = append(res, values[best])
	}

	return res
}

// validKs gives every k for which nums has exactly one answer. The task
// promises such a k, so the tests use only these.
func validKs(nums []int) []int {
	freq := map[int]int{}
	for _, v := range nums {
		freq[v]++
	}
	counts := make([]int, 0, len(freq))
	for _, c := range freq {
		counts = append(counts, c)
	}
	slices.Sort(counts)
	slices.Reverse(counts)

	var ks []int
	for k := 1; k <= len(counts); k++ {
		if k == len(counts) || counts[k-1] > counts[k] {
			ks = append(ks, k)
		}
	}

	return ks
}

// check calls the solution on a copy of nums and compares the answer with
// want, in any order. It also reads the copy back, so it catches a solution
// that sorts or moves the input. It gives back false on an error.
func check(t *testing.T, nums []int, k int, want []int) bool {
	t.Helper()

	in := slices.Clone(nums)
	got := topKFrequent(in, k)
	ok := true

	sortedGot := slices.Sorted(slices.Values(got))
	sortedWant := slices.Sorted(slices.Values(want))
	if !slices.Equal(sortedGot, sortedWant) {
		t.Errorf("topKFrequent(%v, %d) = %v, want %v (in any order)", nums, k, got, want)
		ok = false
	}
	if !slices.Equal(in, nums) {
		t.Errorf("topKFrequent(%v, %d) changed the array to %v", nums, k, in)
		ok = false
	}

	return ok
}

// repeat gives an array that holds every value of values, each one count
// times, in the order of values.
func repeat(values []int, count int) []int {
	nums := make([]int, 0, len(values)*count)
	for _, v := range values {
		for range count {
			nums = append(nums, v)
		}
	}
	return nums
}

func TestTopKFrequent(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		k    int
		want []int
	}{
		{"Example1", []int{1, 2, 2, 3, 3, 3}, 2, []int{2, 3}},
		{"Example2", []int{7, 7}, 1, []int{7}},
		{"OneNumber", []int{5}, 1, []int{5}},
		{"TwoDifferentTakeBoth", []int{5, 6}, 2, []int{5, 6}},
		{"TwoNumbersTakeOne", []int{4, 6, 6}, 1, []int{6}},
		{"MostFrequentFirst", []int{9, 9, 9, 1, 2, 2}, 1, []int{9}},
		{"MostFrequentLast", []int{1, 2, 2, 9, 9, 9}, 1, []int{9}},
		{"MostFrequentInTheMiddle", []int{1, 2, 9, 9, 9, 3}, 1, []int{9}},
		{"MostFrequentSpread", []int{9, 1, 9, 2, 9, 3}, 1, []int{9}},
		{"RareValueComesFirst", []int{1, 2, 3, 3, 4, 4, 4}, 2, []int{3, 4}},
		{"TakeAll", []int{3, 1, 2, 1}, 3, []int{1, 2, 3}},
		{"TakeAllWithEqualCounts", []int{4, 8, 6}, 3, []int{4, 6, 8}},
		{"TieInsideTheAnswer", []int{1, 1, 2, 2, 3}, 2, []int{1, 2}},
		{"TieInsideTheAnswerWithAWinner", []int{5, 5, 5, 1, 1, 2, 2, 3}, 3, []int{1, 2, 5}},
		{"TakeAllButOne", []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4}, 3, []int{2, 3, 4}},
		{"AllSame", []int{8, 8, 8, 8}, 1, []int{8}},
		{"Zero", []int{0, 0, 1}, 1, []int{0}},
		{"Negatives", []int{-1, -2, -2, -3, -3, -3}, 2, []int{-2, -3}},
		{"SignPair", []int{-5, 5, 5}, 1, []int{5}},
		{"SignPairNegativeWins", []int{-5, -5, 5}, 1, []int{-5}},
		{"LowLimit", []int{-1000, -1000, 999}, 1, []int{-1000}},
		{"HighLimit", []int{1000, 1000, -999}, 1, []int{1000}},
		{"BothLimits", []int{-1000, 1000, 1000, 0}, 1, []int{1000}},
		{"BothLimitsTakeAll", []int{-1000, 1000}, 2, []int{-1000, 1000}},
		{"CountMatchesValue", []int{1, 2, 2, 3, 3, 3, 7}, 1, []int{3}},
		{"CountsOrderedAgainstValues", []int{1, 1, 1, 1, 2, 2, 2, 3, 3, 4}, 2, []int{1, 2}},
		{"LongRunAtTheEnd", []int{1, 2, 3, 4, 5, 6, 6, 6}, 1, []int{6}},
		{"Interleaved", []int{1, 2, 1, 2, 1, 3, 2, 1}, 2, []int{1, 2}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			check(t, c.nums, c.k, c.want)
		})
	}
}

// TestTopKFrequentEveryShort walks every array of up to 7 values from a small
// value set, and every k with exactly one answer. It catches an answer that
// holds one value too many or one too few, and a count that is off by one.
func TestTopKFrequentEveryShort(t *testing.T) {
	sets := []struct {
		values []int
		maxLen int
	}{
		{[]int{0, 1, 2}, 7},
		{[]int{-1000, -1, 0, 1000}, 6},
	}

	for _, s := range sets {
		for n := 1; n <= s.maxLen; n++ {
			idx := make([]int, n)
			for {
				nums := make([]int, n)
				for i, j := range idx {
					nums[i] = s.values[j]
				}
				for _, k := range validKs(nums) {
					if !check(t, nums, k, brute(nums, k)) {
						return
					}
				}

				// Go to the next array, like a counter in base len(values).
				i := 0
				for i < n && idx[i] == len(s.values)-1 {
					idx[i] = 0
					i++
				}
				if i == n {
					break
				}
				idx[i]++
			}
		}
	}
}

// TestTopKFrequentRandom compares the solution with brute on random arrays.
// A narrow value range gives long runs and close counts. A wide range gives
// many values that appear once.
func TestTopKFrequentRandom(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))

	shapes := []struct {
		name   string
		lo, hi int
	}{
		{"Narrow", 0, 4},
		{"Medium", -50, 50},
		{"Wide", -1000, 1000},
	}

	for _, s := range shapes {
		t.Run(s.name, func(t *testing.T) {
			for range 300 {
				n := 1 + rnd.Intn(300)
				nums := make([]int, n)
				for i := range nums {
					nums[i] = s.lo + rnd.Intn(s.hi-s.lo+1)
				}
				ks := validKs(nums)
				k := ks[rnd.Intn(len(ks))]

				if !check(t, nums, k, brute(nums, k)) {
					return
				}
			}
		})
	}
}

// TestTopKFrequentOrderMeansNothing shuffles the same array many times.
// The counts stay the same, so every shuffle gives the same answer. It catches
// a solution that depends on the place of a value in the array.
func TestTopKFrequentOrderMeansNothing(t *testing.T) {
	rnd := rand.New(rand.NewSource(2))

	nums := []int{4, 4, 4, 4, -2, -2, -2, 7, 7, 0, 9}
	want := []int{4, -2, 7}

	for range 100 {
		rnd.Shuffle(len(nums), func(i, j int) {
			nums[i], nums[j] = nums[j], nums[i]
		})
		if !check(t, nums, 3, want) {
			return
		}
	}
}

// TestTopKFrequentLarge builds arrays of the largest size that the
// constraints allow. The construction gives the answer, so the test does not
// need brute. The errors do not print the array, because it is too long.
func TestTopKFrequentLarge(t *testing.T) {
	rnd := rand.New(rand.NewSource(3))

	run := func(t *testing.T, nums []int, k int, want []int) {
		t.Helper()

		in := slices.Clone(nums)
		got := slices.Sorted(slices.Values(topKFrequent(in, k)))
		want = slices.Sorted(slices.Values(want))

		if !slices.Equal(got, want) {
			if len(got) > 20 || len(want) > 20 {
				t.Errorf("topKFrequent of %d numbers, k = %d, gave %d values, want %d values", len(nums), k, len(got), len(want))
			} else {
				t.Errorf("topKFrequent of %d numbers, k = %d, = %v, want %v", len(nums), k, got, want)
			}
		}
		if !slices.Equal(in, nums) {
			t.Errorf("topKFrequent of %d numbers changed the array", len(nums))
		}
	}

	// Every value from -1000 to 1000 appears once. Then 120 random values
	// get extra copies: the value at place i of top gets 120-i more. So the
	// counts of the top 120 values are all different, and every k up to 120
	// has exactly one answer.
	t.Run("Steps", func(t *testing.T) {
		all := make([]int, 0, 2001)
		for v := -1000; v <= 1000; v++ {
			all = append(all, v)
		}
		rnd.Shuffle(len(all), func(i, j int) { all[i], all[j] = all[j], all[i] })
		top := all[:120]

		nums := slices.Clone(all)
		for i, v := range top {
			for range 120 - i {
				nums = append(nums, v)
			}
		}
		rnd.Shuffle(len(nums), func(i, j int) { nums[i], nums[j] = nums[j], nums[i] })

		for _, k := range []int{1, 2, 10, 64, 119, 120} {
			run(t, nums, k, top[:k])
		}
	})

	// Every value from -1000 to 1000 appears 4 times, so k must be 2001.
	t.Run("EveryValueTakeAll", func(t *testing.T) {
		values := make([]int, 0, 2001)
		for v := -1000; v <= 1000; v++ {
			values = append(values, v)
		}
		nums := repeat(values, 4)
		rnd.Shuffle(len(nums), func(i, j int) { nums[i], nums[j] = nums[j], nums[i] })

		run(t, nums, len(values), values)
	})

	// One value fills the whole array.
	for _, v := range []int{-1000, 0, 1000} {
		t.Run("AllSame", func(t *testing.T) {
			run(t, repeat([]int{v}, 10000), 1, []int{v})
		})
	}

	// Two values at the limits, one more often than the other.
	t.Run("TwoLimits", func(t *testing.T) {
		nums := append(repeat([]int{1000}, 4999), repeat([]int{-1000}, 5001)...)
		run(t, nums, 1, []int{-1000})
		run(t, nums, 2, []int{-1000, 1000})
	})
}

func BenchmarkTopKFrequentExample(b *testing.B) {
	nums := []int{1, 2, 2, 3, 3, 3}

	for b.Loop() {
		topKFrequent(nums, 2)
	}
}

// BenchmarkTopKFrequentWide uses the largest array with values over the whole
// range, and asks for many of them.
func BenchmarkTopKFrequentWide(b *testing.B) {
	rnd := rand.New(rand.NewSource(4))
	nums := make([]int, 10000)
	for i := range nums {
		nums[i] = rnd.Intn(2001) - 1000
	}
	ks := validKs(nums)
	k := ks[len(ks)/2]

	for b.Loop() {
		topKFrequent(nums, k)
	}
}

// BenchmarkTopKFrequentAllSame uses the largest array of one value.
func BenchmarkTopKFrequentAllSame(b *testing.B) {
	nums := repeat([]int{7}, 10000)

	for b.Loop() {
		topKFrequent(nums, 1)
	}
}
