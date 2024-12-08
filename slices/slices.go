package main

import (
	"fmt"
	"reflect"
)

func main() {

	languages := [9]string{
		"C", "Lisp", "C++", "Java", "Python",
		"JavaScript", "Ruby", "Go", "Rust",
	}

	classics := languages[0:3]
	modern := make([]string, 4)
	modern = languages[3:7]
	new := languages[7:9]

	fmt.Printf("Classics: %v\n", classics)
	fmt.Printf("Modern: %v\n", modern)
	fmt.Printf("New: %v\n", new)

	fmt.Println(" --- ")

	allLangs := languages[:]
	fmt.Println(reflect.TypeOf(allLangs).Kind())

	frameworks := []string{
		"React", "Vue", "Angular", "Svelte",
		"Laravel", "Django", "Flask", "Fiber",
	}

	jsFrameworks := frameworks[0:4:4]
	frameworks = append(frameworks, "Meteor")

	fmt.Printf("All Frameworks: %v\n", frameworks)
	fmt.Printf("JS Frameworks: %v\n", jsFrameworks)
}
