package deque

type RingDeque[T any] struct {
	buf  []T
	head int
	size int
}

func NewRingDeque[T any]() *RingDeque[T] {
	return &RingDeque[T]{make([]T, 1), 0, 0}
}

func (d *RingDeque[T]) PushFront(item T) {
	d.growIfFull()
	d.head = (d.head - 1 + len(d.buf)) % len(d.buf)
	d.buf[d.head] = item
	d.size++
}

func (d *RingDeque[T]) PushBack(item T) {
	d.growIfFull()
	d.buf[d.at(d.size)] = item
	d.size++
}

func (d *RingDeque[T]) PopFront() (T, bool) {
	var zero T
	if d.size == 0 {
		return zero, false
	}

	item := d.buf[d.head]
	d.buf[d.head] = zero
	d.head = (d.head + 1) % len(d.buf)
	d.size--
	return item, true
}

func (d *RingDeque[T]) PopBack() (T, bool) {
	var zero T
	if d.size == 0 {
		return zero, false
	}

	idx := d.at(d.size - 1)
	item := d.buf[idx]
	d.buf[idx] = zero
	d.size--
	return item, true
}

func (d *RingDeque[T]) Len() int {
	return d.size
}

func (d *RingDeque[T]) at(i int) int {
	return (d.head + i) % len(d.buf)
}

func (d *RingDeque[T]) growIfFull() {
	if d.size < len(d.buf) {
		return
	}

	newBuf := make([]T, 2*len(d.buf))

	for i := range len(d.buf) {
		newBuf[i] = d.buf[d.at(i)]
	}
	d.buf = newBuf
	d.head = 0
}
