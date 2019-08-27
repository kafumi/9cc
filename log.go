package main

import (
	"fmt"
	"os"
	"strings"
)

var userInput []rune

func fatal(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func fatalAt(pos int, format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s\n", string(userInput))
	fmt.Fprintf(os.Stderr, strings.Repeat(" ", pos))
	fmt.Fprintf(os.Stderr, "^ ")
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func fatalAtStr(loc []rune, format string, a ...interface{}) {
	fatalAt(len(userInput)-len(loc), format, a...)
}
