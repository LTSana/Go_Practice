package main

import "fmt"

type stack struct {
	index int
	data  [5]int
}

func (s *stack) push(k int) {
	s.data[s.index] = k
	s.index++
}

func (s *stack) pop() int {
	s.index--
	val := s.data[s.index]
	s.data[s.index] = 0
	return val
}

func main() {
	s := new(stack)
	s.push(42)
	s.push(59)
	fmt.Printf("stack: %v\n", *s)
	fmt.Println(s.pop())
	fmt.Printf("stack 2: %v\n", *s)
}
