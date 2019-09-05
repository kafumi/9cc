package main

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"unicode"
)

type TokenKind int

const (
	tkReserved TokenKind = iota // Reserved word or symbol
	tkIdent                     // Identifier
	tkNum                       // Integer
	tkEOF                       // End of input
)

type Token struct {
	kind TokenKind
	next *Token
	str  []rune
	pos  int
	val  int // Valid only if kind is tkNum
}

var token *Token

func peek(op string) bool {
	return token.kind == tkReserved && reflect.DeepEqual(token.str, []rune(op))
}

func consume(op string) bool {
	if token.kind == tkReserved && reflect.DeepEqual(token.str, []rune(op)) {
		token = token.next
		return true
	}
	return false
}

func consumeKind(kind TokenKind) *Token {
	if token.kind == kind {
		consumed := token
		token = token.next
		return consumed
	}
	return nil
}

func expect(op string) {
	if token.kind == tkReserved && reflect.DeepEqual(token.str, []rune(op)) {
		token = token.next
	} else {
		fatalAt(token.pos, "Next token is not \"%s\"", op)
	}
}

func expectKind(kind TokenKind) *Token {
	if token.kind == kind {
		expected := token
		token = token.next
		return expected
	}
	fatalAt(token.pos, "Unexpected next token")
	return nil
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
	head := Token{next: nil}
	cur := &head

	for pos, length := 0, len(p); pos < length; {
		// Space
		if unicode.IsSpace(p[pos]) {
			pos++
			continue
		}

		// 1 or 2 character symbol
		l := isReservedSymbol(p, pos)
		if l > 0 {
			cur = newToken(tkReserved, cur, p[pos:pos+l], pos)
			pos += l
			continue
		}

		// Reserved word (e.g. "if")
		l = isReservedWord(p, pos)
		if l > 0 {
			cur = newToken(tkReserved, cur, p[pos:pos+l], pos)
			pos += l
			continue
		}

		// Variable name
		l = isIdent(p, pos)
		if l > 0 {
			cur = newToken(tkIdent, cur, p[pos:pos+l], pos)
			pos += l
			continue
		}

		// Number
		l = isNumber(p, pos)
		if l > 0 {
			str := p[pos : pos+l]
			num, err := strconv.Atoi(string(str))
			if err != nil {
				fatalAt(pos, "Expect number")
			}
			cur = newToken(tkNum, cur, str, pos)
			cur.val = num
			pos += l
			continue
		}

		fatalAt(pos, "Unable to tokenize")
	}

	newToken(tkEOF, cur, []rune{}, len(p))
	return head.next
}

func isReservedSymbol(p []rune, pos int) int {
	remain := len(p) - pos

	if remain >= 2 {
		switch string(p[pos : pos+2]) {
		case "<=", ">=", "==", "!=":
			return 2
		}
	}

	switch p[pos] {
	case '+', '-', '*', '/', '&', '(', ')', '<', '>', '=', '{', '}', ';', ',':
		return 1
	}

	return 0
}

func isReservedWord(p []rune, pos int) int {
	words := []string{
		"int",
		"if",
		"else",
		"while",
		"for",
		"return",
	}

	remain := len(p) - pos
	for _, word := range words {
		runes := []rune(word)
		l := len(runes)
		if l > remain {
			continue
		}
		if !reflect.DeepEqual(p[pos:pos+l], runes) {
			continue
		}
		if l < remain && isTokenChar(p[pos+l]) {
			continue
		}
		return l
	}
	return 0
}

func isIdent(p []rune, pos int) int {
	if !isTokenFirstChar(p[pos]) {
		return 0
	}

	end := pos + 1
	for end < len(p) && isTokenChar(p[end]) {
		end++
	}
	return end - pos
}

func isNumber(p []rune, pos int) int {
	end := pos
	for end < len(p) && unicode.IsDigit(p[end]) {
		end++
	}
	return end - pos
}

func isTokenFirstChar(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func isTokenChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

func printTokens() {
	for t := token; t != nil; t = t.next {
		fmt.Fprintf(os.Stderr, "token: kind=%d, str=\"%s\", pos=%d, val=%d\n",
			t.kind, string(t.str), t.pos, t.val)
	}
}
