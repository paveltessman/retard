package stack

import (
	"math/rand"
	"testing"
)

// RunStackTests checks one Stack implementation against the interface contract.
// newStack must give back an empty stack on every call.
// values must hold at least 3 distinct items to push.
func RunStackTests[T comparable](t *testing.T, newStack func() Stack[T], values []T) {
	t.Helper()

	if len(values) < 3 {
		t.Fatalf("RunStackTests needs at least 3 values, got %d", len(values))
	}
	var zero T

	t.Run("EmptyStack", func(t *testing.T) {
		s := newStack()

		if got := s.Len(); got != 0 {
			t.Errorf("Len() on a new stack = %d, want 0", got)
		}
		if item, ok := s.Peek(); ok || item != zero {
			t.Errorf("Peek() on a new stack = (%v, %t), want (%v, false)", item, ok, zero)
		}
		if item, ok := s.Pop(); ok || item != zero {
			t.Errorf("Pop() on a new stack = (%v, %t), want (%v, false)", item, ok, zero)
		}
		if got := s.Len(); got != 0 {
			t.Errorf("Len() after a failed Pop = %d, want 0", got)
		}
	})

	t.Run("PushThenPeek", func(t *testing.T) {
		s := newStack()
		s.Push(values[0])

		if got := s.Len(); got != 1 {
			t.Errorf("Len() after one Push = %d, want 1", got)
		}
		// Peek must not remove the item, so two calls give the same answer.
		for call := 1; call <= 2; call++ {
			item, ok := s.Peek()
			if !ok || item != values[0] {
				t.Errorf("Peek() call %d = (%v, %t), want (%v, true)", call, item, ok, values[0])
			}
			if got := s.Len(); got != 1 {
				t.Errorf("Len() after Peek call %d = %d, want 1", call, got)
			}
		}
	})

	t.Run("LastInFirstOut", func(t *testing.T) {
		s := newStack()
		for i, v := range values {
			s.Push(v)
			if got := s.Len(); got != i+1 {
				t.Fatalf("Len() after %d pushes = %d, want %d", i+1, got, i+1)
			}
		}

		for i := len(values) - 1; i >= 0; i-- {
			want := values[i]
			if item, ok := s.Peek(); !ok || item != want {
				t.Fatalf("Peek() at depth %d = (%v, %t), want (%v, true)", i, item, ok, want)
			}
			item, ok := s.Pop()
			if !ok || item != want {
				t.Fatalf("Pop() at depth %d = (%v, %t), want (%v, true)", i, item, ok, want)
			}
			if got := s.Len(); got != i {
				t.Fatalf("Len() after Pop at depth %d = %d, want %d", i, got, i)
			}
		}
	})

	t.Run("ReuseAfterDrain", func(t *testing.T) {
		s := newStack()
		for _, v := range values {
			s.Push(v)
		}
		for range values {
			if _, ok := s.Pop(); !ok {
				t.Fatal("Pop() on a full stack returned false")
			}
		}

		if item, ok := s.Pop(); ok || item != zero {
			t.Errorf("Pop() on a drained stack = (%v, %t), want (%v, false)", item, ok, zero)
		}
		// The stack must work again after it runs empty.
		s.Push(values[1])
		if item, ok := s.Pop(); !ok || item != values[1] {
			t.Errorf("Pop() after a reuse Push = (%v, %t), want (%v, true)", item, ok, values[1])
		}
		if got := s.Len(); got != 0 {
			t.Errorf("Len() after the reuse Pop = %d, want 0", got)
		}
	})

	t.Run("DuplicateItems", func(t *testing.T) {
		s := newStack()
		const count = 4
		for i := 0; i < count; i++ {
			s.Push(values[0])
		}
		if got := s.Len(); got != count {
			t.Errorf("Len() after %d equal pushes = %d, want %d", count, got, count)
		}
		for i := 0; i < count; i++ {
			if item, ok := s.Pop(); !ok || item != values[0] {
				t.Fatalf("Pop() number %d = (%v, %t), want (%v, true)", i+1, item, ok, values[0])
			}
		}
		if got := s.Len(); got != 0 {
			t.Errorf("Len() after all pops = %d, want 0", got)
		}
	})

	t.Run("RandomOperations", func(t *testing.T) {
		s := newStack()
		model := make([]T, 0, 64)
		rng := rand.New(rand.NewSource(1))

		for step := 0; step < 1000; step++ {
			if rng.Intn(2) == 0 {
				v := values[rng.Intn(len(values))]
				s.Push(v)
				model = append(model, v)
			} else {
				item, ok := s.Pop()
				if len(model) == 0 {
					if ok || item != zero {
						t.Fatalf("step %d: Pop() on an empty stack = (%v, %t), want (%v, false)", step, item, ok, zero)
					}
				} else {
					want := model[len(model)-1]
					model = model[:len(model)-1]
					if !ok || item != want {
						t.Fatalf("step %d: Pop() = (%v, %t), want (%v, true)", step, item, ok, want)
					}
				}
			}

			if got := s.Len(); got != len(model) {
				t.Fatalf("step %d: Len() = %d, want %d", step, got, len(model))
			}
			item, ok := s.Peek()
			if len(model) == 0 {
				if ok || item != zero {
					t.Fatalf("step %d: Peek() on an empty stack = (%v, %t), want (%v, false)", step, item, ok, zero)
				}
				continue
			}
			if want := model[len(model)-1]; !ok || item != want {
				t.Fatalf("step %d: Peek() = (%v, %t), want (%v, true)", step, item, ok, want)
			}
		}
	})

	t.Run("IndependentInstances", func(t *testing.T) {
		a := newStack()
		b := newStack()
		a.Push(values[0])

		if got := b.Len(); got != 0 {
			t.Errorf("Len() of a second stack = %d, want 0", got)
		}
		b.Push(values[1])
		if item, ok := a.Peek(); !ok || item != values[0] {
			t.Errorf("Peek() of the first stack = (%v, %t), want (%v, true)", item, ok, values[0])
		}
	})
}

func TestSlicedStackInt(t *testing.T) {
	RunStackTests(t, func() Stack[int] { return &SlicedStack[int]{} }, []int{1, 2, 3, 42, -7})
}

func TestSlicedStackString(t *testing.T) {
	RunStackTests(t, func() Stack[string] { return &SlicedStack[string]{} }, []string{"a", "b", "c", ""})
}

func TestLinkedStackInt(t *testing.T) {
	RunStackTests(t, func() Stack[int] { return &LinkedStack[int]{} }, []int{1, 2, 3, 42, -7})
}

func TestLinkedStackString(t *testing.T) {
	RunStackTests(t, func() Stack[string] { return &LinkedStack[string]{} }, []string{"a", "b", "c", ""})
}
