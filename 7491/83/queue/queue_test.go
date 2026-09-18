package queue

import (
	"math/rand"
	"testing"
)

// RunQueueTests checks one Queue implementation against the interface contract.
// newQueue must give back an empty queue on every call.
// values must hold at least 3 distinct items to enqueue.
func RunQueueTests[T comparable](t *testing.T, newQueue func() Queue[T], values []T) {
	t.Helper()

	if len(values) < 3 {
		t.Fatalf("RunQueueTests needs at least 3 values, got %d", len(values))
	}
	var zero T

	t.Run("EmptyQueue", func(t *testing.T) {
		q := newQueue()

		if got := q.Len(); got != 0 {
			t.Errorf("Len() on a new queue = %d, want 0", got)
		}
		if item, ok := q.Peek(); ok || item != zero {
			t.Errorf("Peek() on a new queue = (%v, %t), want (%v, false)", item, ok, zero)
		}
		if item, ok := q.Dequeue(); ok || item != zero {
			t.Errorf("Dequeue() on a new queue = (%v, %t), want (%v, false)", item, ok, zero)
		}
		if got := q.Len(); got != 0 {
			t.Errorf("Len() after a failed Dequeue = %d, want 0", got)
		}
	})

	t.Run("EnqueueThenPeek", func(t *testing.T) {
		q := newQueue()
		q.Enqueue(values[0])

		if got := q.Len(); got != 1 {
			t.Errorf("Len() after one Enqueue = %d, want 1", got)
		}
		// Peek must not remove the item, so two calls give the same answer.
		for call := 1; call <= 2; call++ {
			item, ok := q.Peek()
			if !ok || item != values[0] {
				t.Errorf("Peek() call %d = (%v, %t), want (%v, true)", call, item, ok, values[0])
			}
			if got := q.Len(); got != 1 {
				t.Errorf("Len() after Peek call %d = %d, want 1", call, got)
			}
		}
	})

	t.Run("FirstInFirstOut", func(t *testing.T) {
		q := newQueue()
		for i, v := range values {
			q.Enqueue(v)
			if got := q.Len(); got != i+1 {
				t.Fatalf("Len() after %d enqueues = %d, want %d", i+1, got, i+1)
			}
		}

		for i, want := range values {
			// Peek must always show the oldest item, not the newest one.
			if item, ok := q.Peek(); !ok || item != want {
				t.Fatalf("Peek() at position %d = (%v, %t), want (%v, true)", i, item, ok, want)
			}
			item, ok := q.Dequeue()
			if !ok || item != want {
				t.Fatalf("Dequeue() at position %d = (%v, %t), want (%v, true)", i, item, ok, want)
			}
			if got, want := q.Len(), len(values)-i-1; got != want {
				t.Fatalf("Len() after Dequeue at position %d = %d, want %d", i, got, want)
			}
		}
	})

	t.Run("InterleavedWrapAround", func(t *testing.T) {
		// A ring buffer must hold the order when the head and the tail wrap.
		q := newQueue()
		model := make([]T, 0, 64)

		for round := 0; round < 5; round++ {
			for i := 0; i < 3; i++ {
				v := values[(round+i)%len(values)]
				q.Enqueue(v)
				model = append(model, v)
			}
			for i := 0; i < 2; i++ {
				want := model[0]
				model = model[1:]
				item, ok := q.Dequeue()
				if !ok || item != want {
					t.Fatalf("round %d: Dequeue() = (%v, %t), want (%v, true)", round, item, ok, want)
				}
			}
			if got := q.Len(); got != len(model) {
				t.Fatalf("round %d: Len() = %d, want %d", round, got, len(model))
			}
		}

		for i, want := range model {
			if item, ok := q.Dequeue(); !ok || item != want {
				t.Fatalf("drain %d: Dequeue() = (%v, %t), want (%v, true)", i, item, ok, want)
			}
		}
		if got := q.Len(); got != 0 {
			t.Errorf("Len() after the drain = %d, want 0", got)
		}
	})

	t.Run("ReuseAfterDrain", func(t *testing.T) {
		q := newQueue()
		for _, v := range values {
			q.Enqueue(v)
		}
		for range values {
			if _, ok := q.Dequeue(); !ok {
				t.Fatal("Dequeue() on a full queue returned false")
			}
		}

		if item, ok := q.Dequeue(); ok || item != zero {
			t.Errorf("Dequeue() on a drained queue = (%v, %t), want (%v, false)", item, ok, zero)
		}
		// The queue must work again after it runs empty.
		q.Enqueue(values[1])
		if item, ok := q.Dequeue(); !ok || item != values[1] {
			t.Errorf("Dequeue() after a reuse Enqueue = (%v, %t), want (%v, true)", item, ok, values[1])
		}
		if got := q.Len(); got != 0 {
			t.Errorf("Len() after the reuse Dequeue = %d, want 0", got)
		}
	})

	t.Run("DuplicateItems", func(t *testing.T) {
		q := newQueue()
		const count = 4
		for i := 0; i < count; i++ {
			q.Enqueue(values[0])
		}
		if got := q.Len(); got != count {
			t.Errorf("Len() after %d equal enqueues = %d, want %d", count, got, count)
		}
		for i := 0; i < count; i++ {
			if item, ok := q.Dequeue(); !ok || item != values[0] {
				t.Fatalf("Dequeue() number %d = (%v, %t), want (%v, true)", i+1, item, ok, values[0])
			}
		}
		if got := q.Len(); got != 0 {
			t.Errorf("Len() after all dequeues = %d, want 0", got)
		}
	})

	t.Run("Growth", func(t *testing.T) {
		// A large run makes a fixed start capacity grow, and makes a ring buffer wrap.
		q := newQueue()
		const count = 1000
		for i := 0; i < count; i++ {
			q.Enqueue(values[i%len(values)])
		}
		if got := q.Len(); got != count {
			t.Fatalf("Len() after %d enqueues = %d, want %d", count, got, count)
		}
		for i := 0; i < count; i++ {
			want := values[i%len(values)]
			if item, ok := q.Dequeue(); !ok || item != want {
				t.Fatalf("Dequeue() number %d = (%v, %t), want (%v, true)", i+1, item, ok, want)
			}
		}
		if got := q.Len(); got != 0 {
			t.Errorf("Len() after all dequeues = %d, want 0", got)
		}
	})

	t.Run("RandomOperations", func(t *testing.T) {
		q := newQueue()
		model := make([]T, 0, 64)
		rng := rand.New(rand.NewSource(1))

		for step := 0; step < 1000; step++ {
			if rng.Intn(2) == 0 {
				v := values[rng.Intn(len(values))]
				q.Enqueue(v)
				model = append(model, v)
			} else {
				item, ok := q.Dequeue()
				if len(model) == 0 {
					if ok || item != zero {
						t.Fatalf("step %d: Dequeue() on an empty queue = (%v, %t), want (%v, false)", step, item, ok, zero)
					}
				} else {
					want := model[0]
					model = model[1:]
					if !ok || item != want {
						t.Fatalf("step %d: Dequeue() = (%v, %t), want (%v, true)", step, item, ok, want)
					}
				}
			}

			if got := q.Len(); got != len(model) {
				t.Fatalf("step %d: Len() = %d, want %d", step, got, len(model))
			}
			item, ok := q.Peek()
			if len(model) == 0 {
				if ok || item != zero {
					t.Fatalf("step %d: Peek() on an empty queue = (%v, %t), want (%v, false)", step, item, ok, zero)
				}
				continue
			}
			if want := model[0]; !ok || item != want {
				t.Fatalf("step %d: Peek() = (%v, %t), want (%v, true)", step, item, ok, want)
			}
		}
	})

	t.Run("IndependentInstances", func(t *testing.T) {
		a := newQueue()
		b := newQueue()
		a.Enqueue(values[0])

		if got := b.Len(); got != 0 {
			t.Errorf("Len() of a second queue = %d, want 0", got)
		}
		b.Enqueue(values[1])
		if item, ok := a.Peek(); !ok || item != values[0] {
			t.Errorf("Peek() of the first queue = (%v, %t), want (%v, true)", item, ok, values[0])
		}
	})
}

func TestSlicedQueueInt(t *testing.T) {
	RunQueueTests(t, func() Queue[int] { return &SlicedQueue[int]{} }, []int{1, 2, 3, 42, -7})
}

func TestSlicedQueueString(t *testing.T) {
	RunQueueTests(t, func() Queue[string] { return &SlicedQueue[string]{} }, []string{"a", "b", "c", ""})
}

func TestLinkedQueueInt(t *testing.T) {
	RunQueueTests(t, func() Queue[int] { return &LinkedQueue[int]{} }, []int{1, 2, 3, 42, -7})
}

func TestLinkedQueueString(t *testing.T) {
	RunQueueTests(t, func() Queue[string] { return &LinkedQueue[string]{} }, []string{"a", "b", "c", ""})
}
