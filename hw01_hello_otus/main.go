package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	outputString := "Hello, OTUS!"

	fmt.Println(reverse.String(outputString))
}
