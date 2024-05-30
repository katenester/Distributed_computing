package main

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

var precedence = map[rune]int{
	'+': 1,
	'-': 1,
	'*': 2,
	'/': 2,
}

func isOperator(c rune) bool {
	_, exists := precedence[c]
	return exists
}

func higherPrecedence(op1, op2 rune) bool {
	return precedence[op1] >= precedence[op2]
}

func toPostfix(expression string) (string, error) {
	var stack []rune
	var output []string
	var numberBuffer strings.Builder

	for i := 0; i < len(expression); {
		r, size := utf8.DecodeRuneInString(expression[i:])
		i += size

		switch {
		case unicode.IsDigit(r) || r == '.':
			numberBuffer.WriteRune(r)
		case r == '(':
			if numberBuffer.Len() > 0 {
				output = append(output, numberBuffer.String())
				numberBuffer.Reset()
			}
			stack = append(stack, r)
		case r == ')':
			if numberBuffer.Len() > 0 {
				output = append(output, numberBuffer.String())
				numberBuffer.Reset()
			}
			for len(stack) > 0 && stack[len(stack)-1] != '(' {
				output = append(output, string(stack[len(stack)-1]))
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				return "", fmt.Errorf("mismatched parentheses")
			}
			stack = stack[:len(stack)-1]
		case isOperator(r):
			if numberBuffer.Len() > 0 {
				output = append(output, numberBuffer.String())
				numberBuffer.Reset()
			}
			for len(stack) > 0 && isOperator(stack[len(stack)-1]) && higherPrecedence(stack[len(stack)-1], r) {
				output = append(output, string(stack[len(stack)-1]))
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, r)
		case unicode.IsSpace(r):
			if numberBuffer.Len() > 0 {
				output = append(output, numberBuffer.String())
				numberBuffer.Reset()
			}
		default:
			return "", fmt.Errorf("invalid character: %v", r)
		}
	}

	if numberBuffer.Len() > 0 {
		output = append(output, numberBuffer.String())
	}

	for len(stack) > 0 {
		if stack[len(stack)-1] == '(' {
			return "", fmt.Errorf("mismatched parentheses")
		}
		output = append(output, string(stack[len(stack)-1]))
		stack = stack[:len(stack)-1]
	}

	return strings.Join(output, " "), nil
}

func main() {
	expressions := []string{
		"3+4*2/(1-5)",
		"(1+2)*(3/4)-(5+6)",
		"10.5+2*6",
	}

	for _, expr := range expressions {
		postfix, err := toPostfix(expr)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("Infix: %s => Postfix: %s\n", expr, postfix)
		}
	}
}
