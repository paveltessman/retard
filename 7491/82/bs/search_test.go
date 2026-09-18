package bs

import (
	"math"
	"math/rand"
	"slices"
	"testing"
	"time"
)

// contract tells the test what the index of a binary search means.
type contract int

const (
	// anyMatch asks only that the index points to a number equal to val on a
	// hit. Any of the equal numbers is fine. On a miss the test looks at the
	// bool alone and ignores the index.
	anyMatch contract = iota

	// firstMatch is the contract of slices.BinarySearch, of lower_bound of C++
	// and of bisect_left of Python. The index is the count of the numbers
	// smaller than val. On a hit it points to the first number equal to val.
	// On a miss it points to the place where val goes in.
	firstMatch

	// afterMatch is the contract of upper_bound of C++ and of bisect_right of
	// Python. The index is the count of the numbers that are val or smaller.
	// On a hit it points one place after the last number equal to val. On a
	// miss it points to the place where val goes in, as firstMatch does.
	afterMatch

	// lastMatch asks for the last number equal to val on a hit. On a miss the
	// test looks at the bool alone and ignores the index. Pick this contract
	// for a search that gives the last number that is val or smaller, which is
	// the afterMatch index minus one.
	lastMatch
)

// searchCase names one binary search of the package and gives the function.
//
// The contract of the function: nums is sorted from small to large.
// search(nums, val) gives back true if nums holds val, and false if not. The
// function must not change nums.
//
// The third field tells the test what the index means:
//
//	{"BinarySearch", binarySearch, anyMatch}  // any of the equal numbers
//	{"LowerBound", lowerBound, firstMatch}    // the first equal number
//	{"UpperBound", upperBound, afterMatch}    // one place after the last equal number
//	{"UpperBound", upperBound, lastMatch}     // the last equal number
//
// A generic function fits this shape as binarySearch[int].
type searchCase struct {
	name     string
	search   func(nums []int, val int) (int, bool)
	contract contract
}

// searches holds every binary search under test. Add one line for each new
// function, then run go test ./2/bs.
var searches = []searchCase{
	{"BinarySearch", binarySearch, anyMatch},
	{"LowerBound", lowerBound, firstMatch},
	{"UpperBound", upperBound, afterMatch},
}

// countAtMost gives the count of the numbers of nums that are val or smaller.
// That count is the answer that the afterMatch contract asks for.
func countAtMost(nums []int, val int) int {
	i, _ := slices.BinarySearchFunc(nums, val, func(num, val int) int {
		if num <= val {
			return -1
		}
		return 1
	})
	return i
}

// checkSearch looks for val in nums and compares the answer with the answer of
// the standard library. It also checks that the search leaves nums as it is.
// It gives back false on the first wrong answer, so a loop over many arrays
// can stop early.
func checkSearch(t *testing.T, s searchCase, nums []int, val int) bool {
	t.Helper()

	firstIdx, wantFound := slices.BinarySearch(nums, val)

	in := slices.Clone(nums)
	gotIdx, gotFound := s.search(in, val)

	if !slices.Equal(in, nums) {
		t.Errorf("%s(%v, %d) changed the array to %v", s.name, nums, val, in)
		return false
	}
	if gotFound != wantFound {
		t.Errorf("%s(%v, %d) = %d, %t, want the bool %t", s.name, nums, val, gotIdx, gotFound, wantFound)
		return false
	}

	// wantIdx holds the one index that the contract allows. It stays at -1
	// while the contract allows more than one index.
	wantIdx := -1

	switch s.contract {
	case firstMatch:
		wantIdx = firstIdx
	case afterMatch:
		wantIdx = countAtMost(nums, val)
	case lastMatch:
		if gotFound {
			wantIdx = countAtMost(nums, val) - 1
		}
	}

	if wantIdx >= 0 {
		if gotIdx != wantIdx {
			t.Errorf("%s(%v, %d) = %d, %t, want the index %d", s.name, nums, val, gotIdx, gotFound, wantIdx)
			return false
		}
		return true
	}

	// The contract allows any index that points to a number equal to val.
	if !gotFound {
		return true
	}
	if gotIdx < 0 || gotIdx >= len(nums) {
		t.Errorf("%s(%v, %d) = %d, %t, want an index from 0 to %d", s.name, nums, val, gotIdx, gotFound, len(nums)-1)
		return false
	}
	if nums[gotIdx] != val {
		t.Errorf("%s(%v, %d) = %d, %t, but the number at that index is %d", s.name, nums, val, gotIdx, gotFound, nums[gotIdx])
		return false
	}
	return true
}

// checkEveryValue looks for every number that can tell two searches apart: the
// numbers of nums, the neighbours of those numbers, and the two extreme
// values.
func checkEveryValue(t *testing.T, s searchCase, nums []int) bool {
	t.Helper()

	for _, val := range interestingValues(nums) {
		if !checkSearch(t, s, nums, val) {
			return false
		}
	}
	return true
}

// interestingValues gives the numbers of nums, the number below and the number
// above each of them, and the two extreme values. The neighbours fall in the
// gaps between the numbers of nums and outside both ends.
func interestingValues(nums []int) []int {
	vals := make([]int, 0, 3*len(nums)+3)

	for _, v := range nums {
		if v > math.MinInt {
			vals = append(vals, v-1)
		}
		vals = append(vals, v)
		if v < math.MaxInt {
			vals = append(vals, v+1)
		}
	}
	vals = append(vals, 0, math.MinInt, math.MaxInt)

	return vals
}

// sortedNums gives n sorted numbers between -max and max. A small max gives
// many equal numbers.
func sortedNums(n, max int, rnd *rand.Rand) []int {
	nums := make([]int, n)
	for i := range nums {
		nums[i] = rnd.Intn(2*max+1) - max
	}
	slices.Sort(nums)
	return nums
}

// everySortedArray calls f once for every sorted array of n numbers taken from
// vals. vals must be sorted.
func everySortedArray(n int, vals []int, f func([]int)) {
	nums := make([]int, n)

	var build func(i, from int)
	build = func(i, from int) {
		if i == n {
			f(nums)
			return
		}
		for j := from; j < len(vals); j++ {
			nums[i] = vals[j]
			build(i+1, j)
		}
	}

	build(0, 0)
}

// runs gives an array that holds each number of vals count times, one block
// after the other.
func runs(vals []int, count int) []int {
	nums := make([]int, 0, len(vals)*count)
	for _, v := range vals {
		for range count {
			nums = append(nums, v)
		}
	}
	return nums
}

func TestSearchesAreRegistered(t *testing.T) {
	if len(searches) == 0 {
		t.Fatal("no binary search is registered in searches")
	}
}

// TestSearch looks for every interesting number in arrays of many shapes. The
// empty array, the array of one number and the array of equal numbers catch
// the bugs of the first loop step. The gaps catch a miss that reports a hit.
func TestSearch(t *testing.T) {
	cases := []struct {
		name string
		nums []int
	}{
		{"Nil", nil},
		{"Empty", []int{}},
		{"One", []int{5}},
		{"TwoDistinct", []int{1, 2}},
		{"TwoEqual", []int{7, 7}},
		{"TwoFarApart", []int{-100, 100}},
		{"ThreeDistinct", []int{1, 2, 3}},
		{"ThreeEqual", []int{4, 4, 4}},
		{"ThreeWithTheMiddleDoubled", []int{1, 2, 2}},
		{"FourDistinct", []int{1, 2, 3, 4}},
		{"FiveDistinct", []int{1, 2, 3, 4, 5}},
		{"EvenLength", []int{1, 3, 5, 7, 9, 11}},
		{"OddLength", []int{1, 3, 5, 7, 9}},
		{"AllEqualLong", []int{2, 2, 2, 2, 2, 2, 2, 2}},
		{"TwoValuesOnly", []int{0, 0, 0, 1, 1, 1}},
		{"FirstValueAlone", []int{0, 1, 1, 1, 1, 1}},
		{"LastValueAlone", []int{1, 1, 1, 1, 1, 2}},
		{"Duplicates", []int{1, 1, 3, 3, 3, 5, 5}},
		{"LongRunOfEqualNumbers", runs([]int{1, 2, 3}, 20)},
		{"Negatives", []int{-5, -4, -3, -2, -1}},
		{"MixedSigns", []int{-7, -2, -1, 0, 1, 2, 3}},
		{"BigGaps", []int{-1000000, -10, 0, 10, 1000000}},
		{"ExtremeValues", []int{math.MinInt, -1, 0, 1, math.MaxInt}},
		{"ExtremeValuesDoubled", []int{math.MinInt, math.MinInt, math.MaxInt, math.MaxInt}},
		{"Long", []int{2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30}},
	}

	for _, s := range searches {
		t.Run(s.name, func(t *testing.T) {
			for _, c := range cases {
				t.Run(c.name, func(t *testing.T) {
					checkEveryValue(t, s, c.nums)
				})
			}
		})
	}
}

// TestSearchEveryShortArray walks every sorted array of up to 7 numbers taken
// from 0, 2 and 4, and looks for every number from -1 to 5. The gaps at 1 and
// 3 catch a search that reports a hit for a number that is not there. The
// equal numbers catch a loop that steps over a block of equal numbers.
func TestSearchEveryShortArray(t *testing.T) {
	const maxLen = 7

	vals := []int{0, 2, 4}

	for _, s := range searches {
		t.Run(s.name, func(t *testing.T) {
			for n := 0; n <= maxLen; n++ {
				failed := false

				everySortedArray(n, vals, func(nums []int) {
					if failed {
						return
					}
					for val := -1; val <= 5; val++ {
						if !checkSearch(t, s, nums, val) {
							failed = true
							return
						}
					}
				})

				if failed {
					return
				}
			}
		})
	}
}

// TestSearchEveryIndexOfEveryLength looks for every number of a sorted array of
// distinct numbers, for every length up to 64. A wrong midpoint shows up at one
// length only, so the test needs many lengths.
func TestSearchEveryIndexOfEveryLength(t *testing.T) {
	const maxLen = 64

	for _, s := range searches {
		t.Run(s.name, func(t *testing.T) {
			for n := 0; n <= maxLen; n++ {
				nums := make([]int, n)
				for i := range nums {
					nums[i] = 2 * i
				}

				if !checkEveryValue(t, s, nums) {
					return
				}
			}
		})
	}
}

// TestSearchRandom looks for random numbers in random sorted arrays. A small
// max gives many equal numbers. A large max gives mostly misses.
func TestSearchRandom(t *testing.T) {
	maxes := []int{1, 3, 20, 1000000}

	for _, s := range searches {
		t.Run(s.name, func(t *testing.T) {
			rnd := rand.New(rand.NewSource(1))

			for _, max := range maxes {
				for range 300 {
					nums := sortedNums(rnd.Intn(200), max, rnd)

					// A number of the array gives a hit. A random number
					// gives a hit or a miss.
					if len(nums) > 0 {
						if !checkSearch(t, s, nums, nums[rnd.Intn(len(nums))]) {
							return
						}
					}
					if !checkSearch(t, s, nums, rnd.Intn(4*max+1)-2*max) {
						return
					}
				}
			}
		})
	}
}

// TestSearchKeepsNeighbours gives a part of a larger array to the search. A
// search that reads or writes outside of nums fails here.
func TestSearchKeepsNeighbours(t *testing.T) {
	for _, s := range searches {
		t.Run(s.name, func(t *testing.T) {
			arr := []int{111, 222, 1, 3, 5, 7, 9, 333, 444}
			nums := arr[2:7]

			for _, val := range []int{0, 1, 2, 5, 8, 9, 10, 111, 333} {
				checkSearch(t, s, nums, val)
			}

			if want := []int{111, 222, 1, 3, 5, 7, 9, 333, 444}; !slices.Equal(arr, want) {
				t.Errorf("%s changed the larger array to %v, want %v", s.name, arr, want)
			}
		})
	}
}

// TestSearchLarge works on large arrays. The array of equal numbers is the
// hard shape for a search that gives the first of the equal numbers.
func TestSearchLarge(t *testing.T) {
	const n = 100000

	rnd := rand.New(rand.NewSource(2))

	even := make([]int, n)
	for i := range even {
		even[i] = 2 * i
	}

	shapes := []struct {
		name string
		nums []int
	}{
		{"EvenNumbers", even},
		{"AllEqual", make([]int, n)},
		{"FewValues", sortedNums(n, 3, rnd)},
		{"Random", sortedNums(n, 1000000, rnd)},
	}

	for _, s := range searches {
		t.Run(s.name, func(t *testing.T) {
			for _, sh := range shapes {
				t.Run(sh.name, func(t *testing.T) {
					// The first, the last and the middle number.
					for _, i := range []int{0, len(sh.nums) / 2, len(sh.nums) - 1} {
						if !checkSearch(t, s, sh.nums, sh.nums[i]) {
							return
						}
					}
					for range 200 {
						val := rnd.Intn(4000000) - 2000000
						if !checkSearch(t, s, sh.nums, val) {
							return
						}
					}
				})
			}
		})
	}
}

// TestSearchIsLogarithmic makes many searches in a very large array. A binary
// search needs some milliseconds. A linear scan needs hours, so it runs into
// the limit.
func TestSearchIsLogarithmic(t *testing.T) {
	if testing.Short() {
		t.Skip("the test needs a large array")
	}

	const n = 1 << 22
	const queries = 100000
	const limit = 10 * time.Second

	nums := make([]int, n)
	for i := range nums {
		nums[i] = 2 * i
	}

	rnd := rand.New(rand.NewSource(3))

	for _, s := range searches {
		t.Run(s.name, func(t *testing.T) {
			start := time.Now()

			for range queries {
				val := rnd.Intn(2 * n)

				idx, found := s.search(nums, val)

				if want := val%2 == 0; found != want {
					t.Fatalf("%s(nums, %d) = %d, %t, want the bool %t", s.name, val, idx, found, want)
				}
				if time.Since(start) > limit {
					t.Fatalf("%s needed more than %v for the searches in an array of %d numbers, so it is not logarithmic", s.name, limit, n)
				}
			}
		})
	}
}

func BenchmarkSearch(b *testing.B) {
	const n = 100000

	rnd := rand.New(rand.NewSource(4))

	even := make([]int, n)
	for i := range even {
		even[i] = 2 * i
	}

	shapes := []struct {
		name string
		nums []int
		val  int
	}{
		{"Hit", even, 2 * (n / 3)},
		{"Miss", even, 2*(n/3) + 1},
		{"FirstNumber", even, even[0]},
		{"LastNumber", even, even[n-1]},
		{"BelowAll", even, -1},
		{"AboveAll", even, 2 * n},
		{"AllEqual", make([]int, n), 0},
		{"FewValues", sortedNums(n, 3, rnd), 0},
	}

	for _, s := range searches {
		b.Run(s.name, func(b *testing.B) {
			for _, sh := range shapes {
				b.Run(sh.name, func(b *testing.B) {
					for b.Loop() {
						s.search(sh.nums, sh.val)
					}
				})
			}
		})
	}
}
