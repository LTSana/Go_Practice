package main

import "fmt"

var a int

var (
	b bool
	c float32
	d string
)

func main() {
	a = 24
	b = true
	c, d = 32.6, "Hello"

	fmt.Println(a, b, c, d)
}
