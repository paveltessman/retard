package largestrect

import (
	"math/rand"
	"slices"
	"testing"
	"time"
)

// brute gives the answer of the task in the slow and plain way: it walks over
// every run of neighbour bars, keeps the lowest bar of the run, and takes the
// largest area. The tests compare the solution with it.
func brute(heights []int) int {
	best := 0

	for i := range heights {
		low := heights[i]

		for j := i; j < len(heights); j++ {
			low = min(low, heights[j])

			if area := low * (j - i + 1); area > best {
				best = area
			}
		}
	}

	return best
}

// check calls largestRectangleArea on a copy of heights and compares the answer
// with want. It also reads the copy back, because the input must stay as it
// was. It gives back false on the first error, so a loop over thousands of
// inputs can stop.
func check(t *testing.T, heights []int, want int) bool {
	t.Helper()

	got := slices.Clone(heights)
	ok := true

	if area := largestRectangleArea(got); area != want {
		t.Errorf("largestRectangleArea(%v) = %d, want %d", heights, area, want)
		ok = false
	}
	if !slices.Equal(got, heights) {
		t.Errorf("largestRectangleArea(%v) changed the input to %v", heights, got)
		ok = false
	}

	return ok
}

// checkLarge does the work of check for an input that is too large to print. It
// names the shape of the input instead. It gives back the time the solution
// needed, so a test can compare it with a limit.
func checkLarge(t *testing.T, name string, heights []int, want int) time.Duration {
	t.Helper()

	got := slices.Clone(heights)

	start := time.Now()
	area := largestRectangleArea(got)
	took := time.Since(start)

	if area != want {
		t.Fatalf("largestRectangleArea(%s of %d bars) = %d, want %d", name, len(heights), area, want)
	}
	if !slices.Equal(got, heights) {
		t.Fatalf("largestRectangleArea(%s of %d bars) changed the input", name, len(heights))
	}

	return took
}

// forEveryInput calls f with every histogram of n bars built from values. It
// stops and gives back false as soon as f gives back false.
func forEveryInput(values []int, n int, f func(heights []int) bool) bool {
	pick := make([]int, n)
	heights := make([]int, n)

	for {
		for i, k := range pick {
			heights[i] = values[k]
		}

		if !f(heights) {
			return false
		}

		p := n - 1
		for p >= 0 {
			pick[p]++
			if pick[p] < len(values) {
				break
			}

			pick[p] = 0
			p--
		}

		if p < 0 {
			return true
		}
	}
}

// TestLargestRectangleArea runs the examples of the task and the small shapes
// at the edges of the constraints: one bar, a bar of height 0, equal bars, bars
// that only rise, bars that only fall, and the largest height the constraints
// allow. Some cases hold a wide and short rectangle next to a narrow and tall
// one, so a solution that takes only one of the two gives a wrong answer.
func TestLargestRectangleArea(t *testing.T) {
	cases := []struct {
		name    string
		heights []int
		want    int
	}{
		{"Example1", []int{2, 1, 5, 6, 2, 3}, 10},
		{"Example2", []int{2, 4}, 4},
		{"OneBar", []int{5}, 5},
		{"OneZeroBar", []int{0}, 0},
		{"OneTallestBar", []int{10000}, 10000},
		{"TwoEqualBars", []int{3, 3}, 6},
		{"TwoBarsRising", []int{1, 2}, 2},
		{"TwoBarsFalling", []int{2, 1}, 2},
		{"TwoZeroBars", []int{0, 0}, 0},
		{"ZeroThenBar", []int{0, 5}, 5},
		{"BarThenZero", []int{5, 0}, 5},
		{"AllEqual", []int{4, 4, 4, 4}, 16},
		{"AllZero", []int{0, 0, 0}, 0},
		{"OnlyRising", []int{1, 2, 3, 4, 5}, 9},
		{"OnlyFalling", []int{5, 4, 3, 2, 1}, 9},
		{"ValleyBetweenTallBars", []int{6, 1, 6}, 6},
		{"ZeroBetweenTallBars", []int{6, 0, 6}, 6},
		{"TallBarBetweenShortBars", []int{1, 9, 1}, 9},
		{"PlateauInTheMiddle", []int{2, 5, 5, 5, 2}, 15},
		{"WideAndShortWins", []int{1, 1, 1, 1, 1, 4}, 6},
		{"NarrowAndTallWins", []int{1, 1, 10}, 10},
		{"ZigZagOverZero", []int{2, 0, 2, 0, 2}, 2},
		{"ShortRunOverTallBars", []int{5, 2, 5, 2, 5}, 10},
		{"StairsUpAndDown", []int{1, 2, 3, 2, 1}, 6},
		{"TallestBars", []int{10000, 10000}, 20000},
		{"TallestBarsAroundZero", []int{10000, 0, 10000}, 10000},
		{"LongRunOfOnesWithATallBar", []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 9}, 11},
		{"EqualRunThenTaller", []int{3, 3, 3, 4}, 12},
		{"TallerThenEqualRun", []int{4, 3, 3, 3}, 12},
		{"AnswerAtTheStart", []int{5, 5, 1}, 10},
		{"AnswerAtTheEnd", []int{1, 5, 5}, 10},
		{"ZerosBetweenRuns", []int{0, 3, 0, 3, 3, 0}, 6},
		{"RisingThenDrop", []int{2, 3, 4, 1}, 6},
		{"DropThenRising", []int{1, 4, 3, 2}, 6},
		{"TwoPlateaus", []int{2, 2, 2, 7, 7}, 14},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			check(t, c.heights, c.want)
		})
	}
}

// TestLargestRectangleAreaEveryShortInput walks every histogram of up to 12
// bars over small value sets and compares the solution with brute. The sets are
// small, so every plateau, every tie and every drop to 0 appears. It catches a
// solution that is off by one bar in the width of a rectangle.
func TestLargestRectangleAreaEveryShortInput(t *testing.T) {
	sets := []struct {
		values []int
		maxLen int
	}{
		{[]int{1, 2}, 12},
		{[]int{0, 1}, 12},
		{[]int{0, 1, 2}, 8},
		{[]int{1, 2, 3}, 8},
		{[]int{0, 2, 3, 10000}, 6},
		{[]int{9999, 10000}, 10},
	}

	for _, s := range sets {
		for n := 1; n <= s.maxLen; n++ {
			ok := forEveryInput(s.values, n, func(heights []int) bool {
				return check(t, heights, brute(heights))
			})

			if !ok {
				return
			}
		}
	}
}

// TestLargestRectangleAreaRandom compares the solution with brute on random
// histograms. The narrow ranges give many equal bars, the wide range gives few.
// The seed is fixed, so a failure comes back on the next run.
func TestLargestRectangleAreaRandom(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))

	ranges := []struct {
		low, high int
	}{
		{0, 10000},
		{0, 2},
		{0, 1},
		{1, 3},
		{7, 7},
		{9995, 10000},
	}

	for range 400 {
		for _, r := range ranges {
			n := 1 + rnd.Intn(60)

			heights := make([]int, n)
			for i := range heights {
				heights[i] = r.low + rnd.Intn(r.high-r.low+1)
			}

			if !check(t, heights, brute(heights)) {
				return
			}
		}
	}
}

// TestLargestRectangleAreaRandomShapes compares the solution with brute on
// random histograms that rise, fall, or hold long runs of equal bars. These
// shapes are rare in a uniform random array, and they are the ones that decide
// where a rectangle ends.
func TestLargestRectangleAreaRandomShapes(t *testing.T) {
	rnd := rand.New(rand.NewSource(2))

	for range 500 {
		n := 1 + rnd.Intn(40)
		heights := make([]int, n)

		switch rnd.Intn(4) {
		case 0: // Rising bars, with ties.
			h := 0
			for i := range heights {
				h += rnd.Intn(3)
				heights[i] = h
			}
		case 1: // Falling bars, with ties.
			h := 3 * n
			for i := range heights {
				h -= rnd.Intn(3)
				heights[i] = h
			}
		case 2: // Runs of equal bars.
			h := rnd.Intn(5)
			for i := range heights {
				if rnd.Intn(4) == 0 {
					h = rnd.Intn(5)
				}
				heights[i] = h
			}
		case 3: // A few tall bars in a field of zeros.
			for i := range heights {
				if rnd.Intn(5) == 0 {
					heights[i] = 1 + rnd.Intn(10000)
				}
			}
		}

		if !check(t, heights, brute(heights)) {
			return
		}
	}
}

// TestLargestRectangleAreaKeepsNeighbours gives a part of a larger array to the
// solution. The bars outside of the part are much taller than the bars inside,
// so a solution that reads outside of the input gives a larger area here.
func TestLargestRectangleAreaKeepsNeighbours(t *testing.T) {
	arr := []int{10000, 10000, 3, 1, 4, 2, 10000, 10000}

	middle := arr[2:6] // 3, 1, 4, 2
	check(t, middle, 4)

	head := arr[:3] // 10000, 10000, 3
	check(t, head, 20000)

	tail := arr[5:] // 2, 10000, 10000
	check(t, tail, 20000)

	if !slices.Equal(arr, []int{10000, 10000, 3, 1, 4, 2, 10000, 10000}) {
		t.Errorf("largestRectangleArea changed the array around the input to %v", arr)
	}
}

// TestLargestRectangleAreaLarge uses the largest input of the constraints. The
// answer comes from the way the input is built, not from brute, which is too
// slow at this size.
func TestLargestRectangleAreaLarge(t *testing.T) {
	const n = 100000

	// Every bar is as tall as the constraints allow, so the rectangle covers
	// the whole histogram.
	equal := make([]int, n)
	for i := range equal {
		equal[i] = 10000
	}

	checkLarge(t, "bars of equal height", equal, 10000*n)

	// The bars rise by one, so the bar at index i is the lowest of the run from
	// i to the last bar.
	rising := make([]int, n)
	risingWant := 0

	for i := range rising {
		rising[i] = i + 1

		if area := (i + 1) * (n - i); area > risingWant {
			risingWant = area
		}
	}

	checkLarge(t, "bars that rise by one", rising, risingWant)

	// The same histogram the other way round.
	falling := make([]int, n)
	for i := range falling {
		falling[i] = n - i
	}

	checkLarge(t, "bars that fall by one", falling, risingWant)

	// A tall bar, then a bar of height 0, again and again. No rectangle is
	// wider than one bar.
	zigzag := make([]int, n)
	for i := 1; i < n; i += 2 {
		zigzag[i] = 10000
	}

	checkLarge(t, "tall bars between bars of height 0", zigzag, 10000)

	// Bars of height 1 with one tall bar in the middle. The whole histogram
	// beats the tall bar, because it is wide.
	wide := make([]int, n)
	for i := range wide {
		wide[i] = 1
	}

	wide[n/2] = 10000
	checkLarge(t, "low bars with one tall bar", wide, n)
}

// TestLargestRectangleAreaIsFastOnHardInput runs the largest input of the
// constraints on the shapes that cost the most: bars of equal height, bars that
// rise, and bars that fall. A solution that looks at every run of bars needs
// far more than the limit here.
func TestLargestRectangleAreaIsFastOnHardInput(t *testing.T) {
	if testing.Short() {
		t.Skip("the test needs a very large input")
	}

	const n = 100000
	const limit = time.Second

	// The widest rectangle of a histogram that rises or falls by one.
	rampWant := 0
	for i := range n {
		if area := (i + 1) * (n - i); area > rampWant {
			rampWant = area
		}
	}

	shapes := []struct {
		name  string
		want  int
		build func(i int) int
	}{
		{"bars of equal height", 10000 * n, func(int) int { return 10000 }},
		{"bars that rise by one", rampWant, func(i int) int { return i + 1 }},
		{"bars that fall by one", rampWant, func(i int) int { return n - i }},
		{"bars of height 1", n, func(int) int { return 1 }},
	}

	for _, s := range shapes {
		heights := make([]int, n)
		for i := range heights {
			heights[i] = s.build(i)
		}

		took := checkLarge(t, s.name, heights, s.want)
		if took > limit {
			t.Fatalf("largestRectangleArea needed %v for %d %s, more than the limit of %v", took, n, s.name, limit)
		}
	}
}

func BenchmarkLargestRectangleArea(b *testing.B) {
	const n = 100000

	rnd := rand.New(rand.NewSource(3))

	shapes := []struct {
		name  string
		build func(i int) int
	}{
		{"Equal", func(int) int { return 10000 }},
		{"Rising", func(i int) int { return i + 1 }},
		{"Falling", func(i int) int { return n - i }},
		{"ZigZag", func(i int) int { return 10000 * (i % 2) }},
		{"Random", func(int) int { return rnd.Intn(10001) }},
	}

	for _, s := range shapes {
		heights := make([]int, n)
		for i := range heights {
			heights[i] = s.build(i)
		}

		b.Run(s.name, func(b *testing.B) {
			for b.Loop() {
				largestRectangleArea(heights)
			}
		})
	}
}
