package main

import "fmt"
import "log"
import "os"
import "strconv"
import "strings"
import "unicode"

type TokenKind int

const (
	tkReserved TokenKind = iota
	tkNum
	tkEOF
)

type Token struct {
	kind TokenKind
	next *Token
	val  int
	str  []rune
}

var userInput []rune
var token *Token

func fatal(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func fatalAt(loc []rune, format string, a ...interface{}) {
	pos := len(userInput) - len(loc)
	fmt.Fprintf(os.Stderr, "%s\n", string(userInput))
	fmt.Fprintf(os.Stderr, strings.Repeat(" ", pos))
	fmt.Fprintf(os.Stderr, "^ ")
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func consume(op rune) bool {
	if token.kind == tkReserved && token.str[0] == op {
		token = token.next
		return true
	}
	return false
}

func expect(op rune) {
	if token.kind == tkReserved && token.str[0] == op {
		token = token.next
	} else {
		fatalAt(token.str, "Next character is not '%c'", op)
	}
}

func expectNumber() int {
	if token.kind != tkNum {
		fatalAt(token.str, "Next token is not number")
	}
	val := token.val
	token = token.next
	return val
}

func atEOF() bool {
	return token.kind == tkEOF
}

func newToken(kind TokenKind, cur *Token, str []rune) *Token {
	tok := &Token{
		kind: kind,
		str:  str,
	}
	cur.next = tok
	return tok
}

func tokenize(p []rune) *Token {
	head := Token{
		next: nil,
	}
	cur := &head

	for len(p) > 0 {
		if unicode.IsSpace(p[0]) {
			p = p[1:]
			continue
		}

		if p[0] == '+' || p[0] == '-' {
			cur = newToken(tkReserved, cur, p)
			p = p[1:]
			continue
		}

		if unicode.IsDigit(p[0]) {
			cur = newToken(tkNum, cur, p)
			cur.val, p = readNumber(p)
			continue
		}

		fatalAt(p, "Unable to tokenize")
	}

	newToken(tkEOF, cur, p)
	return head.next
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("usage: 9cc <program>\n")
		os.Exit(1)
	}

	userInput = []rune(os.Args[1])
	token = tokenize(userInput)

	fmt.Printf(".intel_syntax noprefix\n")
	fmt.Printf(".global main\n")
	fmt.Printf("main:\n")

	fmt.Printf("  mov rax, %d\n", expectNumber())

	for !atEOF() {
		if consume('+') {
			fmt.Printf("  add rax, %d\n", expectNumber())
			continue
		}

		if consume('-') {
			fmt.Printf("  sub rax, %d\n", expectNumber())
			continue
		}

		fatalAt(token.str, "Unexpected token")
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
