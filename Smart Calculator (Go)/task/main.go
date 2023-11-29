package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type Stack struct {
	storage []string
}

func (s *Stack) Push(value string) {
	s.storage = append(s.storage, value)
}

func (s *Stack) Pop() (string, error) {
	last := len(s.storage) - 1
	if last <= -1 { // check the size
		return "", errors.New("Stack is empty") // and return error
	}

	value := s.storage[last]     // save the value
	s.storage = s.storage[:last] // remove the last element

	return value, nil // return saved value and nil error
}

func (s *Stack) TopElement() (string, error) {
	if len(s.storage) == 0 {
		return "", errors.New("Empty stack")
	}
	return s.storage[len(s.storage)-1], nil
}

func (s *Stack) IsEmpty() bool {
	return len(s.storage) == 0
}

type IntStack struct {
	storage []int
}

func (s *IntStack) Push(value int) {
	s.storage = append(s.storage, value)
}

func (s *IntStack) Pop() (int, error) {
	last := len(s.storage) - 1
	if last <= -1 { // check the size
		return 0, errors.New("Stack is empty") // and return error
	}

	value := s.storage[last]     // save the value
	s.storage = s.storage[:last] // remove the last element

	return value, nil // return saved value and nil error
}

func (s *IntStack) TopElement() (int, error) {
	if len(s.storage) == 0 {
		return 0, errors.New("Empty stack")
	}
	return s.storage[len(s.storage)-1], nil
}

func (s *IntStack) IsEmpty() bool {
	return len(s.storage) == 0
}

func isLowerOrEqualPrecedence(a string, b string) bool {
	return ((a == "+" || a == "-") && (b == "+" || b == "-" || b == "*" || b == "/")) ||
		((a == "+" || a == "-" || a == "*" || a == "/") && (b == "+" || b == "-" || b == "*" || b == "/"))
}

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
	myInput = strings.ReplaceAll(myInput, "*", " * ")
	myInput = strings.ReplaceAll(myInput, "/", " / ")
	myInput = strings.ReplaceAll(myInput, "(", " ( ")
	myInput = strings.ReplaceAll(myInput, ")", " ) ")
	initialTokens := strings.Split(myInput, " ")

	// now deal with consecutive operators
	// will probably run into problems with negative numbers
	var stack Stack
	for _, token := range initialTokens {
		topElement, err := stack.TopElement()
		if token == "" {
			continue
		} else if err != nil {
			stack.Push(token)
		} else if token == "+" {
			if topElement == "+" {
				continue
			}
			stack.Push(token)
		} else if token == "-" {
			if topElement == "+" {
				stack.Pop()
				stack.Push(token)
			} else if topElement == "-" {
				stack.Pop()
				stack.Push("+")
			} else {
				stack.Push(token)
			}
		} else {
			stack.Push(token)
		}
	}

	//var reverseStack Stack
	//
	//for !stack.IsEmpty() {
	//	val, err := stack.Pop()
	//	if err != nil {
	//		log.Fatal("WTF")
	//	}
	//	reverseStack.Push(val)
	//}

	return stack.storage
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

		if token == "+" || token == "-" || token == "*" || token == "/" {

			// let's deal with this consecutive operator stuff when we're operating on the postfix
			//var lastPostfixElem string
			//if len(postfix) > 0 {
			//	lastPostfixElem = postfix[len(postfix)-1]
			//} else {
			//	lastPostfixElem = ""
			//}
			//
			//// check for consecutive * or /
			//if (token == "*" || token == "/") &&
			//	(lastPostfixElem == "*" || lastPostfixElem == "/") {
			//	fmt.Println("Invalid expression")
			//	return postfix, errors.New("invalid expression")
			//}
			//
			//// check for consecutive + or -
			//if token == "+" && lastPostfixElem == "+" {
			//	// TODO next this is not accurate bc of the way it's built
			//	continue
			//}
			//if token == "-" && lastPostfixElem == "-" {
			//	if len(postfix) > 1 {
			//		secondToLastPostfixElem := postfix[len(postfix)-2]
			//		if secondToLastPostfixElem == "-" {
			//			// since I have three - in a row, remove one
			//			postfix = postfix[:len(postfix)-1]
			//			continue
			//		}
			//	}
			//}

			// non-consecutive operator
			// 2. If the stack is empty or contains a left parenthesis on top, push the incoming operator on the stack.
			topStackElement, err := stack.TopElement()
			if err != nil || topStackElement == "(" {
				stack.Push(token)
				continue
			}

			// 3. If the incoming operator has higher precedence than the top of the stack, push it on the stack.
			if (token == "*" || token == "/") && (topStackElement == "+" || topStackElement == "-") {
				stack.Push(token)
				continue
			}

			// 4. If the precedence of the incoming operator is lower
			// than or equal to that of the top of the stack,
			if isLowerOrEqualPrecedence(token, topStackElement) {
				// pop the stack and add operators to the result until
				// you see an operator that has smaller precedence or a
				// left parenthesis on the top of the stack;
				// then add the incoming operator to the stack.
				popped, err := stack.Pop()
				if err != nil {
					log.Fatal("Should not find empty stack here 1?")
				}
				postfix = append(postfix, popped)

				for {
					topStackElement, err = stack.TopElement()
					if err != nil {
						// stack must be empty now so push it and move on
						stack.Push(token)
						break
					}

					if !isLowerOrEqualPrecedence(token, topStackElement) || token == "(" {
						stack.Push(token)
						break
					} else {
						popped, err := stack.Pop()
						if err != nil {
							log.Fatal("Should not find empty stack here 3?")
						}
						postfix = append(postfix, popped)
					}
				}
			}
			continue
		}

		// 5. If the incoming element is a left parenthesis, push it on the stack.
		if token == "(" {
			stack.Push(token)
			continue
		}

		// 6. If the incoming element is a right parenthesis,
		// pop the stack and add operators to the result until you see a left parenthesis.
		// Discard the pair of parentheses.
		if token == ")" {
			popped, err := stack.Pop()
			if err != nil {
				return postfix, errors.New("problem after )")
				//log.Fatal("problem after )")
			}
			for popped != "(" {
				postfix = append(postfix, popped)
				if err != nil {
					log.Fatal("problem after ) 1")
				}
				popped, err = stack.Pop()
				if err != nil {
					return postfix, errors.New("problem after ) 2")
					//log.Fatal("problem after ) 2")
				}
			}
		}
	}

	// 7. At the end of the expression, pop the stack and add all operators to the result.
	for !stack.IsEmpty() {
		popped, err := stack.Pop()
		if err != nil {
			log.Fatal("problem in step 7")
		}
		postfix = append(postfix, popped)
	}

	return postfix, nil
}

func processPostfix(postfix []string, variables map[string]int) {
	var stack IntStack
	answer := 0

	for _, token := range postfix {
		// If the incoming element is a number, push it into the stack (the whole number, not a single digit!).
		iValue, err := strconv.Atoi(token)
		if err == nil {
			stack.Push(iValue)
			continue
		}

		// If the incoming element is the name of a variable, push its value into the stack.
		if isValidVariable(token) {
			val, ok := variables[token]
			if !ok {
				fmt.Println("Unknown variable")
				return
			} else {
				stack.Push(val)
				continue
			}
		}

		// If the incoming element is an operator,
		// then pop twice to get two numbers and perform the operation;
		// push the result on the stack.
		operator := token
		operand2, err := stack.Pop()
		if err != nil {
			fmt.Println("Invalid expression")
			return
		}
		operand1, err := stack.Pop()
		if err != nil {
			fmt.Println("Invalid expression")
			return
		}

		switch operator {
		// TODO: handle multiple +/-
		case "+":
			stack.Push(operand1 + operand2)
		case "-":
			stack.Push(operand1 - operand2)
		case "*":
			stack.Push(operand1 * operand2)
		case "/":
			stack.Push(operand1 / operand2)
		}
	}
	answer, err := stack.Pop()
	if err != nil {
		log.Fatal("Can't get answer")
	}
	fmt.Printf("%d\n", answer)
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
			// TODO: Will they ever be assigning an expression to a variable?
			variables = processAssignment(strings.TrimSpace(input), variables)
		} else {
			postfixExpression, err := convertToPostfix(input)
			if err != nil {
				fmt.Println("Invalid expression")
				continue
			}
			processPostfix(postfixExpression, variables)
		}
	}
}

// FOR TESTING TO INFIX:
// 1. 3 + 2 * 4 -> [3 2 4 * +]
// 2. 3 * 2 + 4 -> [3 2 * 4 +]
// 3. 2 * (3 + 4) + 1 -> [2 3 4 * 1 +] but should be [2 3 4 + * 1 +]
//func main() {
//	scanner := bufio.NewScanner(os.Stdin)
//	for {
//		scanner.Scan()
//		input := scanner.Text()
//
//		if input == "/exit" {
//			fmt.Println("Bye!")
//			os.Exit(0)
//		} else {
//			postfixExpression, err := convertToPostfix(input)
//			if err != nil {
//				continue
//			}
//			fmt.Println(postfixExpression)
//		}
//	}
//}
