package containsduplicate

import (
	"math/rand"
	"slices"
	"testing"
)

// distinct gives n values that are all different.
// The answer of the task for this array is false.
func distinct(n int) []int {
	nums := make([]int, n)
	for i := range nums {
		nums[i] = i * 3
	}
	return nums
}

// shuffle puts the values of nums in a new order, in place.
// The values stay the same, so the answer of the task stays the same.
func shuffle(nums []int, rnd *rand.Rand) {
	rnd.Shuffle(len(nums), func(i, j int) {
		nums[i], nums[j] = nums[j], nums[i]
	})
}

func TestContainsDuplicate(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		want bool
	}{
		{"Example1", []int{1, 2, 3, 1}, true},
		{"Example2", []int{1, 2, 3, 4}, false},
		{"Example3", []int{1, 1, 1, 3, 3, 4, 3, 2, 4, 2}, true},
		{"OneNumber", []int{7}, false},
		{"TwoSame", []int{5, 5}, true},
		{"TwoDifferent", []int{5, 6}, false},
		{"FirstTwo", []int{2, 2, 3, 4, 5}, true},
		{"LastTwo", []int{1, 2, 3, 4, 4}, true},
		{"FirstAndLast", []int{9, 1, 2, 3, 9}, true},
		{"MiddleTwo", []int{1, 2, 6, 6, 3, 4}, true},
		{"AllSame", []int{8, 8, 8, 8}, true},
		{"Zeros", []int{0, 0}, true},
		{"ZeroOnce", []int{0, 1, 2}, false},
		{"Negatives", []int{-1, -2, -3, -2}, true},
		{"NegativesDistinct", []int{-1, -2, -3, -4}, false},
		{"SignPair", []int{-5, 5}, false},
		{"SignPairWithDuplicate", []int{-5, 5, -5}, true},
		{"Limits", []int{1000000000, -1000000000}, false},
		{"LimitsDuplicate", []int{1000000000, -1000000000, 1000000000}, true},
		{"FarApart", []int{4, 1, 2, 3, 5, 6, 7, 8, 9, 4}, true},
		{"Neighbours", []int{1, 3, 5, 7, 7, 9}, true},
		{"SortedDistinct", []int{-3, -1, 0, 2, 4}, false},
		{"Empty", []int{}, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := containsDuplicate(c.nums); got != c.want {
				t.Errorf("containsDuplicate(%v) = %v, want %v", c.nums, got, c.want)
			}
		})
	}
}

// TestContainsDuplicateEveryPair walks every array length up to 12 and puts
// the one repeated value in every pair of positions. The test catches a loop
// that starts or stops one step too late.
func TestContainsDuplicateEveryPair(t *testing.T) {
	const maxLen = 12

	for n := 2; n <= maxLen; n++ {
		for a := 0; a < n; a++ {
			for b := a + 1; b < n; b++ {
				nums := distinct(n)
				nums[b] = nums[a]

				if !containsDuplicate(nums) {
					t.Fatalf("containsDuplicate(%v) = false, want true (same value at %d and %d)", nums, a, b)
				}
			}
		}
	}
}

// TestContainsDuplicateDistinct gives arrays of different values only.
// The array holds negative numbers, zero, and positive numbers.
func TestContainsDuplicateDistinct(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))

	for n := 1; n <= 60; n++ {
		nums := make([]int, n)
		for i := range nums {
			nums[i] = i - n/2
		}
		shuffle(nums, rnd)

		if containsDuplicate(nums) {
			t.Fatalf("containsDuplicate(%v) = true, want false", nums)
		}
	}
}

// TestContainsDuplicateOrderMeansNothing shuffles the same array many times.
// The values stay the same, so every shuffle gives the same answer.
// A solution that compares neighbours only fails here.
func TestContainsDuplicateOrderMeansNothing(t *testing.T) {
	rnd := rand.New(rand.NewSource(2))

	cases := []struct {
		nums []int
		want bool
	}{
		{[]int{1, 2, 3, 1}, true},
		{[]int{1, 2, 3, 4}, false},
		{[]int{1, 1, 1, 3, 3, 4, 3, 2, 4, 2}, true},
		{[]int{-2, 7, 0, 5, 7}, true},
		{[]int{-2, 7, 0, 5, 9}, false},
	}

	for _, c := range cases {
		nums := slices.Clone(c.nums)
		for range 50 {
			shuffle(nums, rnd)

			if got := containsDuplicate(nums); got != c.want {
				t.Fatalf("containsDuplicate(%v) = %v, want %v", nums, got, c.want)
			}
		}
	}
}

// TestContainsDuplicateKeepsInput checks that the solution leaves the array alone.
// A solution that sorts nums must copy the array first.
func TestContainsDuplicateKeepsInput(t *testing.T) {
	nums := []int{8, 3, 11, 2, 7, 5}
	want := slices.Clone(nums)

	containsDuplicate(nums)

	if !slices.Equal(nums, want) {
		t.Errorf("containsDuplicate changed the array to %v, want %v", nums, want)
	}
}

// TestContainsDuplicateLarge uses the largest array of the constraints.
// The two equal values sit at the two ends, so the solution must read
// the whole array.
func TestContainsDuplicateLarge(t *testing.T) {
	const n = 100000

	nums := distinct(n)
	if containsDuplicate(nums) {
		t.Fatalf("containsDuplicate of %d different values = true, want false", n)
	}

	nums[n-1] = nums[0]
	if !containsDuplicate(nums) {
		t.Errorf("containsDuplicate of %d values with the first value repeated at the end = false, want true", n)
	}
}

func BenchmarkContainsDuplicateSmall(b *testing.B) {
	nums := []int{1, 2, 3, 1}

	for b.Loop() {
		containsDuplicate(nums)
	}
}

// BenchmarkContainsDuplicateDistinct gives the worst case: the solution
// must read every value before it answers false.
func BenchmarkContainsDuplicateDistinct(b *testing.B) {
	const n = 100000

	nums := distinct(n)
	shuffle(nums, rand.New(rand.NewSource(3)))

	b.ResetTimer()
	for b.Loop() {
		containsDuplicate(nums)
	}
}

// BenchmarkContainsDuplicateEarlyExit repeats the second value at once.
// A solution that answers as soon as it sees a repeat stops after two steps.
func BenchmarkContainsDuplicateEarlyExit(b *testing.B) {
	const n = 100000

	nums := distinct(n)
	nums[1] = nums[0]

	b.ResetTimer()
	for b.Loop() {
		containsDuplicate(nums)
	}
}
