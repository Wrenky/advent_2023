package main

import (
	_ "embed"
	"fmt"
	"math"
	"strings"

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

type coord struct {
	y int
	x int
}

func (c *coord) String() string {
	return fmt.Sprintf("{%d, %d}", c.y, c.x)
}

func main() {
	grid, start := parseInput(data)
	p1 := findFarthestPoint(grid, start)
	fmt.Printf("Part1: %d\n", p1)
	fmt.Printf("Part2: %d\n", picks(vertexPath(start, grid)))
}

func findFarthestPoint(grid [][]rune, start coord) int {
	v := make(map[coord]int)
	dfs(start, 0, grid, &v)
	return lo.Max(lo.Values(v))
}

func dfs(curr coord, steps int, grid [][]rune, visited *map[coord]int) {
	if v, ok := (*visited)[curr]; ok {
		if steps < v {
			// We are doing better, continue
			(*visited)[curr] = steps
		} else {
			return
		}
	}
	(*visited)[curr] = steps
	for _, next := range nextPoints(curr, grid) {
		dfs(next, steps+1, grid, visited)
	}
	return
}

func vertexPath(start coord, grid [][]rune) []coord {
	res := []coord{start}
	prev := start
	// Just pick one
	curr := (nextPoints(start, grid))[0]
	for curr != start {
		res = append(res, curr)
		next := lo.Filter(nextPoints(curr, grid), func(c coord, _ int) bool {
			return c != prev
		})
		prev = curr
		curr = next[0]
	}
	return res
}

func nextPoints(curr coord, grid [][]rune) []coord {
	var addVal []coord
	switch grid[curr.y][curr.x] {
	case '|':
		addVal = []coord{{1, 0}, {-1, 0}}
	case '-':
		addVal = []coord{{0, 1}, {0, -1}}
	case 'L':
		addVal = []coord{{-1, 0}, {0, 1}}
	case 'J':
		addVal = []coord{{-1, 0}, {0, -1}}
	case '7':
		addVal = []coord{{1, 0}, {0, -1}}
	case 'F':
		addVal = []coord{{1, 0}, {0, 1}}
	case 'S':
		// Gotta figure if any surrounding points connect!
		addVal = surroundingPoints(curr, grid)
	}
	// Only one here
	return lo.Map(addVal, func(n coord, _ int) coord {
		return coord{curr.y + n.y, curr.x + n.x}
	})
}

func surroundingPoints(c coord, grid [][]rune) []coord {
	var res []coord
	up := grid[c.y-1][c.x]
	if up == '|' || up == '7' || up == 'F' {
		res = append(res, coord{-1, 0})
	}
	down := grid[c.y+1][c.x]
	if down == '|' || down == 'L' || down == 'J' {
		res = append(res, coord{1, 0})
	}
	left := grid[c.y][c.x-1]
	if left == '-' || left == 'L' || left == 'F' {
		res = append(res, coord{0, -1})
	}
	right := grid[c.y][c.x+1]
	if right == '-' || right == '7' || right == 'J' {
		res = append(res, coord{0, 1})
	}
	return res
}

// This needs to change to match your actual input
func parseInput(input string) ([][]rune, coord) {
	var start coord
	return lo.Map(strings.Split(input, "\n"), func(line string, y int) []rune {
		if lo.Contains([]rune(line), 'S') {
			start = coord{y: y + 1, x: lo.IndexOf([]rune(line), 'S') + 1}
		}
		return []rune(line)
	}), start
}

// Apparently the way to get the amount of interior points of a grid is called picks formula.
// To calculate picks formula, you need to get the area via the "shoelace" theorum
func picks(c []coord) int {
	return shoelace(c) - (len(c) / 2) + 1
}
func shoelace(c []coord) int {
	sum := 0
	p0 := c[len(c)-1]
	for _, p1 := range c {
		sum += p0.y*p1.x - p0.x*p1.y
		p0 = p1
	}
	res := math.Abs(float64(sum / 2))
	return int(res)
}
