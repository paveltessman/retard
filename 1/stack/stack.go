package stack

type Stack[T any] interface {
	Push(item T)
	Pop() (T, bool)
	Peek() (T, bool)
	Len() int
}
