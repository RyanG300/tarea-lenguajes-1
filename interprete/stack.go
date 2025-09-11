package main

/*
------------------------------------------------------------------------
--------------------------------Pila------------------------------------
------------------------------------------------------------------------
*/

// La pila, tiene los metodos push, pop, top, isEmpty
type Stack struct {
	items []any
}

var stack Stack

/*------------------------------------------------------------------------
---------------------------Atributos de la pila---------------------------
------------------------------------------------------------------------*/

func (s *Stack) push(v any) {
	s.items = append(s.items, v)
}

func (s *Stack) pop() bool {
	if len(s.items) == 0 {
		return false
	}
	i := len(s.items) - 1
	//v := s.items[i]
	s.items = s.items[:i]
	return true
}

func (s *Stack) top() (any, bool) {
	var x any
	if len(s.items) == 0 {
		return x, false
	}
	x = s.items[len(s.items)-1]
	return x, true
}

func (s *Stack) isEmpty() bool {
	return len(s.items) == 0
}
