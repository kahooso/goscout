// task-04 (A0.3): изменяемый стек на pointer receiver.
// Разбор: knowledge-base/topics/go/pointers.md
package main

import "fmt"

type Stack struct {
	items []int
}

func NewStack() *Stack {
	return &Stack{items: make([]int, 0)}
}

func (s *Stack) Push(v int) {
	s.items = append(s.items, v)
}

func (s *Stack) Pop() (int, bool) {
	if len(s.items) == 0 {
		return 0, false
	}
	last := s.items[len(s.items)-1] // сохранить ДО обрезки
	s.items = s.items[:len(s.items)-1]
	return last, true
}

func (s *Stack) Peek() (int, bool) {
	if len(s.items) == 0 {
		return 0, false
	}
	return s.items[len(s.items)-1], true
}

func main() {
	s := NewStack()
	s.Push(1)
	s.Push(4)
	s.Push(5)
	s.Push(0)

	fmt.Println(s.Peek()) // 0 true
	fmt.Println(s.Pop())  // 0 true
	fmt.Println(s.Peek()) // 5 true
}
