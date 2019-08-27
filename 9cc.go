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
	nodes := program()

	genProgramHeader()
	genPrologue()

	for _, node := range nodes {
		gen(node)
		genPop()
	}

	genEpilogue()
}
