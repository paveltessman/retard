package queue

type SlicedQueue[T any] struct {
	items []T
}

func (q *SlicedQueue[T]) Enqueue(item T) {
	q.items = append(q.items, item)
}

func (q *SlicedQueue[T]) Dequeue() (T, bool) {
	var zero T

	if q.Len() == 0 {
		return zero, false
	}

	item := q.items[0]
	q.items[0] = zero
	q.items = q.items[1:]

	return item, true
}

func (q *SlicedQueue[T]) Peek() (T, bool) {
	var zero T

	if q.Len() == 0 {
		return zero, false
	}

	return q.items[0], true

}

func (q *SlicedQueue[T]) Len() int {
	return len(q.items)
}
