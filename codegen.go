package main

import "fmt"

var labelSeq = 0

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
	case ndIf:
		seq := labelSeq
		labelSeq++

		gen(node.test)
		fmt.Printf("  pop rax\n")
		fmt.Printf("  cmp rax, 0\n")
		fmt.Printf("  je  .L%s%d\n", "else", seq)
		gen(node.cons)
		fmt.Printf("  jmp .L%s%d\n", "end", seq)
		fmt.Printf(".L%s%d:\n", "else", seq)
		if node.alt != nil {
			gen(node.alt)
		} else {
			genPush()
		}
		fmt.Printf(".L%s%d:\n", "end", seq)
		return
	case ndWhile:
		seq := labelSeq
		labelSeq++

		fmt.Printf(".L%s%d:\n", "begin", seq)
		gen(node.test)
		fmt.Printf("  pop rax\n")
		fmt.Printf("  cmp rax, 0\n")
		fmt.Printf("  je  .L%s%d\n", "end", seq)
		gen(node.cons)
		genPop()
		fmt.Printf("  jmp .L%s%d\n", "begin", seq)
		fmt.Printf(".L%s%d:\n", "end", seq)
		genPush()
		return
	case ndReturn:
		gen(node.lhs)
		fmt.Printf("  pop rax\n")
		genEpilogue()
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
	fmt.Printf("  sub rsp, %d\n", getLocalVarsOffset())
}

func genEpilogue() {
	fmt.Printf("  mov rsp, rbp\n")
	fmt.Printf("  pop rbp\n")
	fmt.Printf("  ret\n")
}

func genPush() {
	fmt.Printf("  push 0xdb\n")
}

func genPop() {
	fmt.Printf("  pop rax\n")
}
