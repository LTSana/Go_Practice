package main

import (
	"fmt"
)

func main() {

	c := make(chan int)
	for i := 0; i < 5; i++ {
		go cookingGopher(i, c)
	}

	for i := 0; i < 5; i++ {
		gopherID := <-c
		fmt.Printf("Gopher %d finished cooking\n", gopherID)
	}
}

func cookingGopher(id int, c chan int) {
	fmt.Printf("Gopher %d started cooking\n", id)
	c <- id
}
