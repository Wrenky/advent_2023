package main

import (
	_ "embed"
	"fmt"
	"strings"
	"time"

	"github.com/samber/lo"
	lop "github.com/samber/lo/parallel"
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
	parsed := parseInput(data)
	pre := time.Now()
	v, f := make(map[hashed]int), make(map[coord]int)
	beam(parsed, &v, &f, coord{0, 0}, RIGHT)
	ans := len(lo.Keys(f))
	post := time.Now()
	fmt.Printf("Part1 answer: %d, in %s\n", ans, post.Sub(pre))

	pre = time.Now()
	starts := generateStarts(parsed)
	p2 := lo.Max(lop.Map(starts, func(h hashed, _ int) int {
		v, f := make(map[hashed]int), make(map[coord]int)
		beam(parsed, &v, &f, h.c, h.d)
		return len(lo.Keys(f))
	}))
	post = time.Now()
	fmt.Printf("Part2 answer: %d, in %s\n", p2, post.Sub(pre))

}

func generateStarts(grid [][]rune) []hashed {
	res := []hashed{}
	for i := 0; i < len(grid); i++ {
		res = append(res, hashed{coord{0, i}, DOWN})
		res = append(res, hashed{coord{len(grid) - 1, i}, UP})
	}
	for i := 0; i < len(grid[0]); i++ {
		res = append(res, hashed{coord{i, 0}, RIGHT})
		res = append(res, hashed{coord{i, len(grid[0]) - 1}, LEFT})
	}
	return res
}

type coord struct {
	col, row int
}

type hashed struct {
	c coord
	d int
}

const (
	UP = iota
	DOWN
	LEFT
	RIGHT
)

func beam(grid [][]rune, visited *map[hashed]int, agg *map[coord]int, start coord, startDir int) {
	// Then check the next direction till event, process event
	var beamR func(c coord, dir int)
	beamR = func(c coord, dir int) {
		curr := c
		for {
			// First, check if valid
			if !isValid(grid, curr) {
				return
			}
			h := hashed{c: curr, d: dir}
			if _, ok := (*visited)[h]; ok {
				return
			}
			(*visited)[h]++
			(*agg)[curr]++

			if grid[curr.col][curr.row] != '.' {
				break
			}
			curr = next(curr, dir)
		}
		// We have an action point
		switch grid[curr.col][curr.row] {
		case '/':
			fallthrough
		case '\\':
			//fmt.Printf("Found a mirror at %v heading %d\n", curr, dir)
			nextCoord, nextDir := mirror(grid[curr.col][curr.row], curr, dir)
			beamR(nextCoord, nextDir)
		case '|':
			if dir == LEFT || dir == RIGHT {
				beamR(next(curr, UP), UP)
				beamR(next(curr, DOWN), DOWN)
			} else {
				beamR(next(curr, dir), dir)
			}
		case '-':
			if dir == UP || dir == DOWN {
				beamR(next(curr, RIGHT), RIGHT)
				beamR(next(curr, LEFT), LEFT)
			} else {
				beamR(next(curr, dir), dir)
			}
		}
	}
	beamR(start, startDir)
	return
}

// Handle coordinate moves
func next(c coord, dir int) coord {
	u, d, l, r := coord{-1, 0}, coord{1, 0}, coord{0, -1}, coord{0, 1}
	add := func(a, b coord) coord {
		return coord{
			row: a.row + b.row,
			col: a.col + b.col,
		}
	}
	switch dir {
	case UP:
		return add(c, u)
	case DOWN:
		return add(c, d)
	case LEFT:
		return add(c, l)
	case RIGHT:
		return add(c, r)
	}
	return c
}

func mirror(mirror rune, curr coord, dir int) (coord, int) {
	var forward = map[int]int{
		UP:    RIGHT,
		DOWN:  LEFT,
		LEFT:  DOWN,
		RIGHT: UP,
	}
	var back = map[int]int{
		UP:    LEFT,
		DOWN:  RIGHT,
		LEFT:  UP,
		RIGHT: DOWN,
	}
	if mirror == '/' {
		return next(curr, forward[dir]), forward[dir]
	} else {
		return next(curr, back[dir]), back[dir]
	}
}

func isValid(grid [][]rune, c coord) bool {
	if (c.row < len(grid[0]) && c.row >= 0) &&
		(c.col < len(grid) && c.col >= 0) {
		return true
	}
	return false
}

// This needs to change to match your actual input
func parseInput(input string) [][]rune {
	return lo.Map(strings.Split(input, "\n"), func(line string, _ int) []rune {
		return []rune(line)
	})
}
