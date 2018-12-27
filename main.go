package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"text/scanner"
)

// const src = `
// x = 3;
// y = 2;
// ! x;
// ! y;
// `

const src = `x=3;`

// Variable is string
type Variable string

// program holds all of the statements until they are executed
var program = []Statement{}

// Statement is the struct that holds statement information
type Statement struct {
	variable string
	operator string
	value    interface{}
}

// Lexer is the main struct that is worked on
type Lexer struct {
	variables  map[Variable]interface{}
	statements []Statement
	current    *Statement
	next       interface{}
	pos        int
}

func (l *Lexer) clearCurrent() {
	l.current = nil
}

func main() {
	var s scanner.Scanner
	s.Init(strings.NewReader(src))
	s.Filename = "jjdl"
	s.Whitespace ^= 1<<'\t' | 1<<'\n'

	l := &Lexer{}
	item := &Statement{}
	l.current = item

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		current := s.TokenText()
		fmt.Printf("position: %+v\n", s.Position)

		switch current {
		case "?":
			item.operator = "ifzero"
		case "=":
			item.operator = "assignment"
			val := scanner.TokenString(s.Peek())
			// TODO: Need to maybe parse this better?
			item.value = val
		case "+=":
			item.operator = "addition"
		case "-=":
			item.operator = "subtraction"
		case "!":
			item.operator = "print"
		case ";":
			fmt.Println("END OF STATEMENT")
			l.statements = append(l.statements, *item)
			l.current = &Statement{}
		default:
			char, _ := readVariable([]byte(s.TokenText()))
			if char != "" {
				l.current.variable = char
			}
		}
	}
	fmt.Printf("Lexer is %+v\n", l)
}

func (l *Lexer) mapVariable(key string, value interface{}) {
	l.variables[Variable(key)] = value
	return
}

func readVariable(token []byte) (string, error) {
	b, err := regexp.Match("[a-zA-Z][a-zA-Z0-9]*", token)
	if err != nil {
		log.Fatalf("error regexp match %+v\n", err)
		return "", errors.New("Error reading regex")
	}
	if b {
		// fmt.Printf("variable: %s\n", s.TokenText())
		return string(token), nil
	}
	return "", errors.New("Invalid token name")
}
