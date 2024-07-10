package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

/**
 * Write a program that reads two numbers in a loop and prints the sum in the standard output.
 * The program should print the same number if a user enters only a single number. If a user enters an empty line, the program should ignore it and continue the loop until the user enters a number.
 * When the command /exit is entered, the program must print "Bye!" (without quotes), and then stop.
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
		case "":
			break
		default:
			handle(input)
		}
	}
}

func handle(input string) {
	numbers := strings.Split(input, " ")
	if len(numbers) == 1 {
		number, err := strconv.Atoi(input)
		if err != nil {
			return
		}

		fmt.Println(number)
	} else {
		sum := 0
		for _, stringNumber := range numbers {
			number, err := strconv.Atoi(stringNumber)
			if err != nil {
				log.Fatal(err)
			}

			sum += number
		}

		fmt.Println(sum)
	}
}
