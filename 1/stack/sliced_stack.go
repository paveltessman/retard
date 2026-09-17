package stack

type SlicedStack[T any] struct {
	items []T
}

func (s *SlicedStack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *SlicedStack[T]) Pop() (T, bool) {
	var zero T
	length := s.Len()

	if length == 0 {
		return zero, false
	}

	item := s.items[length-1]
	s.items[length-1] = zero
	s.items = s.items[:length-1]

	return item, true
}

func (s *SlicedStack[T]) Peek() (T, bool) {
	var zero T
	length := s.Len()

	if length == 0 {
		return zero, false
	}

	item := s.items[length-1]
	return item, true
}

func (s *SlicedStack[T]) Len() int {
	return len(s.items)
}
