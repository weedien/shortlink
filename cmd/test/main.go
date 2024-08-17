package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

func main() {
	// Create a new Person struct
	p := Person{
		Name: "John",
	}
	// Print the struct
	fmt.Println(p)
}
