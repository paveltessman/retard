package sorts

import (
	"math"
	"math/rand"
	"slices"
	"testing"
)

// mergeCase names one function that merges two sorted slices into one sorted
// slice. The function must not change a or b.
type mergeCase struct {
	name  string
	merge func(a, b []int) []int
}

// merges holds every merge function under test. Add one line for each new
// function.
var merges = []mergeCase{
	{"Merge", merge[int]},
}

// checkMerge merges copies of a and b, then compares the answer with the
// sorted numbers of both slices. It also checks that the function leaves a and
// b as they are. It gives back false on the first wrong answer.
func checkMerge(t *testing.T, name string, m func(a, b []int) []int, a, b []int) bool {
	t.Helper()

	want := make([]int, 0, len(a)+len(b))
	want = append(want, a...)
	want = append(want, b...)
	slices.Sort(want)

	inA, inB := slices.Clone(a), slices.Clone(b)

	got := m(inA, inB)

	if !slices.Equal(got, want) {
		t.Errorf("%s(%v, %v) = %v, want %v", name, a, b, got, want)
		return false
	}
	if !slices.Equal(inA, a) {
		t.Errorf("%s(%v, %v) changed the first slice to %v", name, a, b, inA)
		return false
	}
	if !slices.Equal(inB, b) {
		t.Errorf("%s(%v, %v) changed the second slice to %v", name, a, b, inB)
		return false
	}
	return true
}

// sortedRandom gives n sorted numbers between -max and max.
func sortedRandom(n, max int, rnd *rand.Rand) []int {
	nums := randomNums(n, max, rnd)
	slices.Sort(nums)
	return nums
}

func TestMergesAreRegistered(t *testing.T) {
	if len(merges) == 0 {
		t.Fatal("no merge function is registered in merges")
	}
}

func TestMerge(t *testing.T) {
	cases := []struct {
		name string
		a    []int
		b    []int
	}{
		{"BothNil", nil, nil},
		{"BothEmpty", []int{}, []int{}},
		{"FirstEmpty", []int{}, []int{1, 2, 3}},
		{"SecondEmpty", []int{1, 2, 3}, []int{}},
		{"FirstNil", nil, []int{1}},
		{"SecondNil", []int{1}, nil},
		{"OneAndOneInOrder", []int{1}, []int{2}},
		{"OneAndOneOutOfOrder", []int{2}, []int{1}},
		{"OneAndOneEqual", []int{1}, []int{1}},
		{"Interleaved", []int{1, 3, 5}, []int{2, 4, 6}},
		{"InterleavedFromTheSecond", []int{2, 4, 6}, []int{1, 3, 5}},
		{"FirstAllSmaller", []int{1, 2, 3}, []int{4, 5, 6}},
		{"SecondAllSmaller", []int{4, 5, 6}, []int{1, 2, 3}},
		{"OneLongAndOneShort", []int{1, 2, 3, 4, 5, 6, 7}, []int{4}},
		{"ShortFirst", []int{4}, []int{1, 2, 3, 4, 5, 6, 7}},
		{"EqualNumbersInBoth", []int{1, 2, 2, 3}, []int{2, 2, 4}},
		{"AllNumbersEqual", []int{5, 5, 5}, []int{5, 5}},
		{"DuplicatesInOne", []int{1, 1, 1, 9}, []int{2}},
		{"Negatives", []int{-9, -5, -1}, []int{-7, -3}},
		{"MixedSigns", []int{-5, 0, 4}, []int{-2, -1, 7}},
		{"ExtremeValues", []int{math.MinInt, 0}, []int{0, math.MaxInt}},
		{"LongTailInTheFirst", []int{1, 2, 3, 4, 5}, []int{0}},
		{"LongTailInTheSecond", []int{0}, []int{1, 2, 3, 4, 5}},
		{"SameNumberAtTheEnd", []int{1, 7}, []int{3, 7}},
		{"SameNumberAtTheStart", []int{1, 7}, []int{1, 9}},
	}

	for _, m := range merges {
		t.Run(m.name, func(t *testing.T) {
			for _, c := range cases {
				t.Run(c.name, func(t *testing.T) {
					checkMerge(t, m.name, m.merge, c.a, c.b)
				})
			}
		})
	}
}

// TestMergeEveryShortPair walks every pair of sorted arrays of up to 4 numbers
// taken from 0, 1 and 2. It catches an index that stops one step too early and
// a wrong order for equal numbers.
func TestMergeEveryShortPair(t *testing.T) {
	const maxLen = 4

	vals := []int{0, 1, 2}

	for _, m := range merges {
		t.Run(m.name, func(t *testing.T) {
			failed := false

			for na := 0; na <= maxLen; na++ {
				for nb := 0; nb <= maxLen; nb++ {
					everyArray(na, vals, func(a []int) {
						if failed || !slices.IsSorted(a) {
							return
						}
						as := slices.Clone(a)

						everyArray(nb, vals, func(b []int) {
							if failed || !slices.IsSorted(b) {
								return
							}
							if !checkMerge(t, m.name, m.merge, as, b) {
								failed = true
							}
						})
					})
				}
			}
		})
	}
}

// TestMergeRandom merges random sorted arrays of many sizes and value ranges.
func TestMergeRandom(t *testing.T) {
	maxes := []int{1, 3, 20, 1000000}

	for _, m := range merges {
		t.Run(m.name, func(t *testing.T) {
			rnd := rand.New(rand.NewSource(1))

			for _, max := range maxes {
				for range 200 {
					a := sortedRandom(rnd.Intn(100), max, rnd)
					b := sortedRandom(rnd.Intn(100), max, rnd)

					if !checkMerge(t, m.name, m.merge, a, b) {
						return
					}
				}
			}
		})
	}
}

// TestMergeKeepsNeighbours gives parts of larger arrays to the merge. A merge
// that writes outside of a or b fails here.
func TestMergeKeepsNeighbours(t *testing.T) {
	for _, m := range merges {
		t.Run(m.name, func(t *testing.T) {
			arr := []int{111, 1, 3, 5, 222, 2, 4, 6, 333}
			a, b := arr[1:4], arr[5:8]

			got := m.merge(a, b)

			if want := []int{1, 2, 3, 4, 5, 6}; !slices.Equal(got, want) {
				t.Errorf("%s(%v, %v) = %v, want %v", m.name, a, b, got, want)
			}
			if want := []int{111, 1, 3, 5, 222, 2, 4, 6, 333}; !slices.Equal(arr, want) {
				t.Errorf("%s changed the larger array to %v, want %v", m.name, arr, want)
			}
		})
	}
}

// TestMergeLarge merges large arrays.
func TestMergeLarge(t *testing.T) {
	const n = 100000

	rnd := rand.New(rand.NewSource(2))

	for _, m := range merges {
		t.Run(m.name, func(t *testing.T) {
			a := sortedRandom(n, 1000000, rnd)
			b := sortedRandom(n, 1000000, rnd)

			want := make([]int, 0, 2*n)
			want = append(want, a...)
			want = append(want, b...)
			slices.Sort(want)

			got := m.merge(a, b)

			if !slices.Equal(got, want) {
				t.Errorf("%s gave a wrong answer for two arrays of %d numbers", m.name, n)
			}
		})
	}
}

func BenchmarkMerge(b *testing.B) {
	const n = 10000

	rnd := rand.New(rand.NewSource(3))

	pairs := []struct {
		name string
		a    []int
		b    []int
	}{
		{"Interleaved", sortedRandom(n, 1000000, rnd), sortedRandom(n, 1000000, rnd)},
		{"FirstAllSmaller", sortedRandom(n, 100, rnd), sortedRandom(n, 1000000, rnd)},
		{"FewValues", sortedRandom(n, 3, rnd), sortedRandom(n, 3, rnd)},
		{"OneShort", sortedRandom(n, 1000000, rnd), sortedRandom(4, 1000000, rnd)},
	}

	for _, m := range merges {
		b.Run(m.name, func(b *testing.B) {
			for _, p := range pairs {
				b.Run(p.name, func(b *testing.B) {
					for b.Loop() {
						m.merge(p.a, p.b)
					}
				})
			}
		})
	}
}
