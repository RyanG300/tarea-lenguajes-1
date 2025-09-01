package main

import (
	"fmt"
)

type Stack struct {
	items []any
}

func (s *Stack) Push(v any) {
	s.items = append(s.items, v)
}

func (s *Stack) Pop() (any, bool) {
	if len(s.items) == 0 {
		return nil, false
	}
	i := len(s.items) - 1
	v := s.items[i]
	s.items = s.items[:i]
	return v, true
}

func (s *Stack) Top() (any, bool) {
	var x any
	if len(s.items) == 0 {
		return x, false
	}
	x = s.items[len(s.items)-1]
	return x, true
}

func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

func main() {
	var s Stack
	if s.IsEmpty() {
		fmt.Println("Pingas")
	}
}



