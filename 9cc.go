package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("usage: 9cc <program>\n")
		os.Exit(1)
	}

	userInput = []rune(os.Args[1])
	token = tokenize(userInput)
	node := expr()

	fmt.Printf(".intel_syntax noprefix\n")
	fmt.Printf(".global main\n")
	fmt.Printf("main:\n")

	gen(node)

	fmt.Printf("  pop rax\n")
	fmt.Printf("  ret\n")
}
