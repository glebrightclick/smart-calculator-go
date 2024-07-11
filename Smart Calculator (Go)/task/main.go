package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

/*
 * Write a program that reads two numbers in a loop and prints the sum in the standard output.
 * The program should print the same number if a user enters only a single number.
 * If a user enters an empty line, the program should ignore it and continue the loop until the user enters a number.
 * When the command /exit is entered, the program must print "Bye!" (without quotes), and then stop.
 * Add to the calculator the ability to read an unlimited sequence of numbers.
 * Add a /help command to print some information about the program.
 * If you encounter an empty line, do not output anything.
 */
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()

		switch input {
		case "/exit":
			{
				fmt.Print("Bye!")
				return
			}
		case "/help":
			{
				fmt.Print(
					"The program calculates the sum of numbers\n" +
						"Write down an expression contains of numbers (positive and negative) and operators\n",
				)
			}
		case "":
			break
		default:
			handle(input)
		}
	}
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

func handle(input string) {
	// remove duplicate space chars
	input = format(input)

	expression := strings.Split(input, " ")
	if len(expression) == 1 {
		number, err := strconv.Atoi(input)
		if err != nil {
			return
		}

		fmt.Println(number)
	} else {
		sum, operator := 0, "+"
		for _, element := range expression {
			number, err := strconv.Atoi(element)
			if err != nil {
				// if number was evaluated, it's operator
				operator = element
				continue
			}

			switch operator {
			case "+":
				sum += number
			case "-":
				sum -= number
			default:
				log.Fatal("Unexpected operator")
			}
		}

		fmt.Println(sum)
	}
}
