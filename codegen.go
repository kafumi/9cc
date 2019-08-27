package main

import "fmt"

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
