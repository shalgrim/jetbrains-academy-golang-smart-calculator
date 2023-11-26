package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func processCommand(command string) {
	switch command {
	case "exit":
		fmt.Println("Bye!")
		os.Exit(0)
	case "help":
		fmt.Println("The program is a calculator. You can add, subtract, multiply, and divide. It also takes parentheses. It also keeps track of variable assignments")
	default:
		fmt.Println("Unknown command")
	}
}

func processAssignment(input string, variables map[string]int) map[string]int {
	split := strings.Split(input, "=")
	if len(split) != 2 {
		fmt.Println("Invalid assignment")
		return variables
	}
	variable := strings.TrimSpace(split[0])
	sValue := strings.TrimSpace(split[1])

	if isValidVariable(sValue) {
		val, ok := variables[sValue]
		if ok {
			variables[variable] = val
		} else {
			fmt.Println("Unknown variable")
		}
	} else if !isValidVariable(variable) {
		fmt.Println("Invalid identifier")
	} else {
		iValue, err := strconv.Atoi(sValue)
		if err != nil {
			fmt.Println("Invalid assignment")
		} else {
			variables[variable] = iValue
		}
	}
	return variables
}

func isValidVariable(input string) bool {
	for _, char := range input {
		if !unicode.IsLetter(char) {
			return false
		}
	}
	return true
}

func processExpression(input string, variables map[string]int) {
	tokens := strings.Split(input, " ")
	sum := 0
	sign := 1
	numOrVarAcceptable := true
	for _, token := range tokens {
		if token == "" {
			continue
		} else if isValidVariable(token) {
			if !numOrVarAcceptable {
				fmt.Println("error A")
				return
			}
			val, ok := variables[token]
			if !ok {
				fmt.Println("Unknown variable")
				return
			}
			sum += sign * val
			sign = 1
			numOrVarAcceptable = false
		} else if token == "-" {
			sign *= -1
			numOrVarAcceptable = true
		} else if token == "+" {
			numOrVarAcceptable = true
		} else {
			if !numOrVarAcceptable {
				fmt.Println("error B")
				return
			}
			val, err := strconv.Atoi(token)
			if err != nil {
				fmt.Println("unknown error")
			}
			sum += val * sign
			sign = 1
			numOrVarAcceptable = false
		}
	}
	fmt.Println(sum)
}

func faukenizer(input string) []string {
	myInput := strings.ReplaceAll(input, "+", " + ")
	myInput = strings.ReplaceAll(myInput, "-", " - ")
	myInput = strings.ReplaceAll(input, "*", " * ")
	myInput = strings.ReplaceAll(myInput, "/", " / ")
	return strings.Split(myInput, " ")
}

func convertToPostfix(infix string) ([]string, error) {
	var stack Stack
	postfix := make([]string, 0, 100)
	tokens := faukenizer(infix)

	for _, token := range tokens {

		// empty token
		if token == "" {
			continue
		}

		// token is number or variable
		// 1. Add operands (numbers and variables) to the result (postfix notation) as they arrive.
		_, err := strconv.Atoi(token)
		if err == nil || isValidVariable(token) {
			postfix = append(postfix, token)
			continue
		}

		// token is +, -, *, or /
		if token == "+" || token == "-" || token == "*" || token == "/" {
			var lastPostfixElem string
			if len(postfix) > 0 {
				lastPostfixElem = postfix[len(postfix)-1]
			} else {
				lastPostfixElem = ""
			}

			// check for consecutive * or /
			if (token == "*" || token == "/") &&
				(lastPostfixElem == "*" || lastPostfixElem == "/") {
				fmt.Println("Invalid expression")
				return postfix, errors.New("invalid expression")
			}

			// check for consecutive + or -
			if token == "+" && lastPostfixElem == "+" {
				continue
			}
			if token == "-" && lastPostfixElem == "-" {
				if len(postfix) > 1 {
					secondToLastPostfixElem := postfix[len(postfix)-2]
					if secondToLastPostfixElem == "-" {
						// since I have three - in a row, remove one
						postfix = postfix[:len(postfix)-1]
						continue
					}
				}
			}

			// non-consecutive operator
			// 2. If the stack is empty or contains a left parenthesis on top, push the incoming operator on the stack.
			topStackElement, err := stack.TopElement()
			if err != nil || topStackElement == "(" {
				stack.Push(token)
			}

			// 3. If the incoming operator has higher precedence than the top of the stack, push it on the stack.
			if (token == "*" || token == "/") && (topStackElement == "+" || topStackElement == "-") {
				stack.Push(token)
			}

			// 4. If the precedence of the incoming operator is lower
			// than or equal to that of the top of the stack,
			// pop the stack and add operators to the result until
			// you see an operator that has smaller precedence or a
			// left parenthesis on the top of the stack;
			// then add the incoming operator to the stack.
		}
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	variables := make(map[string]int)

	for {
		scanner.Scan()
		input := scanner.Text()

		if input == "" {
			continue
		} else if string(input[0]) == "/" {
			processCommand(input[1:])
		} else if strings.Contains(input, "=") {
			variables = processAssignment(strings.TrimSpace(input), variables)
		} else {
			postfixExpression, err := convertToPostfix(input)
			if err != nil {
				continue
			}
			processPostfix(postfixExpression)
		}
	}
}
