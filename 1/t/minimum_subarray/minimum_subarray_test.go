package minimumsubarray

import (
	"math/rand"
	"slices"
	"testing"
)

// brute gives the answer of the task in the slow and plain way: it starts a
// subarray at every position, adds the numbers one by one, and stops as soon
// as the sum reaches the target. The tests compare the solution with it.
func brute(target int, nums []int) int {
	best := 0

	for i := range len(nums) {
		sum := 0
		for j := i; j < len(nums); j++ {
			sum += nums[j]
			if sum >= target {
				if best == 0 || j-i+1 < best {
					best = j - i + 1
				}
				break
			}
		}
	}

	return best
}

// randomNums gives n numbers between 1 and max.
func randomNums(n, max int, rnd *rand.Rand) []int {
	nums := make([]int, n)
	for i := range nums {
		nums[i] = 1 + rnd.Intn(max)
	}
	return nums
}

// everyArray calls f once for every array of n numbers between 1 and max.
func everyArray(n, max int, f func([]int)) {
	nums := make([]int, n)

	var build func(i int)
	build = func(i int) {
		if i == n {
			f(nums)
			return
		}
		for v := 1; v <= max; v++ {
			nums[i] = v
			build(i + 1)
		}
	}

	build(0)
}

func TestMinSubArrayLen(t *testing.T) {
	cases := []struct {
		name   string
		target int
		nums   []int
		want   int
	}{
		{"Example1", 7, []int{2, 3, 1, 2, 4, 3}, 2},
		{"Example2", 4, []int{1, 4, 4}, 1},
		{"Example3", 11, []int{1, 1, 1, 1, 1, 1, 1, 1}, 0},
		{"OneNumberEnough", 5, []int{5}, 1},
		{"OneNumberTooSmall", 5, []int{4}, 0},
		{"OneNumberAbove", 5, []int{9}, 1},
		{"TwoNumbers", 3, []int{1, 2}, 2},
		{"TwoNumbersTooSmall", 4, []int{1, 2}, 0},
		{"WholeArray", 8, []int{1, 2, 2, 3}, 4},
		{"WholeArrayOneShort", 9, []int{1, 2, 2, 3}, 0},
		{"SumIsExactlyTheTarget", 7, []int{1, 2, 4}, 3},
		{"FirstNumber", 3, []int{3, 1, 1}, 1},
		{"LastNumber", 3, []int{1, 1, 3}, 1},
		{"BestWindowFirst", 6, []int{3, 4, 1, 1, 1, 1}, 2},
		{"BestWindowLast", 6, []int{1, 1, 1, 1, 3, 4}, 2},
		{"BestWindowInTheMiddle", 6, []int{1, 1, 4, 3, 1, 1}, 2},
		{"WindowStartMovesSeveralSteps", 8, []int{1, 1, 1, 1, 5, 3}, 2},
		{"AllEqual", 6, []int{2, 2, 2, 2}, 3},
		{"AllEqualExact", 8, []int{2, 2, 2, 2}, 4},
		{"TargetOne", 1, []int{1, 2, 3}, 1},
		{"OneBigNumber", 100, []int{1, 1, 100, 1}, 1},
		{"BigNumberTooEarly", 100, []int{100, 1, 1, 1}, 1},
		{"GrowingNumbers", 12, []int{1, 2, 3, 4, 5}, 3},
		{"FallingNumbers", 12, []int{5, 4, 3, 2, 1}, 3},
		{"TwoWindowsOfTheSameLength", 5, []int{2, 3, 1, 2, 3}, 2},
		{"LargeTarget", 1000000000, []int{10000, 10000, 10000}, 0},
		{"LargeNumbers", 20000, []int{10000, 10000, 10000}, 2},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := minSubArrayLen(c.target, c.nums); got != c.want {
				t.Errorf("minSubArrayLen(%d, %v) = %d, want %d", c.target, c.nums, got, c.want)
			}
		})
	}
}

// TestMinSubArrayLenEveryShortArray walks every array of up to 7 numbers
// between 1 and 3, and every target from 1 to the sum of the array plus 1.
// It catches a window that starts or stops one step too late, and a window
// start that moves too far.
func TestMinSubArrayLenEveryShortArray(t *testing.T) {
	const (
		maxLen = 7
		maxVal = 3
	)

	for n := 0; n <= maxLen; n++ {
		failed := false

		everyArray(n, maxVal, func(nums []int) {
			if failed {
				return
			}
			for target := 1; target <= n*maxVal+1; target++ {
				if got, want := minSubArrayLen(target, nums), brute(target, nums); got != want {
					t.Errorf("minSubArrayLen(%d, %v) = %d, want %d", target, nums, got, want)
					failed = true
					return
				}
			}
		})
	}
}

// TestMinSubArrayLenRandom compares the solution with brute on random arrays.
// Small numbers give long windows. Large numbers give short windows.
func TestMinSubArrayLenRandom(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	maxVals := []int{1, 2, 10, 1000, 10000}

	for _, max := range maxVals {
		for range 200 {
			nums := randomNums(1+rnd.Intn(200), max, rnd)

			// The largest target is above the sum of the array, so the
			// answer is 0 for a part of the cases.
			sum := 0
			for _, v := range nums {
				sum += v
			}
			target := 1 + rnd.Intn(sum+sum/4+1)

			if got, want := minSubArrayLen(target, nums), brute(target, nums); got != want {
				t.Fatalf("minSubArrayLen(%d, %v) = %d, want %d", target, nums, got, want)
			}
		}
	}
}

// TestMinSubArrayLenEveryTarget uses one array and walks every target from 1
// to the sum of the array plus 5. The answer grows one step at a time, so a
// wrong border shows up at once.
func TestMinSubArrayLenEveryTarget(t *testing.T) {
	arrays := [][]int{
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{2, 3, 1, 2, 4, 3},
		{5, 1, 3, 5, 10, 7, 4, 9, 2, 8},
		{10000, 1, 10000, 1, 10000},
	}

	for _, nums := range arrays {
		sum := 0
		for _, v := range nums {
			sum += v
		}

		for target := 1; target <= sum+5; target++ {
			if got, want := minSubArrayLen(target, nums), brute(target, nums); got != want {
				t.Errorf("minSubArrayLen(%d, %v) = %d, want %d", target, nums, got, want)
			}
		}
	}
}

// TestMinSubArrayLenBestWindowEverywhere puts the one short window at every
// position of the array. The other numbers are 1, so a window of ones needs
// many more numbers to reach the target.
func TestMinSubArrayLenBestWindowEverywhere(t *testing.T) {
	const (
		target = 10
		n      = 12
	)

	for at := 0; at < n; at++ {
		nums := make([]int, n)
		for i := range nums {
			nums[i] = 1
		}
		// This number alone is 9 of the target, so the best window is this
		// number and one neighbour.
		nums[at] = 9

		if got, want := minSubArrayLen(target, nums), 2; got != want {
			t.Errorf("minSubArrayLen(%d, %v) = %d, want %d", target, nums, got, want)
		}
	}
}

// TestMinSubArrayLenKeepsInput checks that minSubArrayLen leaves the array
// alone. A solution that writes the running sums back into nums fails here.
func TestMinSubArrayLenKeepsInput(t *testing.T) {
	nums := []int{2, 3, 1, 2, 4, 3}
	want := slices.Clone(nums)

	minSubArrayLen(7, nums)

	if !slices.Equal(nums, want) {
		t.Errorf("minSubArrayLen changed the array to %v, want %v", nums, want)
	}
}

// TestMinSubArrayLenLarge uses the largest array of the constraints.
func TestMinSubArrayLenLarge(t *testing.T) {
	const n = 100000

	t.Run("AllOnes", func(t *testing.T) {
		nums := make([]int, n)
		for i := range nums {
			nums[i] = 1
		}

		if got, want := minSubArrayLen(n, nums), n; got != want {
			t.Errorf("minSubArrayLen(%d, %d ones) = %d, want %d", n, n, got, want)
		}
		if got, want := minSubArrayLen(n+1, nums), 0; got != want {
			t.Errorf("minSubArrayLen(%d, %d ones) = %d, want %d", n+1, n, got, want)
		}
	})

	t.Run("BestWindowAtTheEnd", func(t *testing.T) {
		// Every number is 1 but the last two, so the only short window sits
		// at the end and the solution must read the whole array.
		nums := make([]int, n)
		for i := range nums {
			nums[i] = 1
		}
		nums[n-2] = 10000
		nums[n-1] = 10000

		if got, want := minSubArrayLen(20000, nums), 2; got != want {
			t.Errorf("minSubArrayLen(20000, %d numbers) = %d, want %d", n, got, want)
		}
	})

	t.Run("LargestTarget", func(t *testing.T) {
		// The sum of the array is 10^9, and that is the largest target of
		// the constraints. The answer is the whole array.
		nums := make([]int, n)
		for i := range nums {
			nums[i] = 10000
		}

		if got, want := minSubArrayLen(1000000000, nums), n; got != want {
			t.Errorf("minSubArrayLen(1000000000, %d numbers) = %d, want %d", n, got, want)
		}
	})
}

func BenchmarkMinSubArrayLenSmall(b *testing.B) {
	nums := []int{2, 3, 1, 2, 4, 3}

	for b.Loop() {
		minSubArrayLen(7, nums)
	}
}

// BenchmarkMinSubArrayLenAllOnes gives one value many times. The window start
// moves as often as the window end.
func BenchmarkMinSubArrayLenAllOnes(b *testing.B) {
	const n = 100000

	nums := make([]int, n)
	for i := range nums {
		nums[i] = 1
	}

	b.ResetTimer()
	for b.Loop() {
		minSubArrayLen(500, nums)
	}
}

// BenchmarkMinSubArrayLenNoAnswer uses a target above the sum of the array,
// so the solution reads the whole array and finds nothing.
func BenchmarkMinSubArrayLenNoAnswer(b *testing.B) {
	const n = 100000

	nums := make([]int, n)
	for i := range nums {
		nums[i] = 1
	}

	b.ResetTimer()
	for b.Loop() {
		minSubArrayLen(n+1, nums)
	}
}

// BenchmarkMinSubArrayLenRandom gives random numbers, so the windows have
// different lengths.
func BenchmarkMinSubArrayLenRandom(b *testing.B) {
	const n = 100000

	nums := randomNums(n, 10000, rand.New(rand.NewSource(2)))

	b.ResetTimer()
	for b.Loop() {
		minSubArrayLen(1000000, nums)
	}
}
