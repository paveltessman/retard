package deque

type Deque[T any] interface {
	PushFront(T)
	PushBack(T)
	PopFront() (T, bool)
	PopBack() (T, bool)
	Len() int
}
