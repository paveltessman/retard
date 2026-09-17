package ll

type node[T any] struct {
	next *node[T]
	val  T
}
