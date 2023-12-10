package helpers

import (
	"fmt"
	"strconv"
)

// Math helpers!
func GCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
func LCM(a, b int) int {
	return ((a * b) / GCD(a, b))
}

// Atoi in AOC is usually only used in parsing, and after a regexp/split so you know its an int.
func Atoi(in string) int {
	i, err := strconv.Atoi(in)
	if err != nil {
		panic(fmt.Sprintf("helpers.Atoi recieved non integer string: %s", err))
	}
	return i
}

//Graphs: https://github.com/dominikbraun/graph
