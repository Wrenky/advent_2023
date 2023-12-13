package main

import (
	"advent/helpers"
	_ "embed"
	"fmt"
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
	blockPatterns := parseInput(data)

	pre := time.Now()
	p1 := lo.Sum(lo.Map(blockPatterns, func(grid [][]rune, _ int) int {
		return detectMirror(grid, false)
	}))
	post := time.Now()
	fmt.Printf("Part1: %d in %s\n", p1, post.Sub(pre))

	// Part 2 is just part 1 with swaps enabled
	pre = time.Now()
	p2 := lo.Sum(lo.Map(blockPatterns, func(grid [][]rune, _ int) int {
		return detectMirror(grid, true)
	}))
	post = time.Now()
	fmt.Printf("Part2: %d in %s\n", p2, post.Sub(pre))
}

func detectMirror(grid [][]rune, swaps bool) int {
	vgrid := helpers.Transpose(grid)
	hv := lo.Map([][][]rune{grid, vgrid}, func(grid [][]rune, _ int) int {
		return lo.Sum(lo.Map(lo.RangeFrom(1, len(grid)-1), func(i int, _ int) int {
			if checkMirrorPoint(grid, i, swaps) {
				return i
			}
			return 0
		}))
	})
	// Its always either 0 or the column/row, so just sum them.
	return hv[0]*100 + hv[1]
}

func checkMirrorPoint(grid [][]rune, start int, swaps bool) bool {
	checkRange := min(start, len(grid)-start)
	low, high := start-1, start
	swap := false
	for checkRange != 0 {
		checkRange--
		l, r := grid[low], grid[high]
		for i, v := range l {
			if r[i] != v {
				if !swap && swaps {
					// Only swap once, if swaps enabled
					swap = true
				} else {
					return false
				}
			}
		}
		low--
		high++
	}
	if swaps && !swap {
		// If swaps are enabled, we need a swap
		// for a run to be valid
		return false
	}
	return true
}

func parseInput(input string) [][][]rune {
	return lo.Map(strings.Split(input, "\n\n"), func(block string, _ int) [][]rune {
		return lo.Map(strings.Split(block, "\n"), func(line string, _ int) []rune {
			return []rune(line)
		})
	})
}
