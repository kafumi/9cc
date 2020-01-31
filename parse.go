package main

type TypeKind int

const (
	tyInt = iota
	tyChar
	tyPtr
	tyArray
)

type Type struct {
	kind      TypeKind
	size      int   // sizeof
	ptrTo     *Type // tyPtr: referenced type, tyArray: element type
	arraySize int   // num of elements of array
}

var typeInt = &Type{kind: tyInt, size: 4}
var typeChar = &Type{kind: tyChar, size: 1}

func typePtrTo(ptrTo *Type) *Type {
	return &Type{
		kind:  tyPtr,
		size:  8,
		ptrTo: ptrTo,
	}
}

func typeArray(ptrTo *Type, arraySize int) *Type {
	return &Type{
		kind:      tyArray,
		size:      arraySize * ptrTo.size,
		ptrTo:     ptrTo,
		arraySize: arraySize,
	}
}

type Var struct {
	typ      *Type
	name     []rune
	offset   int // Valid only if isGlobal = false
	isGlobal bool
}

func newLocalVar(typ *Type, name []rune) *Var {
	str := string(name)
	if _, exist := env.vars[str]; exist {
		fatal("Variable \"%s\" is already defined", str)
	}
	v := &Var{
		typ:    typ,
		name:   name,
		offset: env.maxOffset + typ.size,
	}
	env.vars[str] = v
	env.maxOffset = v.offset
	return v
}

func newGlobalVar(typ *Type, name []rune) *Var {
	str := string(name)
	if _, exist := envGlobal.vars[str]; exist {
		fatal("Variable \"%s\" is already defined", str)
	}
	v := &Var{
		typ:      typ,
		name:     name,
		isGlobal: true,
	}
	envGlobal.vars[str] = v
	return v
}

func findVar(name []rune) *Var {
	str := string(name)
	v := env.vars[str]
	if v == nil {
		v = envGlobal.vars[str]
	}
	return v
}

type Env struct {
	vars      map[string]*Var
	maxOffset int
}

var env *Env
var envGlobal *Env

func newEnv() *Env {
	env = &Env{
		vars: make(map[string]*Var),
	}
	return env
}

type DataSeq int

type GlobalData struct {
	strings map[string]DataSeq
	nextSeq DataSeq
}

var globalData *GlobalData

func newGlobalData() *GlobalData {
	return &GlobalData {
		strings: make(map[string]DataSeq)
	}
}

func newString(val string) DataSeq {
	if seq, ok := globalData.strings[val]; ok {
		return seq
	} else {
		seq := globalData.nextSeq
		globalData.strings[val] = seq
		globalData.nextSeq = seq + 1
	}
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
	ndVar           // Variable
	ndNum           // Integer
	ndString        // String
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
	vble *Var

	// Number literal
	val int

	// String literal
	seq DataSeq
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

func newNodeVar(name []rune) *Node {
	v := findVar(name)
	if v == nil {
		fatal("Variable \"%s\" is not defined", string(name))
	}
	return &Node{
		kind: ndVar,
		vble: v,
	}
}

func newNodeNum(val int) *Node {
	return &Node{
		kind: ndNum,
		val:  val,
	}
}

func newNodeString(val string) *Node {
	seq := newString(val)
	return &Node {
		seq: seq,
	}
}

func nodeType(node *Node) *Type {
	switch node.kind {
	case ndEq, ndNe, ndLt, ndLe, ndMul, ndDiv, ndAssign, ndNum:
		return typeInt
	case ndAdd:
		ltype := nodeType(node.lhs)
		rtype := nodeType(node.rhs)
		lptr := (ltype.kind == tyPtr || ltype.kind == tyArray)
		rptr := (rtype.kind == tyPtr || rtype.kind == tyArray)
		if lptr && rptr {
			fatal("Can not add pointer type value to pointer type value")
		}
		if lptr {
			return ltype
		}
		return rtype
	case ndSub:
		ltype := nodeType(node.lhs)
		rtype := nodeType(node.rhs)
		lptr := (ltype.kind == tyPtr || ltype.kind == tyArray)
		rptr := (rtype.kind == tyPtr || rtype.kind == tyArray)
		if !lptr && rptr {
			fatal("Can not subtract pointer type value from non-pointer value")
		}
		if lptr && rptr {
			return typeInt
		}
		return ltype
	case ndAddr:
		return typePtrTo(nodeType(node.lhs))
	case ndDeref:
		derefNodeType := nodeType(node.lhs)
		if derefNodeType.kind != tyPtr && derefNodeType.kind != tyArray {
			fatal("Node %+v should be pointer type: %+v", node.lhs, derefNodeType)
		}
		return derefNodeType.ptrTo
	case ndFcall:
		return typeInt
	case ndVar:
		return node.vble.typ
	case ndString:
		return typePtrTo(typeChar)
	default:
		fatal("Node %+v don't have type", node)
		return nil
	}
}

func program() []*Function {
	envGlobal = newEnv()
	var funcs []*Function
	for !atEOF() {
		f := toplv()
		if f != nil {
			funcs = append(funcs, f)
		}
	}
	return funcs
}

func toplv() *Function {
	topTyp := typ()
	name := expectKind(tkIdent)

	if consume("(") {
		env := newEnv()
		var params []*Var
		firstParam := true
		for !consume(")") {
			if firstParam {
				firstParam = false
			} else {
				expect(",")
			}
			paramTyp := typ()
			ident := expectKind(tkIdent)
			params = append(params, newLocalVar(paramTyp, ident.str))
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

	for consume("[") {
		count := expectNumber()
		expect("]")
		topTyp = typeArray(topTyp, count)
	}
	newGlobalVar(topTyp, name.str)
	expect(";")

	return nil
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
		for consume("[") {
			count := expectNumber()
			expect("]")
			typ = typeArray(typ, count)
		}
		newLocalVar(typ, ident.str)
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

	node := primary()
	if consume("[") {
		node = newNode(ndDeref, newNode(ndAdd, node, expr()), nil)
		expect("]")
	}
	return node
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
		return newNodeVar(token.str)
	}

	token = consumeKind(tkString)
	if token != nil {
		
	}

	return newNodeNum(expectNumber())
}

func typ() *Type {
	token := expectKind(tkReserved)
	var typ *Type
	if token != nil {
		switch string(token.str) {
		case "int":
			typ = typeInt
		case "char":
			typ = typeChar
		}
	}
	if typ == nil {
		fatal("Expect type name but \"%s\" is unknown type name", string(token.str))
	}
	for consume("*") {
		typ = typePtrTo(typ)
	}
	return typ
}

func peekTyp() bool {
	return peek("int") || peek("char")
}
