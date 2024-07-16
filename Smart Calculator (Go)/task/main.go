package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type expression struct {
	// original input from the user
	input string
	// variables and their values
	variables map[string]int
}

/*
 * Now, you need to consider the reaction of the calculator when users enter expressions in the wrong format.
 * The program should only accept numbers, a plus + sign , a minus - sign , and two commands: /exit and /help.
 * It cannot accept all other characters, and it is necessary to warn the user about this.
 */
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	expression := expression{variables: make(map[string]int)}
	for scanner.Scan() {
		input := scanner.Text()

		switch {
		case strings.HasPrefix(input, "/"):
			if commands(input) {
				return
			}
		default:
			// new expression input
			expression.input = input
			// if error occurred, display an error
			if result, err := handle(expression); err != nil {
				fmt.Println(err)
				continue
				// otherwise, display result if not empty
			} else if len(result) > 0 {
				fmt.Println(result)
			}
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
	input = strings.TrimSpace(input)
	input = regexp.MustCompile(`(=)`).ReplaceAllString(input, ` $1 `)
	input = regexp.MustCompile(`\s+`).ReplaceAllString(input, " ")

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

func isIdentifier(identifier string) bool {
	return regexp.MustCompile(`[a-zA-Z]+`).MatchString(identifier)
}

func isValidIdentifier(identifier string) bool {
	return regexp.MustCompile(`^[a-zA-Z]+$`).MatchString(identifier)
}

func isAssignment(input string) bool {
	return strings.Contains(input, "=")
}

func handleAssignment(expression expression) error {
	input := strings.Split(expression.input, " = ")
	if len(input) != 2 {
		return errors.New("Invalid assignment")
	}

	variable, evaluation := input[0], input[1]
	if !isValidIdentifier(variable) {
		return errors.New("Invalid identifier")
	}

	var result int
	var err error
	if result, err = handleEvaluation(expression, evaluation, true); err != nil {
		return err
	}

	// if identifier is correct and err is nil, result variable has correct result
	expression.variables[variable] = result
	return nil
}

func handleEvaluation(expression expression, evaluation string, isAssignment bool) (int, error) {
	parts := strings.Split(evaluation, " ")
	// set result to 0 and default operator is "+" to handle first number
	sum, operator := 0, "+"
	for _, element := range parts {
		number, err := strconv.Atoi(element)
		if err != nil {
			// if element is valid variable, pick its value and save it to number
			if isIdentifier(element) {
				if !isValidIdentifier(element) {
					if isAssignment {
						return 0, errors.New("Invalid assignment")
					}
					return 0, errors.New("Invalid identifier")
				}

				value, ok := expression.variables[element]
				if !ok {
					return 0, errors.New("Unknown variable")
				}

				number = value
				// if number was evaluated, it's operator
			} else if isValidOperator(element) {
				operator = element
				continue
			} else {
				return 0, errors.New("Invalid expression")
			}
		}

		// at this point, number is correct
		switch operator {
		case "+":
			sum += number
		case "-":
			sum -= number
		default:
			return 0, errors.New("Invalid expression")
		}
		// erase operator
		operator = ""
	}

	return sum, nil
}

// returns feedback to display and error
func handle(expression expression) (string, error) {
	if len(expression.input) == 0 {
		return "", nil
	}

	// remove duplicate space chars
	expression.input = format(expression.input)
	// split expression by elements
	if isAssignment(expression.input) {
		return "", handleAssignment(expression)
	} else {
		var result int
		var err error
		if result, err = handleEvaluation(expression, expression.input, false); err != nil {
			return "", err
		}

		return strconv.Itoa(result), nil
	}

	return "", errors.New("Unexpected error")
}
