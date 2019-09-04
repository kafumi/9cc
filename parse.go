package main

type Var struct {
	name   []rune
	offset int
}

func newVar(name []rune) *Var {
	str := string(name)
	if _, exist := env.vars[str]; exist {
		fatal("Variable \"%s\" is already defined", str)
	}
	v := &Var{
		name:   name,
		offset: env.maxOffset + 8,
	}
	env.vars[str] = v
	env.maxOffset = v.offset
	return v
}

func findVar(name []rune) *Var {
	return env.vars[string(name)]
}

type Env struct {
	vars      map[string]*Var
	maxOffset int
}

var env = &Env{
	vars: make(map[string]*Var),
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
	ndIf            // "if"
	ndWhile         // "while"
	ndFor           // "for"
	ndBlock         // { ... }
	ndReturn        // "return"
	ndFcall         // Function call
	ndLvar          // Local variable
	ndNum           // Integer
)

type Node struct {
	kind NodeKind

	lhs *Node // Left-hand side
	rhs *Node // Right-hand side

	// "if", "while", "for" statement
	test *Node
	cons *Node
	alt  *Node
	init *Node
	post *Node

	// Block
	body []*Node

	// Function call
	funcName string
	args     []*Node

	// Variable
	offset int

	// Number literal
	val int
}

func newNode(kind NodeKind, lhs *Node, rhs *Node) *Node {
	return &Node{
		kind: kind,
		lhs:  lhs,
		rhs:  rhs,
	}
}

func newNodeIf(test *Node, cons *Node, alt *Node) *Node {
	return &Node{
		kind: ndIf,
		test: test,
		cons: cons,
		alt:  alt,
	}
}

func newNodeWhile(test *Node, cons *Node) *Node {
	return &Node{
		kind: ndWhile,
		test: test,
		cons: cons,
	}
}

func newNodeFor(init *Node, test *Node, post *Node, cons *Node) *Node {
	return &Node{
		kind: ndFor,
		init: init,
		test: test,
		post: post,
		cons: cons,
	}
}

func newNodeBlock(body []*Node) *Node {
	return &Node{
		kind: ndBlock,
		body: body,
	}
}

func newNodeFcall(name []rune, args []*Node) *Node {
	return &Node{
		kind:     ndFcall,
		funcName: string(name),
		args:     args,
	}
}

func newNodeLVar(name []rune) *Node {
	lvar := findVar(name)
	if lvar == nil {
		lvar = newVar(name)
	}
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
	var node *Node
	if consumeKind(tkIf) != nil {
		expect("(")
		test := expr()
		expect(")")
		cons := stmt()
		var alt *Node
		if consumeKind(tkElse) != nil {
			alt = stmt()
		}
		node = newNodeIf(test, cons, alt)
	} else if consumeKind(tkWhile) != nil {
		expect("(")
		test := expr()
		expect(")")
		cons := stmt()
		node = newNodeWhile(test, cons)
	} else if consumeKind(tkFor) != nil {
		var init, test, post *Node
		expect("(")
		if !consume(";") {
			init = expr()
			expect(";")
		}
		if !consume(";") {
			test = expr()
			expect(";")
		}
		if !consume(")") {
			post = expr()
			expect(")")
		}
		cons := stmt()
		node = newNodeFor(init, test, post, cons)
	} else if consumeKind(tkReturn) != nil {
		node = newNode(ndReturn, expr(), nil)
		expect(";")
	} else if consume("{") {
		var body []*Node
		for !consume("}") {
			body = append(body, stmt())
		}
		node = newNodeBlock(body)
	} else {
		node = expr()
		expect(";")
	}
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

	token := consumeKind(tkIdent)
	if token != nil {
		if consume("(") {
			var args []*Node
			firstArg := true
			for !consume(")") {
				if firstArg {
					firstArg = false
				} else {
					expect(",")
				}
				args = append(args, expr())
			}
			return newNodeFcall(token.str, args)
		}
		return newNodeLVar(token.str)
	}

	return newNodeNum(expectNumber())
}
