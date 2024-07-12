package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

/*
 * Now, you need to consider the reaction of the calculator when users enter expressions in the wrong format.
 * The program should only accept numbers, a plus + sign , a minus - sign , and two commands: /exit and /help.
 * It cannot accept all other characters, and it is necessary to warn the user about this.
 */
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()

		switch {
		case strings.HasPrefix(input, "/"):
			if commands(input) {
				return
			}
		default:
			handle(input)
		}
	}
}

/**
 * returns true if application should be exited
 */
func commands(input string) bool {
	switch input {
	case "/exit":
		{
			fmt.Print("Bye!")
			return true
		}
	case "/help":
		{
			fmt.Println(
				"The program calculates the sum of numbers\n" +
					"Write down an expression contains of numbers (positive and negative) and operators",
			)
		}
	default:
		{
			fmt.Println("Unknown command")
		}
	}

	return false
}

func format(input string) string {
	// 1. space formatting
	spaces := regexp.MustCompile(`\s+`)
	input = spaces.ReplaceAllString(input, " ")

	// 2. minuses / pluses formatting
	pluses := regexp.MustCompile(`\++`)
	doubleMinuses := regexp.MustCompile(`--`)
	plusMinus := regexp.MustCompile(`-\+|\+-`)
	plusDigit := regexp.MustCompile(`\+(\d)`)
	for {
		formatted := input

		// duplicate pluses
		formatted = pluses.ReplaceAllString(input, "+")

		// duplicate minuses
		formatted = doubleMinuses.ReplaceAllString(formatted, "+")

		// plus and minus
		formatted = plusMinus.ReplaceAllString(formatted, "-")

		// replace all pluses before digits
		formatted = plusDigit.ReplaceAllString(formatted, `$1`)

		if formatted == input {
			break
		}

		input = formatted
	}

	return input
}

func isValidOperator(input string) bool {
	return input == "+" || input == "-"
}

func handle(input string) {
	if len(input) == 0 {
		return
	}

	// remove duplicate space chars
	input = format(input)
	// split expression by elements
	expression := strings.Split(input, " ")
	// set result to 0 and default operator is "+" to handle first number
	sum, operator := 0, "+"
	for _, element := range expression {
		number, err := strconv.Atoi(element)
		if err != nil {
			// if number was evaluated, it's operator
			if isValidOperator(element) {
				operator = element
				continue
			} else {
				fmt.Println("Invalid expression")
				return
			}
		}

		// at this point, number is correct
		switch operator {
		case "+":
			sum += number
		case "-":
			sum -= number
		default:
			fmt.Println("Invalid expression")
			return
		}
		// erase operator
		operator = ""
	}

	fmt.Println(sum)
}
