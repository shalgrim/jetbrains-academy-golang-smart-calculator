package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"unicode"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		sum := 0
		sign := 1
		currentNumber := 0
		signNeedsResetting := false
		scanner.Scan()
		input := scanner.Text()

		if input == "" {
			continue
		} else if input == "/exit" {
			fmt.Println("Bye!")
			os.Exit(0)
		} else if input == "/help" {
			fmt.Println("The program calculates the sum or difference of numbers")
		} else {
			for _, char := range input {
				if unicode.IsSpace(char) || unicode.IsControl(char) {
					sum += currentNumber * sign
					currentNumber = 0
					if signNeedsResetting {
						sign = 1
					}
				} else if unicode.IsNumber(char) {
					signNeedsResetting = true
					thisInt, err := strconv.Atoi(string(char))
					if err != nil {
						log.Fatal(err)
					}
					currentNumber *= 10
					currentNumber += thisInt
				} else {
					// must be either + or _
					sum += currentNumber * sign
					if signNeedsResetting {
						sign = 1
						signNeedsResetting = false
					}
					currentNumber = 0
					if char == '-' {
						sign *= -1
					}
					currentNumber = 0
				}
			}
			sum += currentNumber * sign
			fmt.Println(sum)
		}
	}
}
