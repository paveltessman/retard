package subarraysum

import (
	"math/rand"
	"slices"
	"testing"
)

// brute gives the answer of the task in the slow and plain way: it starts a
// subarray at every position, adds the numbers one by one, and counts every
// sum that reaches k. The tests compare the solution with it.
func brute(nums []int, k int) int {
	count := 0

	for i := range nums {
		sum := 0
		for j := i; j < len(nums); j++ {
			sum += nums[j]
			if sum == k {
				count++
			}
		}
	}

	return count
}

// check calls subarraySum on a copy of nums and compares the answer with want.
// It also reads the copy back, because subarraySum must leave the array as it
// was.
func check(t *testing.T, nums []int, k, want int) {
	t.Helper()

	got := slices.Clone(nums)

	if n := subarraySum(got, k); n != want {
		t.Errorf("subarraySum(%v, %d) = %d, want %d", nums, k, n, want)
	}
	if !slices.Equal(got, nums) {
		t.Errorf("subarraySum(%v, %d) changed the array to %v", nums, k, got)
	}
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

// randomNums gives n numbers between -max and max.
func randomNums(n, max int, rnd *rand.Rand) []int {
	nums := make([]int, n)
	for i := range nums {
		nums[i] = rnd.Intn(2*max+1) - max
	}

	return nums
}

func TestSubarraySum(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		k    int
		want int
	}{
		{"Example1", []int{1, 1, 1}, 2, 2},
		{"Example2", []int{1, 2, 3}, 3, 2},
		{"OneNumberMatches", []int{1}, 1, 1},
		{"OneNumberDoesNot", []int{1}, 2, 0},
		{"OneZeroAndKZero", []int{0}, 0, 1},
		{"OneNegative", []int{-1}, -1, 1},
		{"AllZerosAndKZero", []int{0, 0, 0}, 0, 6},
		{"AllZerosAndKOne", []int{0, 0, 0}, 1, 0},
		{"ZerosAroundOne", []int{0, 1, 0}, 1, 4},
		{"NegativeK", []int{-1, -1, -1}, -2, 2},
		{"MixedSignsAndKZero", []int{1, -1, 1, -1}, 0, 4},
		{"PrefixRepeats", []int{1, -1, 0}, 0, 3},
		{"WholeArray", []int{1, 2, 3, 4}, 10, 1},
		{"NoSubarray", []int{1, 2, 3}, 7, 0},
		{"SubarrayStartsAtTheFirstNumber", []int{3, 4, 7}, 3, 1},
		{"SubarrayEndsAtTheLastNumber", []int{3, 4, 7}, 7, 2},
		{"EmptySubarrayDoesNotCount", []int{1, 2, 3}, 0, 0},
		{"LargeNumbers", []int{1000, 1000, 1000}, 3000, 1},
		{"KAboveTheSum", []int{1000, 1000}, 10000000, 0},
		{"KBelowTheSum", []int{-1000, -1000}, -10000000, 0},
		{"Alternating", []int{1, -1, 1, -1, 1}, 1, 6},
		{"EqualNumbers", []int{2, 2, 2, 2}, 4, 3},
		{"KIsOneNumber", []int{5, 5, 5}, 5, 3},
		{"GrowingNumbers", []int{1, 2, 3, 4, 5}, 5, 2},
		{"NegativeAndPositive", []int{-1, 2, -1, 2, -1}, 1, 5},
		{"LongRun", []int{1, 1, 1, 1}, 2, 3},
		{"KIsTheSumOfAll", []int{1, -1, 1, -1, 2}, 2, 3},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			check(t, c.nums, c.k, c.want)
		})
	}
}

// TestSubarraySumEveryShortArray walks every array of up to 6 numbers taken
// from -1, 0, 1 and 2, with every k from -4 to 4. It catches a solution that
// misses the subarrays that start at the first number, a solution that counts
// a prefix sum once when it comes several times, and a solution that counts
// the empty subarray.
func TestSubarraySumEveryShortArray(t *testing.T) {
	const maxLen = 6

	vals := []int{-1, 0, 1, 2}

	for n := 0; n <= maxLen; n++ {
		failed := false

		everyArray(n, vals, func(nums []int) {
			if failed {
				return
			}

			for k := -4; k <= 4; k++ {
				got := subarraySum(slices.Clone(nums), k)

				if want := brute(nums, k); got != want {
					t.Errorf("subarraySum(%v, %d) = %d, want %d", nums, k, got, want)
					failed = true
					return
				}
			}
		})
	}
}

// TestSubarraySumRandom compares the solution with brute on random arrays. A
// small range of numbers gives many equal prefix sums. A large range gives
// few, so most answers are 0 and a wrong count shows up on the other ones.
func TestSubarraySumRandom(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	maxVals := []int{1, 2, 5, 50, 1000}

	for _, max := range maxVals {
		for range 200 {
			nums := randomNums(1+rnd.Intn(200), max, rnd)

			// A small k keeps the answer above 0 for a part of the cases.
			k := rnd.Intn(4*max+1) - 2*max

			got := subarraySum(slices.Clone(nums), k)

			if want := brute(nums, k); got != want {
				t.Fatalf("subarraySum(%v, %d) = %d, want %d", nums, k, got, want)
			}
		}
	}
}

// TestSubarraySumPositiveRandom compares the solution with brute on arrays of
// positive numbers only. Here the prefix sums grow, so every prefix sum comes
// once and the answer stays small.
func TestSubarraySumPositiveRandom(t *testing.T) {
	rnd := rand.New(rand.NewSource(2))

	for range 200 {
		n := 1 + rnd.Intn(200)

		nums := make([]int, n)
		for i := range nums {
			nums[i] = 1 + rnd.Intn(10)
		}

		sum := 0
		for _, v := range nums {
			sum += v
		}

		k := rnd.Intn(sum + 1)

		got := subarraySum(slices.Clone(nums), k)

		if want := brute(nums, k); got != want {
			t.Fatalf("subarraySum(%v, %d) = %d, want %d", nums, k, got, want)
		}
	}
}

// TestSubarraySumEveryK uses one array and walks every k from below the
// smallest sum to above the largest sum. The sum of every k is the number of
// subarrays of the array, so a lost or a double count shows up.
func TestSubarraySumEveryK(t *testing.T) {
	arrays := [][]int{
		{1, 1, 1, 1, 1, 1},
		{0, 0, 0, 0, 0},
		{1, -1, 1, -1, 1, -1},
		{3, -2, 5, -1, 4, -6, 2},
		{-1, -2, -3, -4},
		{1000, -1000, 1000, -1000},
	}

	for _, nums := range arrays {
		low, high := 0, 0
		for _, v := range nums {
			if v < 0 {
				low += v
			} else {
				high += v
			}
		}

		total := 0

		for k := low - 2; k <= high+2; k++ {
			got := subarraySum(slices.Clone(nums), k)

			if want := brute(nums, k); got != want {
				t.Errorf("subarraySum(%v, %d) = %d, want %d", nums, k, got, want)
			}

			total += got
		}

		n := len(nums)
		if want := n * (n + 1) / 2; total != want {
			t.Errorf("subarraySum(%v, every k) counted %d subarrays, want %d", nums, total, want)
		}
	}
}

// TestSubarraySumOneMatchEverywhere puts the one subarray that reaches k at
// every position of an array of zeros. The numbers around it are 0, so a
// solution that drops a prefix sum too early misses the match.
func TestSubarraySumOneMatchEverywhere(t *testing.T) {
	const n = 12

	for at := 0; at < n; at++ {
		nums := make([]int, n)
		nums[at] = 7

		// Every subarray that holds the 7 and no other number reaches k. The
		// zeros to the left and to the right can join it, so the count is the
		// number of ways to take them.
		want := (at + 1) * (n - at)

		check(t, nums, 7, want)
	}
}

// TestSubarraySumKeepsNeighbours gives a part of a larger array to
// subarraySum. A solution that reads outside of nums fails here.
func TestSubarraySumKeepsNeighbours(t *testing.T) {
	arr := []int{1000, 1000, 1, 2, 3, 1000, 1000}
	nums := arr[2:5]

	if got, want := subarraySum(nums, 3), 2; got != want {
		t.Errorf("subarraySum(%v, 3) = %d, want %d", nums, got, want)
	}
}

// TestSubarraySumLarge uses the largest array of the constraints. A solution
// that tries every subarray needs much more time here.
func TestSubarraySumLarge(t *testing.T) {
	const n = 20000

	t.Run("AllZeros", func(t *testing.T) {
		// Every subarray of the array reaches 0, so the answer is the number
		// of subarrays.
		nums := make([]int, n)

		if got, want := subarraySum(nums, 0), n*(n+1)/2; got != want {
			t.Errorf("subarraySum(%d zeros, 0) = %d, want %d", n, got, want)
		}
	})

	t.Run("AllZerosAndKOne", func(t *testing.T) {
		nums := make([]int, n)

		if got := subarraySum(nums, 1); got != 0 {
			t.Errorf("subarraySum(%d zeros, 1) = %d, want 0", n, got)
		}
	})

	t.Run("AllOnes", func(t *testing.T) {
		// A subarray reaches 100 when it holds 100 ones, and such a subarray
		// starts at every position but the last 99.
		nums := make([]int, n)
		for i := range nums {
			nums[i] = 1
		}

		if got, want := subarraySum(nums, 100), n-99; got != want {
			t.Errorf("subarraySum(%d ones, 100) = %d, want %d", n, got, want)
		}
	})

	t.Run("Alternating", func(t *testing.T) {
		// The prefix sums are 0 and 1 only, so each of them comes about half
		// of the time and the answer is large.
		nums := make([]int, n)
		for i := range nums {
			if i%2 == 0 {
				nums[i] = 1
			} else {
				nums[i] = -1
			}
		}

		// The prefix sum is 1 after an odd number of numbers, and 0 after an
		// even one. A subarray reaches 0 when both of its ends stand at the
		// same kind of position.
		half := n/2 + 1
		want := half*(half-1)/2 + (n+1-half)*(n-half)/2

		if got := subarraySum(nums, 0); got != want {
			t.Errorf("subarraySum(%d alternating numbers, 0) = %d, want %d", n, got, want)
		}
	})

	t.Run("LargestK", func(t *testing.T) {
		// The sum of the array is 2 * 10^7, above the largest k of the
		// constraints. Only the subarrays of 10000 numbers reach k.
		nums := make([]int, n)
		for i := range nums {
			nums[i] = 1000
		}

		if got, want := subarraySum(nums, 10000000), n-9999; got != want {
			t.Errorf("subarraySum(%d numbers of 1000, 10000000) = %d, want %d", n, got, want)
		}
	})

	t.Run("GrowingPrefixSums", func(t *testing.T) {
		// Every prefix sum comes once, so the map holds n + 1 numbers and the
		// answer is 0.
		nums := make([]int, n)
		for i := range nums {
			nums[i] = 1 + i%1000
		}

		if got := subarraySum(nums, -1); got != 0 {
			t.Errorf("subarraySum(%d growing numbers, -1) = %d, want 0", n, got)
		}
	})
}

// TestSubarraySumMedium compares the solution with brute on an array that
// brute still handles. It joins a run of zeros, a run of equal numbers and a
// random part.
func TestSubarraySumMedium(t *testing.T) {
	const n = 2000

	rnd := rand.New(rand.NewSource(3))

	nums := make([]int, 0, n)
	for range n / 4 {
		nums = append(nums, 0)
	}
	for range n / 4 {
		nums = append(nums, 3)
	}
	for i := range n / 4 {
		if i%2 == 0 {
			nums = append(nums, 5)
		} else {
			nums = append(nums, -5)
		}
	}
	nums = append(nums, randomNums(n/4, 4, rnd)...)

	for _, k := range []int{-5, 0, 3, 6, 15, 1000} {
		got := subarraySum(slices.Clone(nums), k)

		if want := brute(nums, k); got != want {
			t.Errorf("subarraySum(the mixed array, %d) = %d, want %d", k, got, want)
		}
	}
}

func BenchmarkSubarraySumSmall(b *testing.B) {
	nums := []int{1, 2, 3}

	for b.Loop() {
		subarraySum(nums, 3)
	}
}

// BenchmarkSubarraySumRandom gives the largest array of the constraints.
func BenchmarkSubarraySumRandom(b *testing.B) {
	const n = 20000

	nums := randomNums(n, 1000, rand.New(rand.NewSource(4)))

	b.ResetTimer()
	for b.Loop() {
		subarraySum(nums, 0)
	}
}

// BenchmarkSubarraySumAllZeros gives one prefix sum many times, so every step
// finds the number in the map.
func BenchmarkSubarraySumAllZeros(b *testing.B) {
	const n = 20000

	nums := make([]int, n)

	b.ResetTimer()
	for b.Loop() {
		subarraySum(nums, 0)
	}
}

// BenchmarkSubarraySumGrowing gives numbers of one sign, so every prefix sum
// comes once and the map grows to the length of the array.
func BenchmarkSubarraySumGrowing(b *testing.B) {
	const n = 20000

	nums := make([]int, n)
	for i := range nums {
		nums[i] = 1 + i%1000
	}

	b.ResetTimer()
	for b.Loop() {
		subarraySum(nums, 500000)
	}
}
