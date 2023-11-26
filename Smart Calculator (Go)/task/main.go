package main

import (
	"bufio"
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
		fmt.Println("The program calculates the sum or difference of numbers and keeps track of variable assignments")
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

func processVariableExpression(input string, variables map[string]int) {
	if isValidVariable(input) {
		val, ok := variables[input]
		if ok {
			fmt.Println(val)
		} else {
			fmt.Println("Unknown variable")
		}
	} else {
		fmt.Println("Invalid identifier")
	}
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
	return strings.Split(myInput, " ")
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
			tokens := faukenizer(input)
			processExpression(strings.Join(tokens, " "), variables)
		}
	}
}
