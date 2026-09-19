package minstack

type MinStack struct {
	data []int
	mins []int
}

func NewMinStack() *MinStack {
	return &MinStack{}
}

func (s *MinStack) Push(value int) {
	s.data = append(s.data, value)

	if len(s.mins) == 0 {
		s.mins = append(s.mins, value)
		return
	}
	s.mins = append(s.mins, min(s.runningMin(), value))
}

func (s *MinStack) Pop() {
	s.data = s.data[:len(s.data)-1]
	s.mins = s.mins[:len(s.mins)-1]
}

func (s *MinStack) Top() int {
	return s.data[len(s.data)-1]
}

func (s *MinStack) GetMin() int {
	return s.runningMin()
}

func (s *MinStack) runningMin() int {
	return s.mins[len(s.mins)-1]
}
