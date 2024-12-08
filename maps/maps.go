package main

import "fmt"

func main() {

	firstReleases := map[string]int{
		"C": 1972, "C++": 1985, "Java": 1996,
		"Python": 1991, "JavaScript": 1996, "Go": 2012,
	}

	for index, val := range firstReleases {
		fmt.Printf("%s was first released on %d\n", index, val)
	}
}
