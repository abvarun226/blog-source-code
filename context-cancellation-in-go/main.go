package main

import (
	"flag"
	"fmt"
)

func main() {
	example := flag.Int("example", 1, "the example to run")
	flag.Parse()

	switch *example {
	case 1:
		fmt.Printf("===== Running example 1 =====\n\n")
		exampleOne()

	case 2:
		fmt.Printf("===== Running example 2 =====\n\n")
		exampleTwo()
	}
}
