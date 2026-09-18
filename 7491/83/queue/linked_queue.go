package queue

type node[T any] struct {
	next *node[T]
	item T
}

type LinkedQueue[T any] struct {
	first, last *node[T]
	len         int
}

func (q *LinkedQueue[T]) Enqueue(item T) {
	new := &node[T]{next: nil, item: item}
	q.len++
	if q.first == nil {
		q.first = new
		q.last = q.first
		return
	}
	q.last.next = new
	q.last = new
}

func (q *LinkedQueue[T]) Dequeue() (T, bool) {
	first := q.first
	if first == nil {
		var zero T
		return zero, false
	}

	q.first = q.first.next
	q.len--
	return first.item, true
}

func (q *LinkedQueue[T]) Peek() (T, bool) {
	if q.first == nil {
		var zero T
		return zero, false
	}
	return q.first.item, true
}

func (q *LinkedQueue[T]) Len() int {
	if q.len < 0 {
		panic("fuk")
	}
	return q.len
}
