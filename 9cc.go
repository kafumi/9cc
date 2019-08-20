package main

import "fmt"
import "os"
import "strconv"
import "strings"
import "reflect"
import "unicode"

var userInput []rune

func fatal(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func fatalAt(pos int, format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s\n", string(userInput))
	fmt.Fprintf(os.Stderr, strings.Repeat(" ", pos))
	fmt.Fprintf(os.Stderr, "^ ")
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func fatalAtStr(loc []rune, format string, a ...interface{}) {
	fatalAt(len(userInput)-len(loc), format, a...)
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
	pos  int
	str  []rune
}

var token *Token

func consume(op string) bool {
	if token.kind == tkReserved && reflect.DeepEqual(token.str, []rune(op)) {
		token = token.next
		return true
	}
	return false
}

func expect(op string) {
	if token.kind == tkReserved && reflect.DeepEqual(token.str, []rune(op)) {
		token = token.next
	} else {
		fatalAt(token.pos, "Next character is not '%s'", op)
	}
}

func expectNumber() int {
	if token.kind != tkNum {
		fatalAt(token.pos, "Next token is not number")
	}
	val := token.val
	token = token.next
	return val
}

func atEOF() bool {
	return token.kind == tkEOF
}

func newToken(kind TokenKind, cur *Token, str []rune, tokenLen int) *Token {
	tok := &Token{
		kind: kind,
		pos:  len(userInput) - len(str),
		str:  str[:tokenLen],
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
		case '=', '!':
			if p[1] != '=' {
				fatalAtStr(p[1:], "Next character is not '='")
			}
			cur = newToken(tkReserved, cur, p, 2)
			p = p[2:]
			continue
		case '<', '>':
			if p[1] == '=' {
				cur = newToken(tkReserved, cur, p, 2)
				p = p[2:]
			} else {
				cur = newToken(tkReserved, cur, p, 1)
				p = p[1:]
			}
			continue
		case '+', '-', '*', '/', '(', ')':
			cur = newToken(tkReserved, cur, p, 1)
			p = p[1:]
			continue
		}

		if unicode.IsDigit(p[0]) {
			val, newP := readNumber(p)
			cur = newToken(tkNum, cur, p, len(p)-len(newP))
			cur.val = val
			p = newP
			continue
		}

		fatalAtStr(p, "Unable to tokenize")
	}

	newToken(tkEOF, cur, p, 0)
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
		fatalAtStr(program, "Expect number")
	}

	return number, remaining
}

type NodeKind int

const (
	ndEq = iota
	ndNe
	ndLt
	ndLe
	ndAdd
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
	return equality()
}

func equality() *Node {
	node := relational()

	for {
		if consume("==") {
			node = newNode(ndEq, node, relational())
		} else if consume("!=") {
			node = newNode(ndNe, node, relational())
		} else {
			return node
		}
	}
}

func relational() *Node {
	node := add()

	for {
		if consume("<") {
			node = newNode(ndLt, node, add())
		} else if consume("<=") {
			node = newNode(ndLe, node, add())
		} else if consume(">") {
			node = newNode(ndLt, add(), node)
		} else if consume(">=") {
			node = newNode(ndLe, add(), node)
		} else {
			return node
		}
	}
}

func add() *Node {
	node := mul()

	for {
		if consume("+") {
			node = newNode(ndAdd, node, mul())
		} else if consume("-") {
			node = newNode(ndSub, node, mul())
		} else {
			return node
		}
	}
}

func mul() *Node {
	node := unary()

	for {
		if consume("*") {
			node = newNode(ndMul, node, unary())
		} else if consume("/") {
			node = newNode(ndDiv, node, unary())
		} else {
			return node
		}
	}
}

func unary() *Node {
	if consume("+") {
		return primary()
	}
	if consume("-") {
		return newNode(ndSub, newNodeNum(0), primary())
	}
	return primary()
}

func primary() *Node {
	if consume("(") {
		node := expr()
		expect(")")
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
	case ndEq, ndNe, ndLt, ndLe:
		genCmp(node)
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

func genCmp(node *Node) {
	fmt.Printf("  cmp rax, rdi\n")
	switch node.kind {
	case ndEq:
		fmt.Printf("  sete al\n")
	case ndNe:
		fmt.Printf("  setne al\n")
	case ndLt:
		fmt.Printf("  setl al\n")
	case ndLe:
		fmt.Printf("  setle al\n")
	}
	fmt.Printf("  movzb rax, al\n")
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
