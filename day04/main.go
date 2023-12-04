package main

import (
	_ "embed"
	"fmt"
	"math"
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
	pre, ans, post := time.Now(), calculateScorecards(parsed), time.Now()
	fmt.Printf("Part1 answer: %d, in %s\n", ans, post.Sub(pre))
	pre, ans, post = time.Now(), countScratchCards(parsed), time.Now()
	fmt.Printf("Part2 answer: %d, in %s\n", ans, post.Sub(pre))
}

// Returns an array of cards [], where each index is two(so another []) []ints (winning numbers, my numbers).
func parseInput(input string) [][][]int {
	return lo.Map(strings.Split(input, "\n"), func(line string, _ int) [][]int {
		// Split each newline by ":", the toss the first half and feed
		//   the second to the next map
		res := strings.Split(line, ":")
		return lo.Map(strings.Split(res[1], "|"), func(cardNumbers string, _ int) []int {
			// now we have our two halves, but split on " " is messy and creates null values.
			// use a filtermap to ignore those and convert the rest to integers.
			return lo.FilterMap(strings.Split(cardNumbers, " "), func(digit string, _ int) (int, bool) {
				if digit == "" {
					return 0, false
				}
				i, _ := strconv.Atoi(digit)
				return i, true
			})
		})
	})
}

// Part 1: Just find out how many matches per card, then just
// run 2^(matches-1). 4 matches, 2^3, etc. Sum the cards.
func calculateScorecards(in [][][]int) int {
	// works on "all" scorecards
	return lo.Reduce(in, func(agg int, line [][]int, _ int) int {
		winners, numbers := line[0], line[1]
		matches := lo.Filter(numbers, func(cand int, _ int) bool {
			return lo.Contains(winners, cand)
		})
		return agg + int(math.Pow(float64(2), float64((len(matches)-1))))
	}, 0)
}

// Part2: Create map of [card] count, and set each to 1.
// Get the amount of winning numbers then update the next
// cards with those new cards * the current card
// Then sum just the map values!
func countScratchCards(in [][][]int) int {
	cardMap := make(map[int]int)
	for i := range in {
		cardMap[i] = 1
	}
	for card, line := range in {
		winners, numbers := line[0], line[1]
		matches := lo.Filter(numbers, func(cand int, _ int) bool {
			return lo.Contains(winners, cand)
		})
		for i := card + 1; i <= (len(matches) + card); i++ {
			cardMap[i] = cardMap[i] + (1 * cardMap[card])
		}
	}

	return lo.Sum(lo.Values(cardMap))
}
