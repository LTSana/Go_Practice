package main

import "fmt"

func average(x []float64) (avg float64) {
	total := 0.00
	if len(x) <= 0 {
		avg = 0
	} else {
		for _, val := range x {
			total += val
		}
		avg = total / float64(len(x))
	}
	return
}

func main() {
	x := []float64{2.15, 3.14, 42.0, 29.5}
	fmt.Printf("Average: %f\n", average(x))
}
