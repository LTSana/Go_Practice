package main

import "fmt"

func main() {
	var address *int
	number := 42
	address = &number
	fmt.Println("Address of number variable is:", address)
	fmt.Println("Value of the number variable: ", *address)

	changeNum(address)
	fmt.Println("Changed value: ", number)
}

func changeNum(num *int) {
	*num = 100
}
