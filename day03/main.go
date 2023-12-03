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

	pre, ans, post := time.Now(), sumPartNumbers(parsed), time.Now()
	fmt.Printf("Part1 answer: %d, in %s\n", ans, post.Sub(pre))

	pre, ans, post = time.Now(), sumGearRatio(parsed), time.Now()
	fmt.Printf("Part2 answer: %d, in %s\n", ans, post.Sub(pre))

}

// Just make it a string array
func parseInput(input string) []string {
	return strings.Split(input, "\n")
}

type coord struct {
	col int
	row int
}

// This is the part1 code. Just sum up the gears that are touching symbols!
func sumPartNumbers(in []string) int {
	return lo.Reduce(filterCoord(in, isSymbol), func(agg int, c coord, _ int) int {
		// First, get all the adjacent parts to our current symbol
		adjacentParts := lo.Filter(
			surroundingPoints(c),
			func(c coord, _ int) bool {
				return isDigit(in[c.col][c.row])
			},
		)
		partNumbers := lo.Uniq(lo.Map(
			adjacentParts,
			func(c coord, _ int) int {
				return getPart(in, c)
			},
		))

		// Then add them to the sum
		return agg + lo.Sum(partNumbers)
	}, 0)
}

// This is the part2 code. Find a "gear" (* touching two parts), then multiply those parts. Sum all gears!
func sumGearRatio(in []string) int {
	return lo.Reduce(filterCoord(in, isCog), func(agg int, c coord, _ int) int {
		// Fetch all parts around our cog
		adjacentParts := lo.Filter(
			surroundingPoints(c),
			func(c coord, _ int) bool {
				return isDigit(in[c.col][c.row])
			},
		)
		// Expand the touching parts, toss duplicates
		partNumbers := lo.Uniq(lo.Map(
			adjacentParts,
			func(c coord, _ int) int {
				return getPart(in, c)
			},
		))
		// If they are two, add them to the sumation!
		if len(partNumbers) == 2 {
			return agg + (partNumbers[0] * partNumbers[1])
		}
		return agg
	}, 0)
}

// Just iterate up/down based on offset, checking if its valid and a digit.
func findBounds(s string, start int, offset int) int {
	if (start > 0) && (start < len(s)-1) && isDigit(s[start+offset]) {
		return findBounds(s, start+offset, offset)
	}
	return start
}

// Given a coordinate we know is a digit, find the full number
func getPart(in []string, c coord) int {
	lower, upper := findBounds(in[c.col], c.row, -1), findBounds(in[c.col], c.row, 1)
	i, _ := strconv.Atoi(in[c.col][lower : upper+1])
	return i
}

// FilterCoord takes a predicate and returns all matching coordinates.
func filterCoord(in []string, pred func(byte) bool) []coord {
	ans := lo.Flatten(lo.FilterMap(in, func(ele string, col int) ([]coord, bool) {
		res := lo.FilterMap([]byte(ele), func(r byte, row int) (coord, bool) {
			return coord{col, row}, pred(r)
		})
		return res, (len(res) != 0)
	}))
	return ans
}

// Helper predicates for filters
func isCog(char byte) bool {
	return char == '*'
}
func isSymbol(char byte) bool {
	return char != '.' && !isDigit(char)
}
func isDigit(c byte) bool {
	return ((c >= 48) && (c <= 57))
}

// Really we are just looking for the 8 points around my input point,
// so we just create a slice of modifications and add it to the point.
func surroundingPoints(c coord) []coord {
	addVal := []coord{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}
	return lo.Map(addVal, func(n coord, _ int) coord {
		return coord{c.col + n.col, c.row + n.row}
	})
}
