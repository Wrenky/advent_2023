package main

import (
	aoc "advent/helpers"
	_ "embed"
	"fmt"
	"regexp"
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
	moves := parseInput(data)
	pre := time.Now()
	coords := FollowMoves(moves)
	p1 := len(coords) + aoc.PicksInnerPoints(coords)
	post := time.Now()
	fmt.Printf("Part1: %d in %s\n", p1, post.Sub(pre))

	pre = time.Now()
	moves = ConvertMoves(moves)
	coords = FollowMoves(moves)
	p2 := len(coords) + aoc.PicksInnerPoints(coords)
	post = time.Now()
	fmt.Printf("Part2: %d in %s\n", p2, post.Sub(pre))

}

func FollowMoves(moves []move) []aoc.Coord {
	curr := aoc.Coord{X: 0, Y: 0}
	res := []aoc.Coord{curr}
	for _, v := range moves {
		addC := coordMap[v.dir]
		for i := 0; i < v.amount; i++ {
			next := curr.Add(addC)
			if (next == aoc.Coord{X: 0, Y: 0}) {
				return res
			}
			res = append(res, next)
			curr = next
		}
	}
	return res
}

func ConvertMoves(moves []move) []move {
	newMoves := []move{}
	for _, v := range moves {
		r := []rune(v.color)
		distH, dir := string(r[0:len(r)-1]), aoc.Atoi(string(r[len(r)-1]))
		dist, _ := strconv.ParseInt(distH, 16, 64)
		newMoves = append(newMoves, move{
			dir:    dir,
			amount: int(dist),
			color:  v.color,
		})
	}
	return newMoves
}

type move struct {
	dir    int
	amount int
	color  string
}

// modified to match part2 lol
const (
	RIGHT = iota
	DOWN
	LEFT
	UP
)

var dirMap = map[string]int{
	"U": UP,
	"D": DOWN,
	"L": LEFT,
	"R": RIGHT,
}
var coordMap = map[int]aoc.Coord{
	UP:    {X: -1, Y: 0},
	DOWN:  {X: 1, Y: 0},
	LEFT:  {X: 0, Y: -1},
	RIGHT: {X: 0, Y: 1},
}

// This needs to change to match your actual input
func parseInput(input string) []move {
	pattern := regexp.MustCompile(`\(#(.*?)\)`)
	return lo.Map(strings.Split(input, "\n"), func(line string, _ int) move {
		entries := strings.Split(line, " ")
		color := pattern.FindStringSubmatch(entries[2])
		return move{
			dir:    dirMap[entries[0]],
			amount: aoc.Atoi(entries[1]),
			color:  color[1],
		}
	})
}
