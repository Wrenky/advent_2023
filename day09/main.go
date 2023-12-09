package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
)

//go:embed input
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

	// Part 1: Make the chain, then get the last digit. Sum all the last digits.
	pre = time.Now()
	p1 := lo.Reduce(parsed, func(agg int, row []int, _ int) int {
		return agg + getLast(chain(row))
	}, 0)
	post = time.Now()
	fmt.Printf("Part1 answer: %d, in %s\n", p1, post.Sub(pre))

	pre = time.Now()
	// Part 1: Make the chain, then get the first digit. Sum all the first digits.
	p2 := lo.Reduce(parsed, func(agg int, row []int, _ int) int {
		return agg + getFirst(chain(row))
	}, 0)
	post = time.Now()
	fmt.Printf("Part2 answer: %d, in %s\n", p2, post.Sub(pre))
}

// These two either get the first or last digit of the input set.
func getLast(in [][]int) int {
	return lo.Reduce(lo.Reverse(in), func(agg int, row []int, _ int) int {
		return agg + row[len(row)-1]
	}, 0)
}
func getFirst(in [][]int) int {
	return lo.Reduce(lo.Reverse(in), func(agg int, row []int, _ int) int {
		return row[0] - agg
	}, 0)
}

// Check that every element in the slice is zero.
func is0(in []int) bool {
	return lo.EveryBy(in, func(x int) bool {
		return x == 0
	})
}

// Get the cascading set of differences until they are all zeros.
func chain(in []int) [][]int {
	var res [][]int
	curr := in
	for !is0(curr) {
		res = append(res, curr)
		curr = difference(curr)
	}
	return res
}

// Get the difference between each element in the list
func difference(in []int) []int {
	var res []int
	for i := 1; i < len(in); i++ {
		res = append(res, in[i]-in[i-1])
	}
	return res
}

func parseInput(input string) [][]int {
	return lo.Map(strings.Split(input, "\n"), func(line string, _ int) []int {
		return lo.FilterMap(strings.Split(line, " "), func(ele string, _ int) (int, bool) {
			i, _ := strconv.Atoi(ele)
			return i, (ele != "")
		})
	})
}
