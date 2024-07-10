package main

import (
	"fmt"
)

/**
 * Write a program that reads two integer numbers from the same line and prints their sum in the standard output. Numbers can be positive, negative, or zero.
 */
func main() {
	var num1, num2 int
	fmt.Scan(&num1, &num2)
	fmt.Println(num1 + num2)
}
