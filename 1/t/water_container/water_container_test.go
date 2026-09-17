package watercontainer

import (
	"math/rand"
	"slices"
	"testing"
)

// brute gives the answer of the task in the slow and plain way: it tries every
// pair of lines and keeps the largest area. The tests compare the solution
// with it.
func brute(height []int) int {
	best := 0

	for i := range height {
		for j := i + 1; j < len(height); j++ {
			h := height[i]
			if height[j] < h {
				h = height[j]
			}

			if area := h * (j - i); area > best {
				best = area
			}
		}
	}

	return best
}

// check calls maxArea on a copy of height and compares the answer with want.
// It also reads the copy back, because maxArea must leave the array as it was.
func check(t *testing.T, height []int, want int) {
	t.Helper()

	got := slices.Clone(height)

	if n := maxArea(got); n != want {
		t.Errorf("maxArea(%v) = %d, want %d", height, n, want)
	}
	if !slices.Equal(got, height) {
		t.Errorf("maxArea(%v) changed the array to %v", height, got)
	}
}

// everyArray calls f once for every array of n numbers taken from vals.
func everyArray(n int, vals []int, f func([]int)) {
	height := make([]int, n)

	var build func(i int)
	build = func(i int) {
		if i == n {
			f(height)
			return
		}
		for _, v := range vals {
			height[i] = v
			build(i + 1)
		}
	}

	build(0)
}

// randomHeights gives n heights between 0 and max.
func randomHeights(n, max int, rnd *rand.Rand) []int {
	height := make([]int, n)
	for i := range height {
		height[i] = rnd.Intn(max + 1)
	}

	return height
}

func TestMaxArea(t *testing.T) {
	cases := []struct {
		name   string
		height []int
		want   int
	}{
		{"Example1", []int{1, 8, 6, 2, 5, 4, 8, 3, 7}, 49},
		{"Example2", []int{1, 1}, 1},
		{"TwoLinesDifferent", []int{1, 2}, 1},
		{"TwoLinesTallThenShort", []int{5, 1}, 1},
		{"TwoLinesWithZero", []int{0, 5}, 0},
		{"TwoLinesZero", []int{0, 0}, 0},
		{"AllZero", []int{0, 0, 0, 0}, 0},
		{"AllTheSame", []int{3, 3, 3, 3}, 9},
		{"Increasing", []int{1, 2, 3, 4, 5}, 6},
		{"Decreasing", []int{5, 4, 3, 2, 1}, 6},
		{"TallEndsLowMiddle", []int{9, 1, 1, 1, 9}, 36},
		{"LowEndsTallMiddle", []int{1, 9, 9, 9, 1}, 18},
		{"PeakInTheMiddle", []int{1, 2, 3, 2, 1}, 4},
		{"ValleyInTheMiddle", []int{3, 2, 1, 2, 3}, 12},
		{"OneTallLine", []int{10, 1, 1, 1, 1}, 4},
		{"LastLineTall", []int{1, 1, 1, 1, 10}, 4},
		{"ZeroBetweenTallLines", []int{4, 0, 4}, 8},
		{"ZeroAtTheStart", []int{0, 3, 3}, 3},
		{"ZeroAtTheEnd", []int{3, 3, 0}, 3},
		{"BestPairIsNeighbours", []int{1, 9, 9, 1}, 9},
		{"BestPairIsTheEnds", []int{9, 1, 1, 9}, 27},
		{"EqualTallLinesFarApart", []int{8, 2, 2, 2, 8}, 32},
		{"Plateau", []int{2, 2, 2, 2, 2, 2}, 10},
		{"SawTooth", []int{1, 5, 1, 5, 1, 5}, 20},
		{"LargestHeights", []int{10000, 10000}, 10000},
		{"LargestHeightAndZero", []int{10000, 0}, 0},
		{"SmallestArray", []int{4, 7}, 4},
		{"WideAndFlat", []int{1, 0, 0, 0, 0, 0, 0, 1}, 7},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			check(t, c.height, c.want)
		})
	}
}

// TestMaxAreaEveryShortArray walks every array of up to 9 heights taken from
// 0, 1 and 2. It catches a solution that moves the wrong pointer when the two
// heights are the same, and a solution that stops one step too early.
func TestMaxAreaEveryShortArray(t *testing.T) {
	const maxLen = 9

	vals := []int{0, 1, 2}

	for n := 2; n <= maxLen; n++ {
		failed := false

		everyArray(n, vals, func(height []int) {
			if failed {
				return
			}

			got := maxArea(slices.Clone(height))

			if want := brute(height); got != want {
				t.Errorf("maxArea(%v) = %d, want %d", height, got, want)
				failed = true
			}
		})
	}
}

// TestMaxAreaRandom compares the solution with brute on random arrays. A small
// range of heights gives many equal lines. A large range gives few.
func TestMaxAreaRandom(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	maxHeights := []int{1, 3, 10, 1000, 10000}

	for _, max := range maxHeights {
		for range 200 {
			height := randomHeights(2+rnd.Intn(200), max, rnd)

			got := maxArea(slices.Clone(height))

			if want := brute(height); got != want {
				t.Fatalf("maxArea(%v) = %d, want %d", height, got, want)
			}
		}
	}
}

// TestMaxAreaSortedRandom compares the solution with brute on sorted arrays.
// A sorted array holds no valley, so it leads a wrong solution to a different
// pair than a random array does.
func TestMaxAreaSortedRandom(t *testing.T) {
	rnd := rand.New(rand.NewSource(2))

	for range 200 {
		height := randomHeights(2+rnd.Intn(100), 1000, rnd)
		slices.Sort(height)

		got := maxArea(slices.Clone(height))

		if want := brute(height); got != want {
			t.Fatalf("maxArea(%v) = %d, want %d", height, got, want)
		}

		slices.Reverse(height)

		got = maxArea(slices.Clone(height))

		if want := brute(height); got != want {
			t.Fatalf("maxArea(%v) = %d, want %d", height, got, want)
		}
	}
}

// TestMaxAreaOneTallPairEverywhere puts one pair of tall lines at every pair of
// positions of an array of low lines. The best pair is the tall pair when the
// two tall lines stand far enough from each other.
func TestMaxAreaOneTallPairEverywhere(t *testing.T) {
	const n = 12

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			height := make([]int, n)
			for k := range height {
				height[k] = 1
			}
			height[i] = 100
			height[j] = 100

			check(t, height, brute(height))
		}
	}
}

// TestMaxAreaKeepsNeighbours gives a part of a larger array to maxArea. A
// solution that reads outside of height fails here.
func TestMaxAreaKeepsNeighbours(t *testing.T) {
	arr := []int{10000, 10000, 1, 8, 6, 2, 5, 4, 8, 3, 7, 10000, 10000}
	height := arr[2:11]

	if got, want := maxArea(height), 49; got != want {
		t.Errorf("maxArea(%v) = %d, want %d", height, got, want)
	}
}

// TestMaxAreaLarge uses the largest array of the constraints. A solution that
// tries every pair is too slow here.
func TestMaxAreaLarge(t *testing.T) {
	const n = 100000

	t.Run("AllTheSame", func(t *testing.T) {
		// Every line has the same height, so the ends give the largest area.
		height := make([]int, n)
		for i := range height {
			height[i] = 10000
		}

		if got, want := maxArea(height), 10000*(n-1); got != want {
			t.Errorf("maxArea of 100000 lines of height 10000 = %d, want %d", got, want)
		}
	})

	t.Run("TallEndsLowMiddle", func(t *testing.T) {
		// Only the two ends hold water, because the lines between them are 0.
		height := make([]int, n)
		height[0] = 10000
		height[n-1] = 10000

		if got, want := maxArea(height), 10000*(n-1); got != want {
			t.Errorf("maxArea of two tall ends = %d, want %d", got, want)
		}
	})

	t.Run("AllZero", func(t *testing.T) {
		height := make([]int, n)

		if got := maxArea(height); got != 0 {
			t.Errorf("maxArea of 100000 lines of height 0 = %d, want 0", got)
		}
	})

	t.Run("Increasing", func(t *testing.T) {
		// The heights grow from 1 to 10000, the largest height of the
		// constraints. The shorter line of a pair is always the left one, and
		// the right end stands furthest away, so the best pair ends at n - 1.
		// A plain loop over the left line gives the answer.
		height := make([]int, n)
		for i := range height {
			height[i] = i*10000/n + 1
		}

		want := 0
		for i := 0; i < n-1; i++ {
			if area := height[i] * (n - 1 - i); area > want {
				want = area
			}
		}

		if got := maxArea(height); got != want {
			t.Errorf("maxArea of 100000 growing lines = %d, want %d", got, want)
		}
	})

	t.Run("Random", func(t *testing.T) {
		// The answer comes from a second run on the reverse array. A container
		// keeps the same area when the array turns around.
		height := randomHeights(n, 10000, rand.New(rand.NewSource(3)))
		want := maxArea(slices.Clone(height))

		reverse := slices.Clone(height)
		slices.Reverse(reverse)

		if got := maxArea(reverse); got != want {
			t.Errorf("maxArea of the reverse array = %d, want %d", got, want)
		}
	})
}

// TestMaxAreaMedium compares the solution with brute on an array that brute
// still handles. It joins a plateau, a valley and a random part.
func TestMaxAreaMedium(t *testing.T) {
	const n = 2000

	rnd := rand.New(rand.NewSource(4))

	height := make([]int, 0, n)
	for i := 0; i < n/4; i++ {
		height = append(height, 5000)
	}
	for i := 0; i < n/4; i++ {
		height = append(height, i)
	}
	for i := 0; i < n/4; i++ {
		height = append(height, n/4-i)
	}
	height = append(height, randomHeights(n/4, 10000, rnd)...)

	got := maxArea(slices.Clone(height))

	if want := brute(height); got != want {
		t.Errorf("maxArea of the mixed array = %d, want %d", got, want)
	}
}

func BenchmarkMaxAreaSmall(b *testing.B) {
	height := []int{1, 8, 6, 2, 5, 4, 8, 3, 7}

	for b.Loop() {
		maxArea(height)
	}
}

// BenchmarkMaxAreaRandom gives the largest array of the constraints.
func BenchmarkMaxAreaRandom(b *testing.B) {
	const n = 100000

	height := randomHeights(n, 10000, rand.New(rand.NewSource(5)))

	b.ResetTimer()
	for b.Loop() {
		maxArea(height)
	}
}

// BenchmarkMaxAreaIncreasing gives growing heights, so a two pointer solution
// moves only the left pointer.
func BenchmarkMaxAreaIncreasing(b *testing.B) {
	const n = 100000

	height := make([]int, n)
	for i := range height {
		height[i] = i*10000/n + 1
	}

	b.ResetTimer()
	for b.Loop() {
		maxArea(height)
	}
}

// BenchmarkMaxAreaAllTheSame gives lines of one height, so every step compares
// two equal lines.
func BenchmarkMaxAreaAllTheSame(b *testing.B) {
	const n = 100000

	height := make([]int, n)
	for i := range height {
		height[i] = 10000
	}

	b.ResetTimer()
	for b.Loop() {
		maxArea(height)
	}
}
