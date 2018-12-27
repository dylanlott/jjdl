package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"text/scanner"
)

const src = `
x = 3;
y = 2;
! x;
! y;
`

// Variable is string
type Variable string

// Statement is the struct that holds statement information
type Statement struct {
	variable string
	operator string
	argument interface{}
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
	l.current = &Statement{}
}

func main() {
	BuildStatements()
}

func (l *Lexer) mapVariable(key string, arg interface{}) {
	l.variables[Variable(key)] = arg
	return
}

func readVariable(token []byte) (string, error) {
	b, err := regexp.Match("[a-zA-Z][a-zA-Z0-9]*", token)
	if err != nil {
		log.Fatalf("error regexp match %+v\n", err)
		return "", errors.New("Error reading regex")
	}
	if b {
		return string(token), nil
	}
	return "", errors.New("Invalid token name")
}

// Statements returns all statements on the Lexer. If the lexer hasn't been run yet
// this won't be populated
func (l *Lexer) Statements() []Statement {
	return l.statements
}

// BuildStatements will run lex the input and build statements
// from the input.
func BuildStatements() {
	var s scanner.Scanner
	s.Init(strings.NewReader(src))
	s.Filename = "jjdl"

	l := &Lexer{}
	l.variables = make(map[Variable]interface{})
	item := &Statement{}
	l.current = item

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		current := s.TokenText()

		switch current {
		case "?":
			item.operator = "IFZERO"
		case "=":
			item.operator = "ASSIGNMENT"
			val := scanner.TokenString(s.Peek())
			// TODO: Need to maybe parse this better?
			item.argument = val
		case "+=":
			item.operator = "ADD"
		case "-=":
			item.operator = "SUB"
		case "!":
			item.operator = "PRINT"
		case ";":
			l.mapVariable(item.variable, item.argument)
			l.statements = append(l.statements, *item)
			l.clearCurrent()
		default:

			char, _ := readVariable([]byte(s.TokenText()))
			if char != "" {
				fmt.Printf("char is %+v\n", char)
				item.variable = char
			}
		}
	}
	fmt.Printf("Lexer is %+v\n", l)
}
