package main

import (
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
	times, distances := parseInput(data)
	tdCombined := lo.Interleave(times, distances)

	// Part1: find out how many wins between each time,distance pair and multiply them
	pre := time.Now()
	p1 := lo.Reduce(lo.Chunk(tdCombined, 2), func(agg int, td []int, _ int) int {
		return agg * amountToWinRace(td[0], td[1])
	}, 1)
	post := time.Now()
	fmt.Printf("Part1 answer: %d in %s\n", p1, post.Sub(pre))

	// Part2: Join all the times and distances and figure out how many wins
	pre = time.Now()
	p2 := amountToWinRace(intJoin(times), intJoin(distances))
	post = time.Now()
	fmt.Printf("Part2 answer: %d in %s\n", p2, post.Sub(pre))
}

// convert to string, then convert to int
func intJoin(nums []int) int {
	i, _ := strconv.Atoi(
		lo.Reduce(nums, func(agg string, n int, _ int) string {
			return agg + strconv.Itoa(n)
		}, ""),
	)
	return i
}

// The actual check- Just find when we get above distance (the "inflection"), then
// just subtract that from both ends of distance (its symmetrical)
func amountToWinRace(t, d int) int {
	i, j, inflection := 1, t-1, 0
	for (i <= j) && (i*j <= d) {
		inflection = i
		i, j = i+1, j-1
	}
	return t - 1 - (2 * inflection)
}

func parseInput(input string) ([]int, []int) {
	pattern := regexp.MustCompile(`Time:\s+(.*?)\nDistance:\s+(.*?)$`)
	matches := pattern.FindStringSubmatch(input)
	res := lo.Map(matches[1:], func(in string, _ int) []int {
		return lo.FilterMap(strings.Split(in, " "), func(d string, _ int) (int, bool) {
			i, _ := strconv.Atoi(d)
			return i, (d != "")
		})
	})
	return res[0], res[1]
}
