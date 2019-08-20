package main

import "fmt"
import "log"
import "os"
import "strconv"
import "strings"
import "unicode"

var userInput []rune

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

var token *Token

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

		switch p[0] {
		case '+', '-', '*', '/', '(', ')':
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

type NodeKind int

const (
	ndAdd = iota
	ndSub
	ndMul
	ndDiv
	ndNum
)

type Node struct {
	kind NodeKind
	lhs  *Node
	rhs  *Node
	val  int
}

func newNode(kind NodeKind, lhs *Node, rhs *Node) *Node {
	return &Node{
		kind: kind,
		lhs:  lhs,
		rhs:  rhs,
	}
}

func newNodeNum(val int) *Node {
	return &Node{
		kind: ndNum,
		val:  val,
	}
}

func expr() *Node {
	node := mul()

	for {
		if consume('+') {
			node = newNode(ndAdd, node, mul())
		} else if consume('-') {
			node = newNode(ndSub, node, mul())
		} else {
			return node
		}
	}
}

func mul() *Node {
	node := primary()

	for {
		if consume('*') {
			node = newNode(ndMul, node, primary())
		} else if consume('/') {
			node = newNode(ndDiv, node, primary())
		} else {
			return node
		}
	}
}

func primary() *Node {
	if consume('(') {
		node := expr()
		expect(')')
		return node
	}

	return newNodeNum(expectNumber())
}

func gen(node *Node) {
	if node.kind == ndNum {
		fmt.Printf("  push %d\n", node.val)
		return
	}

	gen(node.lhs)
	gen(node.rhs)

	fmt.Printf("  pop rdi\n")
	fmt.Printf("  pop rax\n")

	switch node.kind {
	case ndAdd:
		fmt.Printf("  add rax, rdi\n")
	case ndSub:
		fmt.Printf("  sub rax, rdi\n")
	case ndMul:
		fmt.Printf("  imul rax, rdi\n")
	case ndDiv:
		fmt.Printf("  cqo\n")
		fmt.Printf("  idiv rdi\n")
	}

	fmt.Printf("  push rax\n")
}

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
