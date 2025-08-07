package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse" //nolint
)

func main() {
	fmt.Println(reverse.String("Hello, OTUS!"))
}
