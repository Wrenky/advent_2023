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
	grid := parseInput(data)
	pre := time.Now()
	nGrid := TiltNorth(grid)
	p1 := SumGrid(nGrid)
	post := time.Now()
	fmt.Printf("part1: %d in %s\n", p1, post.Sub(pre))

	// PART 2
	cache, pre := make(map[string]int), time.Now()
	g := grid
	i := 0
	var sg string
	for i = 1; i <= 1000000; i++ {
		g = RunCycle(g)
		sg = stringGrid(g)
		if _, ok := cache[sg]; ok {
			break
		}
		cache[sg] = i
	}
	// Run the last few cycles left
	cCount := (1000000000 - i) % (i - cache[sg])
	for j := 0; j < cCount; j++ {
		g = RunCycle(g)
	}
	post = time.Now()
	fmt.Printf("Part2: %d in %s\n", SumGrid(g), post.Sub(pre))
}

func stringGrid(grid [][]rune) string {
	var res string
	for _, v := range grid {
		res += string(v) + "\n"
	}
	return res
}

func SumGrid(grid [][]rune) int {
	return lo.Sum(lo.Map(helpers.Transpose(grid), func(in []rune, _ int) int {
		return sumRow(in)
	}))
}

func TiltNorth(grid [][]rune) [][]rune {
	vgrid := helpers.Transpose(grid)
	nGrid := [][]rune{}
	for _, v := range vgrid {
		n := TiltRow(v)
		nGrid = append(nGrid, n)
	}
	return helpers.Transpose(nGrid)
}
func TiltRow(in []rune) []rune {
	for i := 1; i < len(in); i++ {
		if in[i-1] == '.' && in[i] == 'O' {
			j := i
			for (j-1 >= 0) && (in[j-1] == '.') {
				in[j], in[j-1] = in[j-1], in[j]
				j--
			}
		}
	}
	return in
}

// Run a full cycle
func RunCycle(grid [][]rune) [][]rune {
	for i := 0; i < 4; i++ {
		tilted := TiltNorth(grid)
		grid = helpers.Rotate90(tilted)
	}
	return grid
}

func printGrid(grid [][]rune) {
	fmt.Printf("--------------------------------\n")
	for _, v := range grid {
		fmt.Printf("%s\n", string(v))
	}
	fmt.Printf("--------------------------------\n")
}

func sumRow(in []rune) int {
	answer := 0
	for i := 0; i < len(in); i++ {
		if in[i] == 'O' {
			answer += len(in) - i
		}
	}
	return answer
}

// This needs to change to match your actual input
func parseInput(input string) [][]rune {
	return lo.Map(strings.Split(input, "\n"), func(line string, _ int) []rune {
		return []rune(line)
	})
}
