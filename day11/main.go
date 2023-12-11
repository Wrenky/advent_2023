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
	pre, grid, post := time.Now(), parseInput(data), time.Now()
	fmt.Printf("Data parsing took %s\n", post.Sub(pre))
	eRow, eCol := FindEmpty(grid)
	pairs := generatePairs(grid)
	distances := lo.Map(pairs, func(set []helpers.Coord, _ int) int {
		return expandDistance(set[0], set[1], eCol, eRow, 2)
	})

	fmt.Printf("Part1: %v\n", lo.Sum(distances))
	distances = lo.Map(pairs, func(set []helpers.Coord, _ int) int {
		return expandDistance(set[0], set[1], eCol, eRow, 1000000)
	})
	fmt.Printf("Part1: %v\n", lo.Sum(distances))
}

// This needs to change to match your actual input
func parseInput(input string) [][]rune {
	return lo.Map(strings.Split(input, "\n"), func(line string, _ int) []rune {
		return []rune(line)
	})
}

func expandDistance(a, b helpers.Coord, eCol, eRow []int, expand int) int {
	lowX, highX := min(a.X, b.X), max(a.X, b.X)
	lowY, highY := min(a.Y, b.Y), max(a.Y, b.Y)
	cMod := lo.Reduce(eCol, func(agg, v int, _ int) int {
		if (lowX <= v) && (highX >= v) {
			agg += expand - 1
		}
		return agg
	}, 0)
	rMod := lo.Reduce(eRow, func(agg, v int, _ int) int {
		if (lowY <= v) && (highY >= v) {
			agg += expand - 1
		}
		return agg
	}, 0)
	return a.ManhattanDist(b) + cMod + rMod
}

func generatePairs(grid [][]rune) [][]helpers.Coord {
	galaxies := []helpers.Coord{}
	for c, row := range grid {
		for r, char := range row {
			if char == '#' {
				galaxies = append(galaxies, helpers.Coord{X: r, Y: c})
			}
		}
	}
	// Now make em pairs!
	pairs := [][]helpers.Coord{}
	for i := 0; i < len(galaxies); i++ {
		for j := i + 1; j < len(galaxies); j++ {
			pairs = append(pairs, []helpers.Coord{galaxies[i], galaxies[j]})
		}
	}
	return pairs
}

func FindEmpty(grid [][]rune) ([]int, []int) {
	rows := findEmptyRows(grid)
	grid = transpose(grid)
	cols := findEmptyRows(grid)
	return rows, cols
}

func findEmptyRows(grid [][]rune) []int {
	updateList := []int{}
	for i, v := range grid {
		if !lo.Contains(v, '#') {
			updateList = append(updateList, i)
		}
	}
	return updateList
}

// Make it generic?
func transpose(slice [][]rune) [][]rune {
	xl := len(slice[0])
	yl := len(slice)
	result := make([][]rune, xl)
	for i := range result {
		result[i] = make([]rune, yl)
	}
	for i := 0; i < xl; i++ {
		for j := 0; j < yl; j++ {
			result[i][j] = slice[j][i]
		}
	}
	return result
}
