package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
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
	position   scanner.Position
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
		return "", errors.New("Error reading regex")
	}
	if b {
		return string(token), nil
	}
	return "", errors.New("Invalid token name")
}

func readInteger(token []byte) (int, error) {
	matches, err := regexp.Match("[0-9]+", token)
	if err != nil {
		return 0, err
	}
	if matches {
		i, err := strconv.Atoi(string(token))
		return i, err
	}
	return 0, errors.New("Invalid token")
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
		l.next = s.Peek()
		l.position = s.Pos()

		switch s.TokenText() {
		case "?":
			item.operator = "IFZERO"
		case "=":
			item.operator = "ASSIGNMENT"
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
				item.variable = char
			}

			i, err := readInteger([]byte(s.TokenText()))
			if err == nil {
				item.argument = i
			}
		}
	}
	fmt.Printf("%+v\n", l)
}
