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
	pre := time.Now()
	parsed := parseInput(data)
	post := time.Now()
	parsed_elapsed := post.Sub(pre)
	fmt.Printf("Data parsing took %s\n", parsed_elapsed)

	pre = time.Now()
	p1 := part1(parsed)
	post = time.Now()
	fmt.Printf("Part1 answer: %d, in %s\n", p1, post.Sub(pre))

	pre = time.Now()
	p2 := part2(parsed)
	post = time.Now()
	fmt.Printf("Part2 answer: %d, in %s\n", p2, post.Sub(pre))

}

// This needs to change to match your actual input
func parseInput(input string) string {
	return input
}

func part1(input string) int {
	return 0
}
func part2(input string) int {
	return 0
}
