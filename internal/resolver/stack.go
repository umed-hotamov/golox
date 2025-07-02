package resolver

type Stack struct {
	items []any
}

func NewStack() *Stack {
	return &Stack{
		items: make([]any, 0),
	}
}

func (s *Stack) Push(item any) {
	s.items = append(s.items, item)
}

func (s *Stack) Pop() any {
	if s.IsEmpty() {
		return nil
	}

	last := s.items[:len(s.items)-1]
	s.items = s.items[:len(s.items)-1]

	return last
}

func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

func (s *Stack) Peek() any {
	if s.IsEmpty() {
		return nil
	}

	return s.items[len(s.items)-1]
}

func (s *Stack) Size() int {
	return len(s.items)
}

func (s *Stack) Get(index int) any {
	if s.IsEmpty() {
		return nil
	}

	return s.items[index]
}
