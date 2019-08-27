package main

import "reflect"

type LocalVar struct {
	next   *LocalVar
	name   []rune // Variable name
	offset int    // Offset from RBP
}

var localVars *LocalVar

func newLocalVar(name []rune) *LocalVar {
	prevOffset := 0
	if localVars != nil {
		prevOffset = localVars.offset
	}

	lvar := &LocalVar{
		next:   localVars,
		name:   name,
		offset: prevOffset + 8,
	}
	localVars = lvar
	return lvar
}

func findLocalVar(name []rune) *LocalVar {
	for lvar := localVars; lvar != nil; lvar = lvar.next {
		if reflect.DeepEqual(lvar.name, name) {
			return lvar
		}
	}
	return nil
}

func findOrCreateLocalVar(name []rune) *LocalVar {
	lvar := findLocalVar(name)
	if lvar == nil {
		lvar = newLocalVar(name)
	}
	return lvar
}

func getLocalVarsOffset() int {
	if localVars != nil {
		return localVars.offset
	}
	return 0
}

type NodeKind int

const (
	ndEq     = iota // ==
	ndNe            // !=
	ndLt            // <
	ndLe            // <=
	ndAdd           // +
	ndSub           // -
	ndMul           // *
	ndDiv           // /
	ndAssign        // =
	ndLvar          // Local variable
	ndNum           // Integer
)

type Node struct {
	kind   NodeKind
	lhs    *Node // Left-hand side
	rhs    *Node // Right-hand side
	val    int   // Valid only if kind is ndNum
	offset int   // Valid only if kind is ndLvar
}

func newNode(kind NodeKind, lhs *Node, rhs *Node) *Node {
	return &Node{
		kind: kind,
		lhs:  lhs,
		rhs:  rhs,
	}
}

func newNodeIdent(name []rune) *Node {
	lvar := findOrCreateLocalVar(name)
	return &Node{
		kind:   ndLvar,
		offset: lvar.offset,
	}
}

func newNodeNum(val int) *Node {
	return &Node{
		kind: ndNum,
		val:  val,
	}
}

func program() []*Node {
	var code []*Node
	for !atEOF() {
		code = append(code, stmt())
	}
	return code
}

func stmt() *Node {
	node := expr()
	expect(";")
	return node
}

func expr() *Node {
	return assign()
}

func assign() *Node {
	node := equality()

	if consume("=") {
		node = newNode(ndAssign, node, assign())
	}
	return node
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

	token := consumeIdent()
	if token != nil {
		return newNodeIdent(token.str)
	}

	return newNodeNum(expectNumber())
}
