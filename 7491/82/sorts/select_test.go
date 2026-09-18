package sorts

import (
	"math"
	"math/rand"
	"slices"
	"testing"
)

// selectCase names one function that finds the k-th order statistic
// (quickselect).
//
// The contract: kth(nums, k) gives the number that stands at index k of the
// sorted nums, with k from 0 to len(nums)-1. The function can reorder nums,
// but it must keep the same numbers. For the k-th largest number, the caller
// passes len(nums)-k.
//
// Add one line for each new function:
//
//	{"QuickSelect", quickSelect[int]}
var selectors []selectCase = []selectCase{
	{"QuickSelect", quickSelect},
}

type selectCase struct {
	name string
	kth  func(nums []int, k int) int
}

// checkSelect asks for the number at index k of the sorted nums and compares
// the answer with the answer of the standard library. It also checks that the
// function keeps the same numbers in nums. It gives back false on the first
// wrong answer.
func checkSelect(t *testing.T, name string, kth func([]int, int) int, nums []int, k int) bool {
	t.Helper()

	want := slices.Clone(nums)
	slices.Sort(want)

	in := slices.Clone(nums)

	got := kth(in, k)

	if got != want[k] {
		t.Errorf("%s(%v, %d) = %d, want %d", name, nums, k, got, want[k])
		return false
	}

	left := slices.Clone(in)
	slices.Sort(left)
	if !slices.Equal(left, want) {
		t.Errorf("%s(%v, %d) left the numbers %v, want the same numbers in any order", name, nums, k, in)
		return false
	}
	return true
}

// checkEveryIndex asks for every index of nums, from 0 to len(nums)-1.
func checkEveryIndex(t *testing.T, name string, kth func([]int, int) int, nums []int) bool {
	t.Helper()

	for k := range nums {
		if !checkSelect(t, name, kth, nums, k) {
			return false
		}
	}
	return true
}

func TestSelectorsAreRegistered(t *testing.T) {
	if len(selectors) == 0 {
		t.Skip("no k-th order statistic function is registered in selectors")
	}
}

func TestSelect(t *testing.T) {
	cases := []struct {
		name string
		nums []int
	}{
		{"One", []int{5}},
		{"TwoSorted", []int{1, 2}},
		{"TwoReversed", []int{2, 1}},
		{"TwoEqual", []int{7, 7}},
		{"ThreeReversed", []int{3, 2, 1}},
		{"ThreeEqual", []int{4, 4, 4}},
		{"EvenLength", []int{4, 1, 3, 2}},
		{"OddLength", []int{5, 1, 4, 2, 3}},
		{"Sorted", []int{1, 2, 3, 4, 5, 6, 7, 8}},
		{"Reversed", []int{8, 7, 6, 5, 4, 3, 2, 1}},
		{"AllEqual", []int{2, 2, 2, 2, 2, 2}},
		{"TwoValuesOnly", []int{1, 0, 1, 0, 1, 0}},
		{"Duplicates", []int{5, 1, 5, 1, 3, 3, 5}},
		{"Negatives", []int{-1, -5, -3, -2, -4}},
		{"MixedSigns", []int{3, -1, 0, -7, 2, -2, 1}},
		{"ExtremeValues", []int{math.MaxInt, math.MinInt, 0, math.MaxInt, math.MinInt}},
		{"Sawtooth", sawtooth(30, 5)},
		{"OrganPipe", organPipe(31)},
	}

	for _, s := range selectors {
		t.Run(s.name, func(t *testing.T) {
			for _, c := range cases {
				t.Run(c.name, func(t *testing.T) {
					checkEveryIndex(t, s.name, s.kth, c.nums)
				})
			}
		})
	}
}

// TestSelectEveryShortArray walks every array of up to 6 numbers taken from 0,
// 1 and 2, and asks for every index. The many equal numbers catch a partition
// that moves the wrong side.
func TestSelectEveryShortArray(t *testing.T) {
	const maxLen = 6

	vals := []int{0, 1, 2}

	for _, s := range selectors {
		t.Run(s.name, func(t *testing.T) {
			for n := 1; n <= maxLen; n++ {
				failed := false

				everyArray(n, vals, func(nums []int) {
					if failed {
						return
					}
					if !checkEveryIndex(t, s.name, s.kth, nums) {
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

// TestSelectEveryPermutation walks every order of the numbers 1 to 6 and asks
// for every index. All numbers differ, so a wrong index gives a wrong number.
func TestSelectEveryPermutation(t *testing.T) {
	const maxLen = 6

	for _, s := range selectors {
		t.Run(s.name, func(t *testing.T) {
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
					if !checkEveryIndex(t, s.name, s.kth, nums) {
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

// TestSelectRandom asks for a random index of random arrays. A small max gives
// many equal numbers.
func TestSelectRandom(t *testing.T) {
	maxes := []int{1, 3, 20, 1000000}

	for _, s := range selectors {
		t.Run(s.name, func(t *testing.T) {
			rnd := rand.New(rand.NewSource(1))

			for _, max := range maxes {
				for range 300 {
					nums := randomNums(1+rnd.Intn(200), max, rnd)

					if !checkSelect(t, s.name, s.kth, nums, rnd.Intn(len(nums))) {
						return
					}
				}
			}
		})
	}
}

// TestSelectFirstAndLast asks for the smallest and the largest number. These
// two indexes sit at the edge, so an off by one bug shows up here first.
func TestSelectFirstAndLast(t *testing.T) {
	rnd := rand.New(rand.NewSource(2))

	for _, s := range selectors {
		t.Run(s.name, func(t *testing.T) {
			for range 200 {
				nums := randomNums(1+rnd.Intn(100), 1000, rnd)

				if !checkSelect(t, s.name, s.kth, nums, 0) {
					return
				}
				if !checkSelect(t, s.name, s.kth, nums, len(nums)-1) {
					return
				}
			}
		})
	}
}

// TestSelectKeepsNeighbours gives a part of a larger array to the function. A
// function that writes outside of nums fails here.
func TestSelectKeepsNeighbours(t *testing.T) {
	for _, s := range selectors {
		t.Run(s.name, func(t *testing.T) {
			arr := []int{111, 222, 5, 1, 4, 2, 3, 333, 444}
			nums := arr[2:7]

			got := s.kth(nums, 2)

			if want := 3; got != want {
				t.Errorf("%s(%v, 2) = %d, want %d", s.name, []int{5, 1, 4, 2, 3}, got, want)
			}
			if got, want := arr[:2], []int{111, 222}; !slices.Equal(got, want) {
				t.Errorf("%s changed the numbers before nums to %v, want %v", s.name, got, want)
			}
			if got, want := arr[7:], []int{333, 444}; !slices.Equal(got, want) {
				t.Errorf("%s changed the numbers after nums to %v, want %v", s.name, got, want)
			}
		})
	}
}

// TestSelectLarge works on large arrays. The shapes after Random make a
// quickselect with a weak pivot quadratic, so a very slow subtest points to
// the pivot choice.
func TestSelectLarge(t *testing.T) {
	const n = 100000
	const hard = 20000

	rnd := rand.New(rand.NewSource(3))

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
		{"OrganPipe", organPipe(hard)},
	}

	for _, s := range selectors {
		t.Run(s.name, func(t *testing.T) {
			for _, sh := range shapes {
				t.Run(sh.name, func(t *testing.T) {
					want := slices.Clone(sh.nums)
					slices.Sort(want)

					// The middle index needs the most work.
					k := len(sh.nums) / 2

					if got := s.kth(slices.Clone(sh.nums), k); got != want[k] {
						t.Errorf("%s gave %d for index %d of the %s array of %d numbers, want %d", s.name, got, k, sh.name, len(sh.nums), want[k])
					}
				})
			}
		})
	}
}

func BenchmarkSelect(b *testing.B) {
	const n = 10000

	rnd := rand.New(rand.NewSource(4))

	sorted := make([]int, n)
	for i := range sorted {
		sorted[i] = i
	}

	shapes := []struct {
		name string
		nums []int
	}{
		{"Random", randomNums(n, 1000000, rnd)},
		{"Sorted", sorted},
		{"AllEqual", make([]int, n)},
		{"FewValues", randomNums(n, 3, rnd)},
	}

	for _, s := range selectors {
		b.Run(s.name, func(b *testing.B) {
			for _, sh := range shapes {
				b.Run(sh.name, func(b *testing.B) {
					work := make([]int, len(sh.nums))
					k := len(sh.nums) / 2

					for b.Loop() {
						copy(work, sh.nums)
						s.kth(work, k)
					}
				})
			}
		})
	}
}
