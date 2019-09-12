package main

type TypeKind int

const (
	tyInt = iota
	tyPtr
)

type Type struct {
	kind  TypeKind
	size  int
	ptrTo *Type
}

var typeInt = &Type{kind: tyInt, size: 4}

func typePtrTo(ptrTo *Type) *Type {
	return &Type{
		kind:  tyPtr,
		size:  8,
		ptrTo: ptrTo,
	}
}

type Var struct {
	typ    *Type
	name   []rune
	offset int
}

func newVar(typ *Type, name []rune) *Var {
	str := string(name)
	if _, exist := env.vars[str]; exist {
		fatal("Variable \"%s\" is already defined", str)
	}
	v := &Var{
		typ:    typ,
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

var env *Env

func newEnv() *Env {
	env = &Env{
		vars: make(map[string]*Var),
	}
	return env
}

type Function struct {
	name   []rune
	env    *Env
	params []*Var
	body   *Node
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
	ndAddr          // unary &
	ndDeref         // unary *
	ndIf            // "if"
	ndWhile         // "while"
	ndFor           // "for"
	ndBlock         // { ... }
	ndReturn        // "return"
	ndFcall         // Function call
	ndLvar          // Local variable
	ndNum           // Integer
	ndNull          // Null statement
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
	lvar *Var

	// Number literal
	val int
}

var nullNode = &Node{
	kind: ndNull,
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
		fatal("Variable \"%s\" is not defined", string(name))
	}
	return &Node{
		kind: ndLvar,
		lvar: lvar,
	}
}

func newNodeNum(val int) *Node {
	return &Node{
		kind: ndNum,
		val:  val,
	}
}

func nodeType(node *Node) *Type {
	switch node.kind {
	case ndEq, ndNe, ndLt, ndLe, ndMul, ndDiv, ndAssign, ndNum:
		return typeInt
	case ndAdd:
		ltype := nodeType(node.lhs)
		rtype := nodeType(node.rhs)
		if ltype.kind == tyPtr && rtype.kind == tyPtr {
			fatal("Can not add pointer type value to pointer type value")
		}
		if ltype.kind == tyPtr {
			return ltype
		}
		return rtype
	case ndSub:
		ltype := nodeType(node.lhs)
		rtype := nodeType(node.rhs)
		if ltype.kind == tyInt && rtype.kind == tyPtr {
			fatal("Can not subtract pointer type value from int value")
		}
		if ltype.kind == tyPtr && rtype.kind == tyPtr {
			return typeInt
		}
		return ltype
	case ndAddr:
		return typePtrTo(nodeType(node.lhs))
	case ndDeref:
		derefNodeType := nodeType(node.lhs)
		if derefNodeType.kind != tyPtr {
			fatal("Node %+v should be pointer type", node.lhs)
		}
		return derefNodeType.ptrTo
	case ndFcall:
		return typeInt
	case ndLvar:
		return node.lvar.typ
	default:
		fatal("Node %+v don't have type", node)
		return nil
	}
}

func program() []*Function {
	var funcs []*Function
	for !atEOF() {
		funcs = append(funcs, funct())
	}
	return funcs
}

func funct() *Function {
	typ()
	name := expectKind(tkIdent)
	env := newEnv()

	expect("(")
	var params []*Var
	firstParam := true
	for !consume(")") {
		if firstParam {
			firstParam = false
		} else {
			expect(",")
		}
		typ := typ()
		ident := expectKind(tkIdent)
		params = append(params, newVar(typ, ident.str))
	}

	expect("{")
	var stmts []*Node
	for !consume("}") {
		stmts = append(stmts, stmt())
	}
	body := newNodeBlock(stmts)

	return &Function{
		name:   name.str,
		env:    env,
		params: params,
		body:   body,
	}
}

func stmt() *Node {
	var node *Node
	if consume("if") {
		expect("(")
		test := expr()
		expect(")")
		cons := stmt()
		var alt *Node
		if consume("else") {
			alt = stmt()
		}
		node = newNodeIf(test, cons, alt)
	} else if consume("while") {
		expect("(")
		test := expr()
		expect(")")
		cons := stmt()
		node = newNodeWhile(test, cons)
	} else if consume("for") {
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
	} else if consume("return") {
		node = newNode(ndReturn, expr(), nil)
		expect(";")
	} else if peekTyp() {
		typ := typ()
		ident := expectKind(tkIdent)
		newVar(typ, ident.str)
		expect(";")
		node = nullNode
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
	if consume("&") {
		return newNode(ndAddr, unary(), nil)
	}
	if consume("*") {
		return newNode(ndDeref, unary(), nil)
	}
	if consume("sizeof") {
		node := unary()
		return newNodeNum(nodeType(node).size)
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

func typ() *Type {
	expect("int")
	typ := typeInt
	for consume("*") {
		typ = typePtrTo(typ)
	}
	return typ
}

func peekTyp() bool {
	return peek("int")
}
