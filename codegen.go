package main

import "fmt"

func gen(node *Node) {
	switch node.kind {
	case ndNum:
		fmt.Printf("  push %d\n", node.val)
		return
	case ndLvar:
		genLval(node)
		fmt.Printf("  pop rax\n")
		fmt.Printf("  mov rax, [rax]\n")
		fmt.Printf("  push rax\n")
		return
	case ndAssign:
		genLval(node.lhs)
		gen(node.rhs)
		fmt.Printf("  pop rdi\n")
		fmt.Printf("  pop rax\n")
		fmt.Printf("  mov [rax], rdi\n")
		fmt.Printf("  push rdi\n")
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

func genLval(node *Node) {
	if node.kind != ndLvar {
		fatal("Left-hand side of assign expression is not variable")
	}

	fmt.Printf("  mov rax, rbp\n")
	fmt.Printf("  sub rax, %d\n", node.offset)
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

func genProgramHeader() {
	fmt.Printf(".intel_syntax noprefix\n")
	fmt.Printf(".global main\n")
	fmt.Printf("main:\n")
}

func genPrologue() {
	fmt.Printf("  push rbp\n")
	fmt.Printf("  mov rbp, rsp\n")
	fmt.Printf("  sub rsp, 208\n") // 208 = 8 bits * 26 variables
}

func genEpilogue() {
	fmt.Printf("  mov rsp, rbp\n")
	fmt.Printf("  pop rbp\n")
	fmt.Printf("  ret\n")
}

func genPop() {
	fmt.Printf("  pop rax\n")
}
