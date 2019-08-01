package main

import "fmt"
import "log"
import "os"
import "strconv"
import "unicode"

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("usage: 9cc <program>\n")
		os.Exit(1)
	}

	program := []rune(os.Args[1])

	fmt.Printf(".intel_syntax noprefix\n")
	fmt.Printf(".global main\n")
	fmt.Printf("main:\n")

	value, program := readNumber(program)
	fmt.Printf("  mov rax, %d\n", value)

	for len(program) > 0 {
		if program[0] == '+' {
			program = program[1:]
			value, program = readNumber(program)
			fmt.Printf("  add rax, %d\n", value)
			continue
		}

		if program[0] == '-' {
			program = program[1:]
			value, program = readNumber(program)
			fmt.Printf("  sub rax, %d\n", value)
			continue
		}

		fmt.Fprintf(os.Stderr, "Unexpected character: '%c'\n", program[0])
		os.Exit(1)
	}

	fmt.Printf("  ret\n")
}

func readNumber(program []rune) (int, []rune) {
	length := 0
	for length < len(program) && unicode.IsDigit(program[length]) {
		length++
	}

	target := string(program[0:length])
	remaining := program[length:]
	number, err := strconv.Atoi(target)
	if err != nil {
		log.Fatal(err)
	}

	return number, remaining
}
