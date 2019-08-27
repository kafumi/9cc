package main

import (
	"reflect"
	"strconv"
	"unicode"
)

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

func newToken(kind TokenKind, cur *Token, str []rune, pos int) *Token {
	tok := &Token{
		kind: kind,
		pos:  pos,
		str:  str,
	}
	cur.next = tok
	return tok
}

func tokenize(p []rune) *Token {
	lenAll := len(p)
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
			cur = newToken(tkReserved, cur, p[:2], lenAll-len(p))
			p = p[2:]
			continue
		case '<', '>':
			if p[1] == '=' {
				cur = newToken(tkReserved, cur, p[:2], lenAll-len(p))
				p = p[2:]
			} else {
				cur = newToken(tkReserved, cur, p[:1], lenAll-len(p))
				p = p[1:]
			}
			continue
		case '+', '-', '*', '/', '(', ')':
			cur = newToken(tkReserved, cur, p[:1], lenAll-len(p))
			p = p[1:]
			continue
		}

		if unicode.IsDigit(p[0]) {
			val, length := readNumber(p)
			cur = newToken(tkNum, cur, p[:length], lenAll-len(p))
			cur.val = val
			p = p[length:]
			continue
		}

		fatalAtStr(p, "Unable to tokenize")
	}

	newToken(tkEOF, cur, p, lenAll)
	return head.next
}

func readNumber(program []rune) (int, int) {
	length := 0
	for length < len(program) && unicode.IsDigit(program[length]) {
		length++
	}

	target := string(program[0:length])
	number, err := strconv.Atoi(target)
	if err != nil {
		fatalAtStr(program, "Expect number")
	}

	return number, length
}
