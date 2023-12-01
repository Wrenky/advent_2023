package main

import (
	_ "embed"
	"fmt"
	"strings"
	"time"
)

//go:embed demo
var data string

func init() {
	// Strip trailing newline
	data = strings.TrimRight(data, "\n")
	if len(data) == 0 {
		panic("No input file")
	}
}

func main() {
	pre, parsed, post := time.Now(), parseInput(data), time.Now()
	fmt.Printf("Data parsing took %s\n", post.Sub(pre))

	pre, ans, post := time.Now(), len(parsed), time.Now()
	fmt.Printf("Part1 answer: %d, in %s\n", ans, post.Sub(pre))

	pre, ans, post = time.Now(), 0, time.Now()
	fmt.Printf("Part2 answer: %d, in %s\n", ans, post.Sub(pre))

}

// This needs to change to match your actual input
func parseInput(input string) string {
	return input
}
