package main

import (
	_ "embed"
	"fmt"
	"regexp"
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
	instructions, nodes := parseInput(data)
	graphMap := buildGraphMap(nodes)
	// Just traverse the map from AAA to ZZZ
	p1 := traverse("AAA", []string{"ZZZ"}, instructions, graphMap)
	fmt.Printf("Part 1 result: %d\n", p1)

	pre := time.Now()
	startingNodes := lo.Filter(lo.Keys(graphMap), func(i string, _ int) bool { return ([]rune(i))[2] == 'A' })
	endingNodes := lo.Filter(lo.Keys(graphMap), func(i string, _ int) bool { return ([]rune(i))[2] == 'Z' })
	// Find an ending Z for every starting node
	rates := lo.Map(startingNodes, func(start string, _ int) int {
		return traverse(start, endingNodes, instructions, graphMap)
	})
	// Then figure out their least common multiple for when the cycles align
	p2 := lo.Reduce(rates, func(agg int, i int, _ int) int {
		return lcm(agg, i)
	}, rates[0])
	post := time.Now()
	fmt.Printf("Part 2 result: %d in %s\n", p2, post.Sub(pre))
}

type Graph struct {
	Value string
	Right string
	Left  string
}

func buildGraphMap(in [][]string) map[string]Graph {
	gm := make(map[string]Graph)
	for _, v := range in {
		gm[v[0]] = Graph{
			Value: v[0],
			Left:  v[1],
			Right: v[2],
		}
	}
	return gm
}

// Math helpers!
func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
func lcm(a, b int) int {
	return ((a * b) / gcd(a, b))
}

func traverse(startNode string, endNodes []string, ins string, gm map[string]Graph) int {
	var i rune
	seen, next := 0, ""
	instructions, current := []rune(ins), gm[startNode]
	for {
		i, instructions = instructions[0], instructions[1:]
		if len(instructions) == 0 {
			instructions = []rune(ins)
		}
		if lo.Contains(endNodes, current.Value) {
			break
		}
		if i == 'L' {
			// Left path
			next = current.Left
		} else {
			// Right path
			next = current.Right
		}
		seen++
		current = gm[next]
	}
	return seen
}

// This needs to change to match your actual input
func parseInput(input string) (string, [][]string) {

	instructionPattern := regexp.MustCompile(`^([LR]+)\n`)
	instructions := instructionPattern.FindStringSubmatch(input)

	nodePat := regexp.MustCompile(`(\w+)\s=\s\((\w+),\s(\w+)\)`)
	nodeMatches := nodePat.FindAllStringSubmatch(input, -1)
	return instructions[1], lo.Map(nodeMatches, func(m []string, _ int) []string {
		// Strip the first match
		return m[1:]
	})
}
