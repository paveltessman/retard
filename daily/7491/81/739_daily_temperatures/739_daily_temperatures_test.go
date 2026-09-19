package dailytemperatures

import (
	"math/rand"
	"slices"
	"testing"
	"time"
)

// brute gives the answer of the task in the slow and plain way: for every day
// it walks over the later days and stops at the first warmer one. The tests
// compare the solution with it.
func brute(temperatures []int) []int {
	answer := make([]int, len(temperatures))

	for i, v := range temperatures {
		for j := i + 1; j < len(temperatures); j++ {
			if temperatures[j] > v {
				answer[i] = j - i
				break
			}
		}
	}

	return answer
}

// check calls dailyTemperatures on a copy of temperatures and compares the
// answer with want. It also reads the copy back, because the input must stay as
// it was. It gives back false on the first error, so a loop over thousands of
// inputs can stop.
func check(t *testing.T, temperatures, want []int) bool {
	t.Helper()

	got := slices.Clone(temperatures)
	ok := true

	if answer := dailyTemperatures(got); !slices.Equal(answer, want) {
		t.Errorf("dailyTemperatures(%v) = %v, want %v", temperatures, answer, want)
		ok = false
	}
	if !slices.Equal(got, temperatures) {
		t.Errorf("dailyTemperatures(%v) changed the input to %v", temperatures, got)
		ok = false
	}

	return ok
}

// checkLarge does the work of check for an input that is too large to print. It
// names the first day with a wrong answer. It gives back the time the solution
// needed, so a test can compare it with a limit.
func checkLarge(t *testing.T, name string, temperatures, want []int) time.Duration {
	t.Helper()

	got := slices.Clone(temperatures)

	start := time.Now()
	answer := dailyTemperatures(got)
	took := time.Since(start)

	if len(answer) != len(want) {
		t.Fatalf("dailyTemperatures(%s of %d days) gave %d answers, want %d", name, len(temperatures), len(answer), len(want))
	}

	for i := range want {
		if answer[i] != want[i] {
			t.Fatalf("dailyTemperatures(%s of %d days): answer[%d] = %d, want %d", name, len(temperatures), i, answer[i], want[i])
		}
	}

	if !slices.Equal(got, temperatures) {
		t.Fatalf("dailyTemperatures(%s of %d days) changed the input", name, len(temperatures))
	}

	return took
}

// forEveryInput calls f with every array of n days built from values. It stops
// and gives back false as soon as f gives back false.
func forEveryInput(values []int, n int, f func(temperatures []int) bool) bool {
	pick := make([]int, n)
	temperatures := make([]int, n)

	for {
		for i, k := range pick {
			temperatures[i] = values[k]
		}

		if !f(temperatures) {
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

// TestDailyTemperatures runs the examples of the task and the small shapes at
// the edges of the constraints: one day, two days, equal days, days that only
// rise, and days that only fall.
func TestDailyTemperatures(t *testing.T) {
	cases := []struct {
		name         string
		temperatures []int
		want         []int
	}{
		{"Example1", []int{73, 74, 75, 71, 69, 72, 76, 73}, []int{1, 1, 4, 2, 1, 1, 0, 0}},
		{"Example2", []int{30, 40, 50, 60}, []int{1, 1, 1, 0}},
		{"Example3", []int{30, 60, 90}, []int{1, 1, 0}},
		{"OneDay", []int{50}, []int{0}},
		{"OneColdestDay", []int{30}, []int{0}},
		{"OneWarmestDay", []int{100}, []int{0}},
		{"TwoDaysWarmer", []int{30, 31}, []int{1, 0}},
		{"TwoDaysColder", []int{31, 30}, []int{0, 0}},
		{"TwoDaysEqual", []int{30, 30}, []int{0, 0}},
		{"ColdestAndWarmest", []int{30, 100}, []int{1, 0}},
		{"WarmestAndColdest", []int{100, 30}, []int{0, 0}},
		{"AllEqual", []int{50, 50, 50, 50}, []int{0, 0, 0, 0}},
		{"AllColdest", []int{30, 30, 30}, []int{0, 0, 0}},
		{"AllWarmest", []int{100, 100, 100}, []int{0, 0, 0}},
		{"OnlyRising", []int{30, 31, 32, 33, 34}, []int{1, 1, 1, 1, 0}},
		{"OnlyFalling", []int{34, 33, 32, 31, 30}, []int{0, 0, 0, 0, 0}},
		{"WarmerAtTheLastDay", []int{50, 50, 50, 51}, []int{3, 2, 1, 0}},
		{"ColderAtTheLastDay", []int{40, 40, 40, 39}, []int{0, 0, 0, 0}},
		{"EqualRunThenWarmer", []int{30, 30, 30, 30, 30, 31}, []int{5, 4, 3, 2, 1, 0}},
		{"WarmDayAtTheStart", []int{100, 30, 31, 32}, []int{0, 1, 1, 0}},
		{"WarmDayAtTheEnd", []int{30, 31, 32, 100}, []int{1, 1, 1, 0}},
		{"ValleyBetweenEqualDays", []int{50, 30, 50}, []int{0, 1, 0}},
		{"FallThenOneWarmDay", []int{75, 74, 73, 72, 76}, []int{4, 3, 2, 1, 0}},
		{"LongWaitFromTheFirstDay", []int{99, 30, 30, 30, 30, 100}, []int{5, 4, 3, 2, 1, 0}},
		{"ZigZag", []int{30, 40, 30, 40, 30, 40}, []int{1, 0, 1, 0, 1, 0}},
		{"OverAPlateau", []int{50, 60, 60, 60, 70}, []int{1, 3, 2, 1, 0}},
		{"TwoPeaks", []int{70, 60, 65, 50, 80}, []int{4, 1, 2, 1, 0}},
		{"TheSameValueLater", []int{60, 50, 60, 70}, []int{3, 1, 1, 0}},
		{"EqualAfterTheWarmestDay", []int{31, 30, 30, 30}, []int{0, 0, 0, 0}},
		{"NarrowRange", []int{99, 100, 99, 100, 99}, []int{1, 0, 1, 0, 0}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			check(t, c.temperatures, c.want)
		})
	}
}

// TestDailyTemperaturesEveryShortInput walks every array of up to 12 days over
// small value sets and compares the solution with brute. The sets are small, so
// every plateau, every tie and every fall appears. It catches a solution that
// counts an equal day as warmer, and a solution that is off by one day.
func TestDailyTemperaturesEveryShortInput(t *testing.T) {
	sets := []struct {
		values []int
		maxLen int
	}{
		{[]int{50, 51}, 12},
		{[]int{30, 31, 32}, 7},
		{[]int{30, 40, 41, 100}, 6},
		{[]int{99, 100}, 10},
	}

	for _, s := range sets {
		for n := 1; n <= s.maxLen; n++ {
			ok := forEveryInput(s.values, n, func(temperatures []int) bool {
				return check(t, temperatures, brute(temperatures))
			})

			if !ok {
				return
			}
		}
	}
}

// TestDailyTemperaturesRandom compares the solution with brute on random
// arrays. The narrow ranges give many equal days, the wide range gives few. The
// seed is fixed, so a failure comes back on the next run.
func TestDailyTemperaturesRandom(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))

	ranges := []struct {
		low, high int
	}{
		{30, 100},
		{30, 32},
		{30, 31},
		{60, 60},
		{95, 100},
		{40, 45},
	}

	for range 400 {
		for _, r := range ranges {
			n := 1 + rnd.Intn(80)

			temperatures := make([]int, n)
			for i := range temperatures {
				temperatures[i] = r.low + rnd.Intn(r.high-r.low+1)
			}

			if !check(t, temperatures, brute(temperatures)) {
				return
			}
		}
	}
}

// TestDailyTemperaturesKeepsNeighbours gives a part of a larger array to the
// solution. The days outside of the part are warmer than the days inside, so a
// solution that reads outside of the input gives a wrong answer here.
func TestDailyTemperaturesKeepsNeighbours(t *testing.T) {
	arr := []int{100, 99, 50, 60, 55, 40, 100, 98}

	middle := arr[2:6] // 50, 60, 55, 40
	check(t, middle, []int{1, 0, 0, 0})

	head := arr[:3] // 100, 99, 50
	check(t, head, []int{0, 0, 0})

	if !slices.Equal(arr, []int{100, 99, 50, 60, 55, 40, 100, 98}) {
		t.Errorf("dailyTemperatures changed the array around the input to %v", arr)
	}
}

// TestDailyTemperaturesLarge uses the largest input of the constraints. The
// answer comes from the way the input is built, not from brute, which is too
// slow at this size.
func TestDailyTemperaturesLarge(t *testing.T) {
	const n = 100000

	// A rise of one degree per day from 30 to 100, again and again. Every day
	// under 100 waits one day. Every day at 100 waits for ever.
	rising := make([]int, n)
	risingWant := make([]int, n)

	for i := range rising {
		rising[i] = 30 + i%71
	}

	for i := range risingWant {
		if rising[i] < 100 && i+1 < n {
			risingWant[i] = 1
		}
	}

	checkLarge(t, "a rise of one degree per day", rising, risingWant)

	// A fall of one degree per day from 100 to 30, again and again. A day at
	// 100 waits for ever. Every other day waits for the next 100.
	falling := make([]int, n)
	fallingWant := make([]int, n)

	for i := range falling {
		falling[i] = 100 - i%71
	}

	for i := range fallingWant {
		wait := 71 - i%71
		if i%71 != 0 && i+wait < n {
			fallingWant[i] = wait
		}
	}

	checkLarge(t, "a fall of one degree per day", falling, fallingWant)

	// Pairs of equal days. The second day of a pair waits one day, the first
	// day of a pair waits two.
	pairs := make([]int, n)
	pairsWant := make([]int, n)

	for i := range pairs {
		pairs[i] = 30 + (i/2)%71
	}

	for i := range pairsWant {
		wait := 2 - i%2
		if pairs[i] < 100 && i+wait < n {
			pairsWant[i] = wait
		}
	}

	checkLarge(t, "pairs of equal days", pairs, pairsWant)
}

// TestDailyTemperaturesIsFastOnHardInput runs the hardest shapes: every day
// equal, one warm day at the end of a long cold run, and a fall that never
// turns. A solution that walks over the later days of every day needs far more
// than the limit here. The input is larger than the constraints allow, because
// a walk over 10^5 days is still fast on a quick machine.
func TestDailyTemperaturesIsFastOnHardInput(t *testing.T) {
	if testing.Short() {
		t.Skip("the test needs a very large input")
	}

	const n = 200000
	const limit = 2 * time.Second

	// Every day is equal, so no day is warmer than another one.
	equal := make([]int, n)
	for i := range equal {
		equal[i] = 50
	}

	took := checkLarge(t, "equal days", equal, make([]int, n))
	if took > limit {
		t.Fatalf("dailyTemperatures needed %v for %d equal days, more than the limit of %v", took, n, limit)
	}

	// One warm day closes a long cold run, so every day waits for the last one.
	cold := make([]int, n)
	coldWant := make([]int, n)

	for i := range cold {
		cold[i] = 30
	}

	cold[n-1] = 100
	for i := range n - 1 {
		coldWant[i] = n - 1 - i
	}

	took = checkLarge(t, "one warm day at the end", cold, coldWant)
	if took > limit {
		t.Fatalf("dailyTemperatures needed %v for %d days with one warm day at the end, more than the limit of %v", took, n, limit)
	}

	// The temperature falls from 100 to 30 over the whole array, with a plateau
	// of equal days at every step. No day is ever warmer than an earlier one.
	falling := make([]int, n)
	for i := range falling {
		falling[i] = 100 - (i*71)/n
	}

	took = checkLarge(t, "a slow fall", falling, make([]int, n))
	if took > limit {
		t.Fatalf("dailyTemperatures needed %v for a slow fall over %d days, more than the limit of %v", took, n, limit)
	}
}

func BenchmarkDailyTemperatures(b *testing.B) {
	const n = 100000

	rnd := rand.New(rand.NewSource(3))

	shapes := []struct {
		name  string
		build func(i int) int
	}{
		{"Rising", func(i int) int { return 30 + i%71 }},
		{"Falling", func(i int) int { return 100 - i%71 }},
		{"EqualDays", func(int) int { return 50 }},
		{"Random", func(int) int { return 30 + rnd.Intn(71) }},
		{"SlowFall", func(i int) int { return 100 - (i*71)/n }},
	}

	for _, s := range shapes {
		temperatures := make([]int, n)
		for i := range temperatures {
			temperatures[i] = s.build(i)
		}

		b.Run(s.name, func(b *testing.B) {
			for b.Loop() {
				dailyTemperatures(temperatures)
			}
		})
	}
}
