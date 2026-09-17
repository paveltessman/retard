package stack

var _ Stack[int] = (*LinkedStack[int])(nil)

type node[T any] struct {
	next *node[T]
	item T
}

type LinkedStack[T any] struct {
	top *node[T]
	len int
}

func (s *LinkedStack[T]) Push(item T) {
	node := &node[T]{next: s.top, item: item}
	s.top = node
	s.len++
}

func (s *LinkedStack[T]) Pop() (T, bool) {
	var zero T
	if s.Len() == 0 {
		return zero, false
	}
	top := s.top
	s.top = top.next
	s.len--
	return top.item, true
}

func (s *LinkedStack[T]) Peek() (T, bool) {
	var zero T
	if s.Len() == 0 {
		return zero, false
	}
	return s.top.item, true
}

func (s *LinkedStack[T]) Len() int {
	if s.len < 0 {
		panic("fuck")
	}
	return s.len
}
