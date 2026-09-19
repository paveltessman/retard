package minstack

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"testing"
	"time"
)

// The kinds of calls a test can make on the stack.
const (
	opPush = iota
	opPop
	opTop
	opMin
)

// op is one call on the stack. value belongs to opPush only.
type op struct {
	kind  int
	value int
}

func push(value int) op { return op{opPush, value} }

var (
	pop    = op{kind: opPop}
	top    = op{kind: opTop}
	getMin = op{kind: opMin}
)

func (o op) String() string {
	switch o.kind {
	case opPush:
		return fmt.Sprintf("push(%d)", o.value)
	case opPop:
		return "pop()"
	case opTop:
		return "top()"
	default:
		return "getMin()"
	}
}

// format writes a list of calls as one line, so an error message shows the
// whole run that led to the wrong answer.
func format(ops []op) string {
	parts := make([]string, len(ops))
	for i, o := range ops {
		parts[i] = o.String()
	}

	return strings.Join(parts, " ")
}

// refStack gives the answers of the task in the slow and plain way: it keeps
// every element in a slice, and it reads all of them to find the smallest one.
// The tests compare the solution with it.
type refStack struct {
	elems []int
}

func (r *refStack) Push(value int) { r.elems = append(r.elems, value) }

func (r *refStack) Pop() { r.elems = r.elems[:len(r.elems)-1] }

func (r *refStack) Top() int { return r.elems[len(r.elems)-1] }

func (r *refStack) GetMin() int {
	best := r.elems[0]
	for _, v := range r.elems[1:] {
		if v < best {
			best = v
		}
	}

	return best
}

func (r *refStack) Len() int { return len(r.elems) }

// run makes the calls of ops on the solution and on refStack, and compares
// every answer. With checkEvery, it also calls Top and GetMin after each call,
// so a wrong answer shows up at the call that made it. It gives back false on
// the first error, so a loop over thousands of runs can stop.
func run(t *testing.T, ops []op, checkEvery bool) bool {
	t.Helper()

	s := NewMinStack()
	ref := &refStack{}

	for i, o := range ops {
		done := ops[:i+1]

		switch o.kind {
		case opPush:
			s.Push(o.value)
			ref.Push(o.value)
		case opPop:
			s.Pop()
			ref.Pop()
		case opTop:
			if got, want := s.Top(), ref.Top(); got != want {
				t.Errorf("%s = %d, want %d", format(done), got, want)
				return false
			}
		case opMin:
			if got, want := s.GetMin(), ref.GetMin(); got != want {
				t.Errorf("%s = %d, want %d", format(done), got, want)
				return false
			}
		}

		if !checkEvery || ref.Len() == 0 {
			continue
		}

		if got, want := s.Top(), ref.Top(); got != want {
			t.Errorf("%s top() = %d, want %d", format(done), got, want)
			return false
		}
		if got, want := s.GetMin(), ref.GetMin(); got != want {
			t.Errorf("%s getMin() = %d, want %d", format(done), got, want)
			return false
		}
	}

	return true
}

// everyRun calls f once for every list of up to maxLen calls over the given
// values plus pop. It never gives a pop on an empty stack, because the task
// does not allow such a call.
func everyRun(values []int, maxLen int, f func(ops []op) bool) {
	ops := make([]op, 0, maxLen)

	var build func(depth int) bool
	build = func(depth int) bool {
		if len(ops) == maxLen {
			return true
		}

		for _, v := range values {
			ops = append(ops, push(v))
			ok := f(ops) && build(depth+1)
			ops = ops[:len(ops)-1]

			if !ok {
				return false
			}
		}

		if depth == 0 {
			return true
		}

		ops = append(ops, pop)
		ok := f(ops) && build(depth-1)
		ops = ops[:len(ops)-1]

		return ok
	}

	build(0)
}

// randomRun gives a list of n calls. The values come from the range
// [-span, span], so a narrow span gives many equal elements and a wide span
// gives few. It never gives a call on an empty stack.
func randomRun(n, span int, rnd *rand.Rand) []op {
	ops := make([]op, 0, n)
	size := 0

	for range n {
		if size == 0 {
			ops = append(ops, push(rnd.Intn(2*span+1)-span))
			size++

			continue
		}

		switch rnd.Intn(4) {
		case 0:
			ops = append(ops, pop)
			size--
		case 1:
			ops = append(ops, top)
		case 2:
			ops = append(ops, getMin)
		default:
			ops = append(ops, push(rnd.Intn(2*span+1)-span))
			size++
		}
	}

	return ops
}

func TestMinStack(t *testing.T) {
	cases := []struct {
		name string
		ops  []op
	}{
		{"Example1", []op{push(-2), push(0), push(-3), getMin, pop, top, getMin}},
		{"OneElement", []op{push(7), top, getMin}},
		{"OneElementZero", []op{push(0), top, getMin}},
		{"OneNegative", []op{push(-5), top, getMin}},
		{"PushPopPush", []op{push(3), pop, push(9), top, getMin}},
		{"AllPositive", []op{push(4), push(7), push(5), getMin, top}},
		{"AllNegative", []op{push(-4), push(-7), push(-5), getMin, top}},
		{"Increasing", []op{push(1), push(2), push(3), getMin, top, pop, getMin, top}},
		{"Decreasing", []op{push(3), push(2), push(1), getMin, top, pop, getMin, top}},
		{"MinAtTheBottom", []op{push(1), push(5), push(9), pop, getMin, pop, getMin}},
		{"MinAtTheTop", []op{push(9), push(5), push(1), getMin, pop, getMin, pop, getMin}},
		{"MinInTheMiddle", []op{push(5), push(1), push(9), getMin, pop, getMin, pop, getMin}},
		{"MinPoppedAway", []op{push(4), push(2), getMin, pop, getMin}},
		{"TwoEqualMins", []op{push(2), push(2), getMin, pop, getMin, pop}},
		{"TwoEqualMinsApart", []op{push(2), push(5), push(2), getMin, pop, getMin, pop, getMin}},
		{"ThreeEqualMins", []op{push(1), push(1), push(1), pop, getMin, pop, getMin, pop}},
		{"EqualMinsAroundABigger", []op{push(-1), push(3), push(-1), pop, getMin, pop, getMin}},
		{"AllEqual", []op{push(6), push(6), push(6), getMin, top, pop, getMin, top, pop, getMin, top}},
		{"NewMinAfterPop", []op{push(5), push(3), pop, push(4), getMin, top}},
		{"MinRestoredTwice", []op{push(5), push(2), pop, getMin, push(1), getMin, pop, getMin}},
		{"EmptyAgainThenPush", []op{push(2), pop, push(8), getMin, top, pop, push(-8), getMin, top}},
		{"ManyPopsBackToOne", []op{push(9), push(8), push(7), push(6), pop, pop, pop, getMin, top}},
		{"TopDoesNotRemove", []op{push(1), push(2), top, top, top, getMin, top}},
		{"GetMinDoesNotRemove", []op{push(2), push(1), getMin, getMin, getMin, top, pop, top}},
		{"ZeroAndNegative", []op{push(0), push(-1), getMin, pop, getMin, push(-2), getMin}},
		{"MinIsZero", []op{push(3), push(0), push(4), getMin, pop, getMin}},
		{"LargestValue", []op{push(math.MaxInt32), top, getMin}},
		{"SmallestValue", []op{push(math.MinInt32), top, getMin}},
		{"BothEnds", []op{push(math.MaxInt32), push(math.MinInt32), getMin, top, pop, getMin, top}},
		{"BothEndsOtherOrder", []op{push(math.MinInt32), push(math.MaxInt32), getMin, top, pop, getMin, top}},
		{"BigValuesOnly", []op{push(math.MaxInt32), push(math.MaxInt32 - 1), getMin, pop, getMin}},
		{"SmallValuesOnly", []op{push(math.MinInt32), push(math.MinInt32 + 1), getMin, pop, getMin}},
		{"PositiveValuesOnly", []op{push(5), push(9), getMin, pop, getMin}},
		{"NegativeValuesOnly", []op{push(-5), push(-9), getMin, pop, getMin}},
		{"LongMixedRun", []op{
			push(2), push(-1), push(3), getMin, top, pop, getMin, push(-7), getMin,
			pop, getMin, pop, getMin, top, push(0), getMin, top, pop, pop, push(4), getMin, top,
		}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			run(t, c.ops, false)
		})
	}
}

// TestMinStackEveryShortRun walks every short run of pushes and pops over a
// small set of values, and compares the solution with refStack after every
// call. It catches a solution that loses the smallest element after a pop, a
// solution that keeps a smallest element that left the stack, and a solution
// that handles two equal elements wrong.
func TestMinStackEveryShortRun(t *testing.T) {
	sets := []struct {
		values []int
		maxLen int
	}{
		{[]int{0, 1, 2}, 8},
		{[]int{0, 1}, 11},
		{[]int{-2, 0, -3}, 8},
		{[]int{1, 2}, 11},
	}

	for _, set := range sets {
		everyRun(set.values, set.maxLen, func(ops []op) bool {
			return run(t, ops, true)
		})

		if t.Failed() {
			return
		}
	}
}

// TestMinStackRandom compares the solution with refStack on random runs. The
// spans differ, so some runs hold many equal elements and some hold few. The
// test checks the top and the smallest element after every call.
func TestMinStackRandom(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))

	for _, span := range []int{1, 2, 5, 50, 1000000000} {
		for range 300 {
			ops := randomRun(1+rnd.Intn(120), span, rnd)

			if !run(t, ops, true) {
				return
			}
		}
	}

	// A few long runs, where a wrong answer needs many calls to show up.
	for _, span := range []int{2, 30} {
		for range 5 {
			ops := randomRun(5000, span, rnd)

			if !run(t, ops, false) {
				return
			}
		}
	}
}

// TestMinStackIsPerObject uses two stacks at the same time. It catches a
// solution that keeps the elements in a place that both objects share.
func TestMinStackIsPerObject(t *testing.T) {
	a := NewMinStack()
	b := NewMinStack()

	a.Push(5)
	b.Push(-5)

	if got := a.GetMin(); got != 5 {
		t.Errorf("getMin() of the first stack = %d, want 5", got)
	}
	if got := b.GetMin(); got != -5 {
		t.Errorf("getMin() of the second stack = %d, want -5", got)
	}

	a.Push(3)
	b.Push(7)
	b.Pop()

	if got := a.Top(); got != 3 {
		t.Errorf("top() of the first stack = %d, want 3", got)
	}
	if got := b.Top(); got != -5 {
		t.Errorf("top() of the second stack = %d, want -5", got)
	}
	if got := a.GetMin(); got != 3 {
		t.Errorf("getMin() of the first stack = %d, want 3", got)
	}
	if got := b.GetMin(); got != -5 {
		t.Errorf("getMin() of the second stack = %d, want -5", got)
	}
}

// TestMinStackTwoRandomStacks drives two stacks with one random run each, and
// mixes the calls. It catches a solution where one object changes the other.
func TestMinStackTwoRandomStacks(t *testing.T) {
	rnd := rand.New(rand.NewSource(2))

	stacks := [2]*MinStack{NewMinStack(), NewMinStack()}
	refs := [2]*refStack{{}, {}}
	names := [2]string{"the first stack", "the second stack"}

	for range 20000 {
		k := rnd.Intn(2)
		s, ref := stacks[k], refs[k]

		if ref.Len() == 0 {
			v := rnd.Intn(41) - 20
			s.Push(v)
			ref.Push(v)

			continue
		}

		switch rnd.Intn(4) {
		case 0:
			s.Pop()
			ref.Pop()
		case 1:
			if got, want := s.Top(), ref.Top(); got != want {
				t.Fatalf("top() of %s = %d, want %d", names[k], got, want)
			}
		case 2:
			if got, want := s.GetMin(), ref.GetMin(); got != want {
				t.Fatalf("getMin() of %s = %d, want %d", names[k], got, want)
			}
		default:
			v := rnd.Intn(41) - 20
			s.Push(v)
			ref.Push(v)
		}
	}
}

// TestMinStackLarge makes the largest number of calls the constraints allow.
// The answers come from the way the runs are built, so refStack stays out of
// the way. Each shape puts the smallest element at another place.
func TestMinStackLarge(t *testing.T) {
	const n = 10000 // pushes, pops and queries together stay under 3 * 10^4

	t.Run("Decreasing", func(t *testing.T) {
		s := NewMinStack()

		for i := range n {
			s.Push(n - i)

			if got := s.GetMin(); got != n-i {
				t.Fatalf("getMin() after %d decreasing pushes = %d, want %d", i+1, got, n-i)
			}
		}

		for i := range n {
			if got := s.Top(); got != 1+i {
				t.Fatalf("top() after %d pops = %d, want %d", i, got, 1+i)
			}
			if got := s.GetMin(); got != 1+i {
				t.Fatalf("getMin() after %d pops = %d, want %d", i, got, 1+i)
			}

			s.Pop()
		}
	})

	t.Run("Increasing", func(t *testing.T) {
		s := NewMinStack()

		for i := range n {
			s.Push(i)

			if got := s.GetMin(); got != 0 {
				t.Fatalf("getMin() after %d increasing pushes = %d, want 0", i+1, got)
			}
		}

		for i := n - 1; i >= 0; i-- {
			if got := s.Top(); got != i {
				t.Fatalf("top() with %d elements left = %d, want %d", i+1, got, i)
			}
			if got := s.GetMin(); got != 0 {
				t.Fatalf("getMin() with %d elements left = %d, want 0", i+1, got)
			}

			s.Pop()
		}
	})

	t.Run("EqualValues", func(t *testing.T) {
		s := NewMinStack()

		for range n {
			s.Push(-1)
		}

		for i := range n {
			if got := s.GetMin(); got != -1 {
				t.Fatalf("getMin() after %d pops of the same value = %d, want -1", i, got)
			}
			if got := s.Top(); got != -1 {
				t.Fatalf("top() after %d pops of the same value = %d, want -1", i, got)
			}

			s.Pop()
		}
	})

	t.Run("PushAndPopAgain", func(t *testing.T) {
		s := NewMinStack()
		s.Push(0)

		// Every round puts a smaller element on the stack and takes it away
		// again. The smallest element goes back to 0 after each round.
		for i := 1; i <= n/2; i++ {
			s.Push(-i)

			if got := s.GetMin(); got != -i {
				t.Fatalf("getMin() in round %d = %d, want %d", i, got, -i)
			}

			s.Pop()

			if got := s.GetMin(); got != 0 {
				t.Fatalf("getMin() after round %d = %d, want 0", i, got)
			}
		}
	})
}

// runWithLimit drives the stack in another goroutine and waits for the limit.
// It gives back ok = false when the run does not end in time.
func runWithLimit(n int, limit time.Duration) (ok bool, err error) {
	done := make(chan error, 1)

	go func() {
		done <- longRun(n)
	}()

	select {
	case err := <-done:
		return true, err
	case <-time.After(limit):
		return false, nil
	}
}

// longRun pushes n elements, then reads and pops all of them. The values fall
// until the middle of the run and rise again, so the smallest element changes
// with almost every call of the first half. The answers come from the way the
// run is built.
func longRun(n int) error {
	s := NewMinStack()

	value := func(i int) int {
		if i < n/2 {
			return -i
		}

		return -(n - 1 - i)
	}

	for i := range n {
		s.Push(value(i))

		want := -min(i, n/2-1)
		if got := s.GetMin(); got != want {
			return fmt.Errorf("getMin() after %d pushes = %d, want %d", i+1, got, want)
		}
	}

	for i := n - 1; i >= 0; i-- {
		if got := s.Top(); got != value(i) {
			return fmt.Errorf("top() with %d elements left = %d, want %d", i+1, got, value(i))
		}

		want := -min(i, n/2-1)
		if got := s.GetMin(); got != want {
			return fmt.Errorf("getMin() with %d elements left = %d, want %d", i+1, got, want)
		}

		s.Pop()
	}

	return nil
}

// TestMinStackScales makes far more calls than the constraints allow. A
// solution that answers every call in constant time needs a moment here. A
// solution that reads the elements of the stack to answer a call needs hours,
// so it runs into the limit.
func TestMinStackScales(t *testing.T) {
	if testing.Short() {
		t.Skip("the test needs a very long run")
	}

	const n = 300000
	const limit = 10 * time.Second

	ok, err := runWithLimit(n, limit)

	if !ok {
		t.Fatalf("a run of %d pushes, pops and queries needed more than %v, so a call is not constant in time", n, limit)
	}
	if err != nil {
		t.Fatal(err)
	}
}

func BenchmarkMinStack(b *testing.B) {
	const n = 100000

	b.Run("Push", func(b *testing.B) {
		for b.Loop() {
			s := NewMinStack()
			for i := range n {
				s.Push(n - i)
			}
		}
	})

	b.Run("PushAndPop", func(b *testing.B) {
		for b.Loop() {
			s := NewMinStack()
			for i := range n {
				s.Push(i % 100)
			}
			for range n {
				s.Pop()
			}
		}
	})

	b.Run("GetMin", func(b *testing.B) {
		s := NewMinStack()
		for i := range n {
			s.Push(n - i)
		}

		b.ResetTimer()

		for b.Loop() {
			s.GetMin()
		}
	})

	b.Run("Top", func(b *testing.B) {
		s := NewMinStack()
		for i := range n {
			s.Push(i)
		}

		b.ResetTimer()

		for b.Loop() {
			s.Top()
		}
	})

	b.Run("MixedRun", func(b *testing.B) {
		ops := randomRun(n, 30, rand.New(rand.NewSource(3)))

		b.ResetTimer()

		for b.Loop() {
			s := NewMinStack()
			for _, o := range ops {
				switch o.kind {
				case opPush:
					s.Push(o.value)
				case opPop:
					s.Pop()
				case opTop:
					s.Top()
				default:
					s.GetMin()
				}
			}
		}
	})
}
