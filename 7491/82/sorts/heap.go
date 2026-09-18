package sorts

import "slices"

type maxHeap struct {
	data []int
}

func newMaxHeap(data []int) *maxHeap {
	h := maxHeap{data}
	n := len(data)

	for i := n/2 - 1; i >= 0; i-- {
		h.siftDown(i, n)
	}
	return &h
}

func (h *maxHeap) siftUp(i int) {
	for i > 0 {
		parent := (i - 1) / 2

		if h.data[parent] > h.data[i] {
			break
		}

		h.data[parent], h.data[i] = h.data[i], h.data[parent]
		i = parent
	}
}

func (h *maxHeap) siftDown(parent int, n int) {
	for parent < n {
		largest := parent
		left, right := 2*parent+1, 2*parent+2

		if left < n && h.data[left] > h.data[largest] {
			largest = left
		}

		if right < n && h.data[right] > h.data[largest] {
			largest = right
		}

		if largest == parent {
			break
		}

		h.data[parent], h.data[largest] = h.data[largest], h.data[parent]
		parent = largest
	}
}

func (h *maxHeap) pop() (int, bool) {
	n := len(h.data)
	if n == 0 {
		return 0, false
	}

	max := h.data[0]

	h.data[0], h.data[n-1] = h.data[n-1], h.data[0]
	h.data = h.data[:n-1]

	if len(h.data) > 0 {
		h.siftDown(0, len(h.data))
	}
	return max, true
}

func (h *maxHeap) push(val int) {
	h.data = append(h.data, val)
	h.siftUp(len(h.data))
}

func heapSort(nums []int) []int {
	n := len(nums)

	h := newMaxHeap(slices.Clone(nums))

	for end := n - 1; end > 0; end-- {
		h.data[0], h.data[end] = h.data[end], h.data[0]
		h.siftDown(0, end)
	}
	return h.data
}
