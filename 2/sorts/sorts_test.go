package sorts

import (
	"math"
	"math/rand"
	"slices"
	"testing"
)

// sortCase names one sorting algorithm of the package and gives a function
// that sorts a slice of ints.
//
// A sort that gives back a new slice fits this shape as it is:
//
//	{"MergeSort", mergeSort[int]}
//
// A sort that works on the slice of the caller and gives back nothing needs
// the inPlace wrapper:
//
//	{"InsertionSort", inPlace(insertionSort[int])}
//
// A sort that takes the bounds of the part to sort needs a small closure:
//
//	{"QuickSort", func(nums []int) []int { quickSort(nums, 0, len(nums)-1); return nums }}
type sortCase struct {
	name string
	sort func(nums []int) []int
}

// algorithms holds every sorting algorithm under test. Add one line for each
// new algorithm, then run go test ./2/sorts.
var algorithms = []sortCase{
	{"MergeSort", mergeSort[int]},
	{"QuickSort", quickSort[int]},
	{"HeapSort", heapSort},
}

// inPlace turns a sort that works on the slice of the caller into the shape
// that sortCase needs.
func inPlace[T any](sort func(nums []T)) func(nums []T) []T {
	return func(nums []T) []T {
		sort(nums)
		return nums
	}
}

// checkSort sorts a copy of nums and compares the answer with the answer of
// the standard library. It gives back false on the first wrong answer, so a
// loop over many arrays can stop early.
func checkSort(t *testing.T, name string, sort func([]int) []int, nums []int) bool {
	t.Helper()

	want := slices.Clone(nums)
	slices.Sort(want)

	got := sort(slices.Clone(nums))

	if !slices.Equal(got, want) {
		t.Errorf("%s(%v) = %v, want %v", name, nums, got, want)
		return false
	}
	return true
}

// randomNums gives n numbers between -max and max. A small max gives many
// equal numbers.
func randomNums(n, max int, rnd *rand.Rand) []int {
	nums := make([]int, n)
	for i := range nums {
		nums[i] = rnd.Intn(2*max+1) - max
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

// everyPermutation calls f once for every order of the numbers in nums.
func everyPermutation(nums []int, f func([]int)) {
	var build func(i int)
	build = func(i int) {
		if i == len(nums) {
			f(nums)
			return
		}
		for j := i; j < len(nums); j++ {
			nums[i], nums[j] = nums[j], nums[i]
			build(i + 1)
			nums[i], nums[j] = nums[j], nums[i]
		}
	}

	build(0)
}

// sawtooth gives n numbers that count from 0 to period-1 again and again.
func sawtooth(n, period int) []int {
	nums := make([]int, n)
	for i := range nums {
		nums[i] = i % period
	}
	return nums
}

// organPipe gives n numbers that go up to the middle and down again. This
// shape makes a quicksort with a median of three pivot slow.
func organPipe(n int) []int {
	nums := make([]int, n)
	for i := range nums {
		if i < n/2 {
			nums[i] = i
			continue
		}
		nums[i] = n - i
	}
	return nums
}

func TestAlgorithmsAreRegistered(t *testing.T) {
	if len(algorithms) == 0 {
		t.Fatal("no sorting algorithm is registered in algorithms")
	}
}

func TestSort(t *testing.T) {
	cases := []struct {
		name string
		nums []int
	}{
		{"Nil", nil},
		{"Empty", []int{}},
		{"One", []int{5}},
		{"TwoSorted", []int{1, 2}},
		{"TwoReversed", []int{2, 1}},
		{"TwoEqual", []int{7, 7}},
		{"ThreeSorted", []int{1, 2, 3}},
		{"ThreeReversed", []int{3, 2, 1}},
		{"ThreeEqual", []int{4, 4, 4}},
		{"ThreeFirstOutOfPlace", []int{3, 1, 2}},
		{"ThreeLastOutOfPlace", []int{2, 3, 1}},
		{"ThreeMiddleOutOfPlace", []int{1, 3, 2}},
		{"EvenLength", []int{4, 1, 3, 2}},
		{"OddLength", []int{5, 1, 4, 2, 3}},
		{"SortedLong", []int{1, 2, 3, 4, 5, 6, 7, 8}},
		{"ReversedLong", []int{8, 7, 6, 5, 4, 3, 2, 1}},
		{"AllEqualLong", []int{2, 2, 2, 2, 2, 2, 2, 2}},
		{"OnlyTheFirstOutOfPlace", []int{9, 1, 2, 3, 4, 5}},
		{"OnlyTheLastOutOfPlace", []int{1, 2, 3, 4, 5, 0}},
		{"FirstAndLastSwapped", []int{6, 2, 3, 4, 5, 1}},
		{"TwoValuesOnly", []int{1, 0, 1, 0, 1, 0}},
		{"ManyZeroes", []int{0, 3, 0, 1, 0, 2, 0}},
		{"Duplicates", []int{5, 1, 5, 1, 3, 3, 5}},
		{"Negatives", []int{-1, -5, -3, -2, -4}},
		{"MixedSigns", []int{3, -1, 0, -7, 2, -2, 1}},
		{"ExtremeValues", []int{math.MaxInt, math.MinInt, 0, math.MaxInt, math.MinInt}},
		{"Sawtooth", sawtooth(30, 5)},
		{"OrganPipe", organPipe(31)},
		{"OneValueRepeatedWithOneOther", []int{7, 7, 7, 1, 7, 7, 7}},
	}

	for _, a := range algorithms {
		t.Run(a.name, func(t *testing.T) {
			for _, c := range cases {
				t.Run(c.name, func(t *testing.T) {
					checkSort(t, a.name, a.sort, c.nums)
				})
			}
		})
	}
}

// TestSortEveryShortArray walks every array of up to 7 numbers taken from 0, 1
// and 2. The many equal numbers catch a partition that loses a number or that
// runs one step too far.
func TestSortEveryShortArray(t *testing.T) {
	const maxLen = 7

	vals := []int{0, 1, 2}

	for _, a := range algorithms {
		t.Run(a.name, func(t *testing.T) {
			for n := 0; n <= maxLen; n++ {
				failed := false

				everyArray(n, vals, func(nums []int) {
					if failed {
						return
					}
					if !checkSort(t, a.name, a.sort, nums) {
						failed = true
					}
				})

				if failed {
					return
				}
			}
		})
	}
}

// TestSortEveryPermutation walks every order of the numbers 1 to 7. All
// numbers differ, so a wrong index shows up as a lost or a doubled number.
func TestSortEveryPermutation(t *testing.T) {
	const maxLen = 7

	for _, a := range algorithms {
		t.Run(a.name, func(t *testing.T) {
			for n := 1; n <= maxLen; n++ {
				nums := make([]int, n)
				for i := range nums {
					nums[i] = i + 1
				}

				failed := false

				everyPermutation(nums, func(nums []int) {
					if failed {
						return
					}
					if !checkSort(t, a.name, a.sort, nums) {
						failed = true
					}
				})

				if failed {
					return
				}
			}
		})
	}
}

// TestSortRandom sorts random arrays of many sizes. A small max gives many
// equal numbers. A large max gives few equal numbers.
func TestSortRandom(t *testing.T) {
	maxes := []int{1, 3, 20, 1000000}

	for _, a := range algorithms {
		t.Run(a.name, func(t *testing.T) {
			rnd := rand.New(rand.NewSource(1))

			for _, max := range maxes {
				for range 200 {
					nums := randomNums(rnd.Intn(300), max, rnd)

					if !checkSort(t, a.name, a.sort, nums) {
						return
					}
				}
			}
		})
	}
}

// TestSortKeepsNeighbours gives a part of a larger array to the sort. A sort
// that writes outside of nums fails here.
func TestSortKeepsNeighbours(t *testing.T) {
	for _, a := range algorithms {
		t.Run(a.name, func(t *testing.T) {
			arr := []int{111, 222, 5, 1, 4, 2, 3, 333, 444}
			nums := arr[2:7]

			got := a.sort(nums)

			if want := []int{1, 2, 3, 4, 5}; !slices.Equal(got, want) {
				t.Errorf("%s(%v) = %v, want %v", a.name, []int{5, 1, 4, 2, 3}, got, want)
			}
			if got, want := arr[:2], []int{111, 222}; !slices.Equal(got, want) {
				t.Errorf("%s changed the numbers before nums to %v, want %v", a.name, got, want)
			}
			if got, want := arr[7:], []int{333, 444}; !slices.Equal(got, want) {
				t.Errorf("%s changed the numbers after nums to %v, want %v", a.name, got, want)
			}
		})
	}
}

// TestSortLarge sorts large arrays. The shapes after Random are the shapes
// that make a quicksort with a weak pivot quadratic, so a very slow subtest
// points to the pivot choice.
func TestSortLarge(t *testing.T) {
	const n = 100000
	const hard = 20000

	rnd := rand.New(rand.NewSource(2))

	sorted := make([]int, hard)
	for i := range sorted {
		sorted[i] = i
	}
	reversed := make([]int, hard)
	for i := range reversed {
		reversed[i] = hard - i
	}

	shapes := []struct {
		name string
		nums []int
	}{
		{"Random", randomNums(n, 1000000, rnd)},
		{"Sorted", sorted},
		{"Reversed", reversed},
		{"AllEqual", make([]int, hard)},
		{"FewValues", randomNums(hard, 3, rnd)},
		{"Sawtooth", sawtooth(hard, 100)},
		{"OrganPipe", organPipe(hard)},
	}

	for _, a := range algorithms {
		t.Run(a.name, func(t *testing.T) {
			for _, s := range shapes {
				t.Run(s.name, func(t *testing.T) {
					want := slices.Clone(s.nums)
					slices.Sort(want)

					got := a.sort(slices.Clone(s.nums))

					if !slices.Equal(got, want) {
						t.Errorf("%s gave a wrong answer for the %s array of %d numbers", a.name, s.name, len(s.nums))
					}
				})
			}
		})
	}
}

func BenchmarkSort(b *testing.B) {
	const n = 10000

	rnd := rand.New(rand.NewSource(3))

	sorted := make([]int, n)
	for i := range sorted {
		sorted[i] = i
	}
	reversed := make([]int, n)
	for i := range reversed {
		reversed[i] = n - i
	}

	shapes := []struct {
		name string
		nums []int
	}{
		{"Random", randomNums(n, 1000000, rnd)},
		{"Sorted", sorted},
		{"Reversed", reversed},
		{"AllEqual", make([]int, n)},
		{"FewValues", randomNums(n, 3, rnd)},
		{"NearlySorted", nearlySorted(n, 100, rnd)},
		{"Small", randomNums(16, 1000, rnd)},
	}

	for _, a := range algorithms {
		b.Run(a.name, func(b *testing.B) {
			for _, s := range shapes {
				b.Run(s.name, func(b *testing.B) {
					work := make([]int, len(s.nums))

					for b.Loop() {
						copy(work, s.nums)
						a.sort(work)
					}
				})
			}
		})
	}
}

// nearlySorted gives a sorted array of n numbers with n/swaps pairs swapped.
// An insertion sort is fast on this shape.
func nearlySorted(n, swaps int, rnd *rand.Rand) []int {
	nums := make([]int, n)
	for i := range nums {
		nums[i] = i
	}
	for range n / swaps {
		i, j := rnd.Intn(n), rnd.Intn(n)
		nums[i], nums[j] = nums[j], nums[i]
	}
	return nums
}
