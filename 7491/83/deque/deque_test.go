package deque

import (
	"math/rand"
	"testing"
)

// RunDequeTests checks one Deque implementation against the interface contract.
// newDeque must give back an empty deque on every call.
// values must hold at least 3 distinct items to push.
func RunDequeTests[T comparable](t *testing.T, newDeque func() Deque[T], values []T) {
	t.Helper()

	if len(values) < 3 {
		t.Fatalf("RunDequeTests needs at least 3 values, got %d", len(values))
	}
	var zero T

	t.Run("EmptyDeque", func(t *testing.T) {
		d := newDeque()

		if got := d.Len(); got != 0 {
			t.Errorf("Len() on a new deque = %d, want 0", got)
		}
		if item, ok := d.PopFront(); ok || item != zero {
			t.Errorf("PopFront() on a new deque = (%v, %t), want (%v, false)", item, ok, zero)
		}
		if item, ok := d.PopBack(); ok || item != zero {
			t.Errorf("PopBack() on a new deque = (%v, %t), want (%v, false)", item, ok, zero)
		}
		if got := d.Len(); got != 0 {
			t.Errorf("Len() after two failed pops = %d, want 0", got)
		}
	})

	t.Run("SingleItemEachEnd", func(t *testing.T) {
		// One item is both the front and the back, so either pop must find it.
		cases := []struct {
			name string
			push func(Deque[T], T)
			pop  func(Deque[T]) (T, bool)
		}{
			{"PushFrontPopFront", Deque[T].PushFront, Deque[T].PopFront},
			{"PushFrontPopBack", Deque[T].PushFront, Deque[T].PopBack},
			{"PushBackPopFront", Deque[T].PushBack, Deque[T].PopFront},
			{"PushBackPopBack", Deque[T].PushBack, Deque[T].PopBack},
		}

		for _, c := range cases {
			d := newDeque()
			c.push(d, values[0])
			if got := d.Len(); got != 1 {
				t.Errorf("%s: Len() after one push = %d, want 1", c.name, got)
			}
			if item, ok := c.pop(d); !ok || item != values[0] {
				t.Errorf("%s: pop = (%v, %t), want (%v, true)", c.name, item, ok, values[0])
			}
			if got := d.Len(); got != 0 {
				t.Errorf("%s: Len() after the pop = %d, want 0", c.name, got)
			}
		}
	})

	t.Run("PushBackPopFrontIsFIFO", func(t *testing.T) {
		d := newDeque()
		for i, v := range values {
			d.PushBack(v)
			if got := d.Len(); got != i+1 {
				t.Fatalf("Len() after %d pushes = %d, want %d", i+1, got, i+1)
			}
		}

		for i, want := range values {
			item, ok := d.PopFront()
			if !ok || item != want {
				t.Fatalf("PopFront() at position %d = (%v, %t), want (%v, true)", i, item, ok, want)
			}
			if got, want := d.Len(), len(values)-i-1; got != want {
				t.Fatalf("Len() after PopFront at position %d = %d, want %d", i, got, want)
			}
		}
	})

	t.Run("PushFrontPopBackIsFIFO", func(t *testing.T) {
		d := newDeque()
		for _, v := range values {
			d.PushFront(v)
		}

		for i, want := range values {
			item, ok := d.PopBack()
			if !ok || item != want {
				t.Fatalf("PopBack() at position %d = (%v, %t), want (%v, true)", i, item, ok, want)
			}
		}
		if got := d.Len(); got != 0 {
			t.Errorf("Len() after the drain = %d, want 0", got)
		}
	})

	t.Run("PushBackPopBackIsLIFO", func(t *testing.T) {
		d := newDeque()
		for _, v := range values {
			d.PushBack(v)
		}

		for i := len(values) - 1; i >= 0; i-- {
			item, ok := d.PopBack()
			if !ok || item != values[i] {
				t.Fatalf("PopBack() at position %d = (%v, %t), want (%v, true)", i, item, ok, values[i])
			}
		}
		if got := d.Len(); got != 0 {
			t.Errorf("Len() after the drain = %d, want 0", got)
		}
	})

	t.Run("PushFrontPopFrontIsLIFO", func(t *testing.T) {
		d := newDeque()
		for _, v := range values {
			d.PushFront(v)
		}

		for i := len(values) - 1; i >= 0; i-- {
			item, ok := d.PopFront()
			if !ok || item != values[i] {
				t.Fatalf("PopFront() at position %d = (%v, %t), want (%v, true)", i, item, ok, values[i])
			}
		}
		if got := d.Len(); got != 0 {
			t.Errorf("Len() after the drain = %d, want 0", got)
		}
	})

	t.Run("BothEndsTogether", func(t *testing.T) {
		// Build values[2], values[0], values[1], values[3] from the middle out.
		d := newDeque()
		d.PushBack(values[0])
		d.PushFront(values[2])
		d.PushBack(values[1])
		d.PushBack(values[0])
		d.PushFront(values[1])

		want := []T{values[1], values[2], values[0], values[1], values[0]}
		if got := d.Len(); got != len(want) {
			t.Fatalf("Len() after 5 pushes = %d, want %d", got, len(want))
		}

		// Pop from both ends and meet in the middle.
		for len(want) > 0 {
			item, ok := d.PopFront()
			if !ok || item != want[0] {
				t.Fatalf("PopFront() = (%v, %t), want (%v, true)", item, ok, want[0])
			}
			want = want[1:]
			if len(want) == 0 {
				break
			}

			last := len(want) - 1
			item, ok = d.PopBack()
			if !ok || item != want[last] {
				t.Fatalf("PopBack() = (%v, %t), want (%v, true)", item, ok, want[last])
			}
			want = want[:last]
		}
		if got := d.Len(); got != 0 {
			t.Errorf("Len() after the drain = %d, want 0", got)
		}
	})

	t.Run("InterleavedWrapAround", func(t *testing.T) {
		// A ring buffer must hold the order when the head and the tail wrap.
		d := newDeque()
		model := make([]T, 0, 64)

		for round := 0; round < 5; round++ {
			for i := 0; i < 2; i++ {
				v := values[(round+i)%len(values)]
				d.PushBack(v)
				model = append(model, v)
			}
			v := values[round%len(values)]
			d.PushFront(v)
			model = append([]T{v}, model...)

			want := model[0]
			model = model[1:]
			if item, ok := d.PopFront(); !ok || item != want {
				t.Fatalf("round %d: PopFront() = (%v, %t), want (%v, true)", round, item, ok, want)
			}

			last := len(model) - 1
			want = model[last]
			model = model[:last]
			if item, ok := d.PopBack(); !ok || item != want {
				t.Fatalf("round %d: PopBack() = (%v, %t), want (%v, true)", round, item, ok, want)
			}

			if got := d.Len(); got != len(model) {
				t.Fatalf("round %d: Len() = %d, want %d", round, got, len(model))
			}
		}

		for i, want := range model {
			if item, ok := d.PopFront(); !ok || item != want {
				t.Fatalf("drain %d: PopFront() = (%v, %t), want (%v, true)", i, item, ok, want)
			}
		}
		if got := d.Len(); got != 0 {
			t.Errorf("Len() after the drain = %d, want 0", got)
		}
	})

	t.Run("ReuseAfterDrain", func(t *testing.T) {
		d := newDeque()
		for _, v := range values {
			d.PushBack(v)
		}
		for range values {
			if _, ok := d.PopFront(); !ok {
				t.Fatal("PopFront() on a full deque returned false")
			}
		}

		if item, ok := d.PopFront(); ok || item != zero {
			t.Errorf("PopFront() on a drained deque = (%v, %t), want (%v, false)", item, ok, zero)
		}
		if item, ok := d.PopBack(); ok || item != zero {
			t.Errorf("PopBack() on a drained deque = (%v, %t), want (%v, false)", item, ok, zero)
		}
		// The deque must work again after it runs empty.
		d.PushFront(values[1])
		if item, ok := d.PopBack(); !ok || item != values[1] {
			t.Errorf("PopBack() after a reuse PushFront = (%v, %t), want (%v, true)", item, ok, values[1])
		}
		if got := d.Len(); got != 0 {
			t.Errorf("Len() after the reuse pop = %d, want 0", got)
		}
	})

	t.Run("DuplicateItems", func(t *testing.T) {
		d := newDeque()
		const count = 4
		for i := 0; i < count; i++ {
			d.PushBack(values[0])
			d.PushFront(values[0])
		}
		if got := d.Len(); got != 2*count {
			t.Errorf("Len() after %d equal pushes = %d, want %d", 2*count, got, 2*count)
		}
		for i := 0; i < count; i++ {
			if item, ok := d.PopFront(); !ok || item != values[0] {
				t.Fatalf("PopFront() number %d = (%v, %t), want (%v, true)", i+1, item, ok, values[0])
			}
			if item, ok := d.PopBack(); !ok || item != values[0] {
				t.Fatalf("PopBack() number %d = (%v, %t), want (%v, true)", i+1, item, ok, values[0])
			}
		}
		if got := d.Len(); got != 0 {
			t.Errorf("Len() after all pops = %d, want 0", got)
		}
	})

	t.Run("Growth", func(t *testing.T) {
		// A large run makes a fixed start capacity grow, and makes a ring buffer wrap.
		d := newDeque()
		const count = 1000
		for i := 0; i < count; i++ {
			if i%2 == 0 {
				d.PushBack(values[i%len(values)])
			} else {
				d.PushFront(values[i%len(values)])
			}
		}
		if got := d.Len(); got != count {
			t.Fatalf("Len() after %d pushes = %d, want %d", count, got, count)
		}

		// The odd pushes went to the front, so they come back in reverse order.
		for i := count - 1; i >= 1; i -= 2 {
			want := values[i%len(values)]
			if item, ok := d.PopFront(); !ok || item != want {
				t.Fatalf("PopFront() for push %d = (%v, %t), want (%v, true)", i, item, ok, want)
			}
		}
		// The even pushes went to the back, so they come back in push order.
		for i := 0; i < count; i += 2 {
			want := values[i%len(values)]
			if item, ok := d.PopFront(); !ok || item != want {
				t.Fatalf("PopFront() for push %d = (%v, %t), want (%v, true)", i, item, ok, want)
			}
		}
		if got := d.Len(); got != 0 {
			t.Errorf("Len() after all pops = %d, want 0", got)
		}
	})

	t.Run("RandomOperations", func(t *testing.T) {
		d := newDeque()
		model := make([]T, 0, 64)
		rng := rand.New(rand.NewSource(1))

		for step := 0; step < 2000; step++ {
			switch rng.Intn(4) {
			case 0:
				v := values[rng.Intn(len(values))]
				d.PushFront(v)
				model = append([]T{v}, model...)
			case 1:
				v := values[rng.Intn(len(values))]
				d.PushBack(v)
				model = append(model, v)
			case 2:
				item, ok := d.PopFront()
				if len(model) == 0 {
					if ok || item != zero {
						t.Fatalf("step %d: PopFront() on an empty deque = (%v, %t), want (%v, false)", step, item, ok, zero)
					}
					break
				}
				want := model[0]
				model = model[1:]
				if !ok || item != want {
					t.Fatalf("step %d: PopFront() = (%v, %t), want (%v, true)", step, item, ok, want)
				}
			case 3:
				item, ok := d.PopBack()
				if len(model) == 0 {
					if ok || item != zero {
						t.Fatalf("step %d: PopBack() on an empty deque = (%v, %t), want (%v, false)", step, item, ok, zero)
					}
					break
				}
				last := len(model) - 1
				want := model[last]
				model = model[:last]
				if !ok || item != want {
					t.Fatalf("step %d: PopBack() = (%v, %t), want (%v, true)", step, item, ok, want)
				}
			}

			if got := d.Len(); got != len(model) {
				t.Fatalf("step %d: Len() = %d, want %d", step, got, len(model))
			}
		}
	})

	t.Run("IndependentInstances", func(t *testing.T) {
		a := newDeque()
		b := newDeque()
		a.PushBack(values[0])

		if got := b.Len(); got != 0 {
			t.Errorf("Len() of a second deque = %d, want 0", got)
		}
		b.PushFront(values[1])
		if item, ok := a.PopFront(); !ok || item != values[0] {
			t.Errorf("PopFront() of the first deque = (%v, %t), want (%v, true)", item, ok, values[0])
		}
	})
}

func TestRingDequeInt(t *testing.T) {
	RunDequeTests(t, func() Deque[int] { return NewRingDeque[int]() }, []int{1, 2, 3, 42, -7})
}

func TestRingDequeString(t *testing.T) {
	RunDequeTests(t, func() Deque[string] { return NewRingDeque[string]() }, []string{"a", "b", "c", ""})
}
