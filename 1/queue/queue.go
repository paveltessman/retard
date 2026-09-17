package queue

type Queue[T any] interface {
	Enqueue(T)
	Dequeue() (T, bool)
	Peek() (T, bool)
	Len() int
}
