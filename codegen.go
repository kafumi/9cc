package main

import "fmt"

var labelSeq = 0

func gen(node *Node) {
	switch node.kind {
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
	case ndFor:
		seq := labelSeq
		labelSeq++

		if node.init != nil {
			gen(node.init)
			genPop()
		}
		fmt.Printf(".L%s%d:\n", "begin", seq)
		if node.test != nil {
			gen(node.test)
			fmt.Printf("  pop rax\n")
			fmt.Printf("  cmp rax, 0\n")
			fmt.Printf("  je  .L%s%d\n", "end", seq)
		}
		gen(node.cons)
		genPop()
		if node.post != nil {
			gen(node.post)
			genPop()
		}
		fmt.Printf("  jmp .L%s%d\n", "begin", seq)
		fmt.Printf(".L%s%d:\n", "end", seq)
		genPush()
		return
	case ndBlock:
		for _, stmt := range node.body {
			gen(stmt)
			genPop()
		}
		genPush()
		return
	case ndReturn:
		gen(node.lhs)
		fmt.Printf("  pop rax\n")
		genEpilogue()
		return
	case ndFcall:
		seq := labelSeq
		labelSeq++

		argRegs := []string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"}
		for _, arg := range node.args {
			gen(arg)
		}
		for i := len(node.args) - 1; i >= 0; i-- {
			fmt.Printf("  pop %s\n", argRegs[i])
		}

		// We need to make RSP 16 byte aligned when calling function.
		fmt.Printf("  mov rax, rsp\n")
		fmt.Printf("  and rax, 15\n")
		fmt.Printf("  cmp rax, 0\n")
		fmt.Printf("  je  .L%s%d\n", "call", seq)
		fmt.Printf("  sub rsp, 8\n")
		fmt.Printf("  call %s\n", node.funcName)
		fmt.Printf("  add rsp, 8\n")
		fmt.Printf("  jmp .L%s%d\n", "end", seq)
		fmt.Printf(".L%s%d:\n", "call", seq)
		fmt.Printf("  call %s\n", node.funcName)
		fmt.Printf(".L%s%d:\n", "end", seq)
		fmt.Printf("  push rax\n")
		return
	case ndLvar:
		genLval(node)
		fmt.Printf("  pop rax\n")
		fmt.Printf("  mov rax, [rax]\n")
		fmt.Printf("  push rax\n")
		return
	case ndNum:
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
	fmt.Printf("  sub rsp, %d\n", env.maxOffset)
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
