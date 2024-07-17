package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
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
	input = regexp.MustCompile(`(\+|-|=|\*|\(|\)|\^)`).ReplaceAllString(input, ` $1 `)
	input = regexp.MustCompile(`\s+`).ReplaceAllString(input, " ")
	input = strings.TrimSpace(input)

	// 2. minuses / pluses formatting
	for {
		formatted := input
		formatted = regexp.MustCompile(`\+ ?\+`).ReplaceAllString(formatted, "+")
		formatted = regexp.MustCompile(`\* ?\*`).ReplaceAllString(formatted, "**")
		formatted = regexp.MustCompile(`/ ?/`).ReplaceAllString(formatted, "//")
		formatted = regexp.MustCompile(`\++`).ReplaceAllString(formatted, "+")
		formatted = regexp.MustCompile(`(\+ ?-)|(- ?\+)`).ReplaceAllString(formatted, "-")
		formatted = regexp.MustCompile(`- ?-`).ReplaceAllString(formatted, "+")
		formatted = regexp.MustCompile(`\+([0-9]+)`).ReplaceAllString(formatted, "$1")
		if formatted == input {
			break
		}

		input = formatted
	}

	if strings.HasPrefix(input, "-") {
		input = regexp.MustCompile("^- ?(.+)$").ReplaceAllString(input, `-$1`)
	}

	return input
}

func isValidOperator(input string) bool {
	return input == "+" || input == "-" || input == "*" || input == "/" || input == "^"
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

func handleExpression(expression expression) (string, error) {
	result, err := handleEvaluation(expression, expression.input, false)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(result), nil
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

func isValidNumber(input string) bool {
	return regexp.MustCompile(`^-?[0-9]+$`).MatchString(input)
}

func toPostfix(infix []string) ([]string, error) {
	result, stack, i, j := make([]string, len(infix)), make([]string, len(infix)), 0, 0
	precedence := map[string]int{"+": 1, "-": 1, "*": 2, "/": 2, "^": 3}
	for _, element := range infix {
		switch {
		// Add operands (numbers and variables) to the result (postfix notation) as they arrive.
		case isValidNumber(element) || isValidIdentifier(element):
			result[i] = element
			i++
		// If the incoming element is a left parenthesis, push it on the stack.
		case element == "(":
			stack[j] = element
			j++
		// If the incoming element is a right parenthesis, pop the stack and add operators to the result until you see a left parenthesis.
		case element == ")":
			for j >= 0 {
				j--
				if j < 0 {
					return nil, errors.New("Invalid expression")
				}
				stackElement := stack[j]
				if stackElement == "(" {
					// Discard the pair of parentheses.
					break
				}

				result[i] = stackElement
				i++
			}
		case isValidOperator(element):
			// If the stack is empty or contains a left parenthesis on top, push the incoming operator on the stack.
			// If the incoming operator has higher precedence than the top of the stack, push it on the stack.
			for j >= 0 {
				if j == 0 || stack[j-1] == "(" || precedence[element] > precedence[stack[j-1]] {
					stack[j] = element
					j++
					break
				}

				// If the precedence of the incoming operator is lower than or equal to that of the top of the stack,
				// pop the stack and add operators to the result until you see an operator that has smaller precedence
				// or a left parenthesis on the top of the stack; then add the incoming operator to the stack.
				result[i] = stack[j-1]
				j--
				i++
			}
		default:
			return nil, errors.New("Invalid expression")
		}
	}

	// At the end of the expression, pop the stack and add all operators to the result.
	for j > 0 {
		j--
		if stack[j] == "(" {
			return nil, errors.New("Invalid expression")
		}

		result[i] = stack[j]
		i++
	}

	return result, nil
}

func handleEvaluation(expression expression, evaluation string, isAssignment bool) (int, error) {
	parts, err := toPostfix(strings.Split(evaluation, " "))
	if err != nil {
		return 0, err
	}

	stack, i := make([]int, len(parts)), 0
	// When we have an expression in postfix notation, we can calculate it using another stack. To do that, scan the postfix expression from left to right:
	for _, element := range parts {
		if element == "" {
			break
		}

		switch {
		// 1. If the incoming element is a number, push it into the stack (the whole number, not a single digit!).
		case isValidNumber(element):
			value, err := strconv.Atoi(element)
			if err != nil {
				return 0, err
			}

			stack[i] = value
			i++
		// 2. If the incoming element is the name of a variable, push its value into the stack.
		case isIdentifier(element):
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

			stack[i] = value
			i++
		// 3. If the incoming element is an operator, then pop twice to get two numbers and perform the operation; push the result on the stack.
		case isValidOperator(element):
			if i < 2 {
				return 0, errors.New("Invalid expression")
			}

			value, number1, number2 := 0, stack[i-2], stack[i-1]
			switch element {
			case "+":
				value = number1 + number2
			case "-":
				value = number1 - number2
			case "*":
				value = number1 * number2
			case "/":
				value = number1 / number2
			case "^":
				value = int(math.Pow(float64(number1), float64(number2)))
			default:
				return 0, errors.New("Invalid operator")
			}

			stack[i-2] = value
			i--
		default:
			return 0, errors.New("Invalid expression")
		}

	}

	// 4. When the expression ends, the number on the top of the stack is a final result.
	return stack[0], nil
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
		return handleExpression(expression)
	}

	return "", errors.New("Unexpected error")
}
