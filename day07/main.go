package main

import (
	_ "embed"
	"fmt"
	"sort"
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
	// Just a struct to run parts. Kind of like table testing
	type parts struct {
		name  string
		freq  func(string) map[rune]int
		ranks string
	}
	runs := []parts{{"part1", cardFreq, "23456789TJQKA"}, {"part2", cardFreqJokers, "J23456789TQKA"}}
	for _, v := range runs {
		// Parse the structure based on the frequency function. Returns sorted.
		pre, parsed := time.Now(), parseInput(data, v.freq, v.ranks)

		// Now just sum it up
		ans, post := lo.Reduce(parsed, func(agg int, h hand, i int) int {
			return agg + (h.bid * (i + 1))
		}, 0), time.Now()

		fmt.Printf("%s answer: %d, in %s\n", v.name, ans, post.Sub(pre))
	}
}

type hand struct {
	hand     string
	bid      int
	strength int // 7 strengths, 7 is 5 of kind, while 1 is high card
}

// Less functions for hands
func handLess(i hand, j hand, ranks string) bool {
	if i.strength != j.strength {
		return i.strength < j.strength
	}
	// Equal case!
	for k := 0; k < len(i.hand); k++ {
		if i.hand[k] != j.hand[k] {
			iRank := strings.IndexRune(ranks, rune(i.hand[k]))
			jRank := strings.IndexRune(ranks, rune(j.hand[k]))
			return iRank < jRank
		}
	}
	// Shouldnt happen!
	return false
}

func parseInput(input string, counter func(string) map[rune]int, ranks string) []hand {
	res := lo.Map(strings.Split(input, "\n"), func(line string, _ int) hand {
		hb := strings.Split(line, " ")
		bid, _ := strconv.Atoi(hb[1])
		return hand{
			hand:     hb[0],
			bid:      bid,
			strength: getStrength(hb[0], counter),
		}
	})
	sort.Slice(res, func(i, j int) bool {
		return handLess(res[i], res[j], ranks)
	})
	return res
}

// Simple frequencys of the hand.
func cardFreq(hand string) map[rune]int {
	counts := make(map[rune]int)
	for _, c := range hand {
		counts[c]++
	}
	return counts
}

// Card frequencies but we ignore jokers and add them in the
// end to the highest frequency.
func cardFreqJokers(hand string) map[rune]int {
	counts := make(map[rune]int)
	var jokers int
	var maxR rune

	for _, c := range hand {
		if c == 'J' {
			jokers++
			continue
		}
		counts[c]++
		if counts[c] >= counts[maxR] {
			maxR = c
		}
	}

	counts[maxR] += jokers
	return counts
}

// Fetchs the frequencies for the cards, then figures out the hand strength.
func getStrength(hand string, freq func(string) map[rune]int) int {
	counts := freq(hand)
	switch len(counts) {
	case 1:
		return 7 // 5 of a kind
	case 2:
		if lo.Contains(lo.Values(counts), 4) {
			return 6 // 4 of a kind
		}
		return 5 // Full house
	case 3:
		if lo.Contains(lo.Values(counts), 3) {
			return 4 // Three of a kind
		}
		return 3 // Two pairs + extra
	case 4:
		return 2 // Single pair
	default:
		return 1 // High card
	}
}
