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

// variables is a map of variables to interface values
var variables map[Variable]interface{}

// program holds all of the statements until they are executed
var program = []Statement{}

// Statement is the struct that holds statement information
type Statement struct {
	variable string
	operator string
	value    interface{}
}

func main() {
	var s scanner.Scanner
	s.Init(strings.NewReader(src))
	s.Filename = "example"
	s.Whitespace ^= 1<<'\t' | 1<<'\n'

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		current := s.TokenText()
		statement := &Statement{}

		switch current {
		case "?":
			fmt.Println("\n ifzero")
			statement.operator = "?"
		case "=":
			fmt.Println("\n assignment")
			statement.operator = "="
			fmt.Printf("after assignment, statement is %+v\n", statement)
		case "+=":
			fmt.Println("\n addition")
			statement.operator = "+="
		case "-=":
			fmt.Println("\n subtraction")
			statement.operator = "+-"
		case "!":
			fmt.Println("\n print")
			statement.operator = "!"
		case ";":
			fmt.Println("\n endline")
			fmt.Printf("statement is %+v\n", statement)
		default:
			// handle case where it's not a token of ours
			s, _ := readVariable([]byte(s.TokenText()))

			fmt.Printf("variable: %+v\n", s)
		}
	}
}

func mapVariable(key string, value interface{}) {
	variables[Variable(key)] = value
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
