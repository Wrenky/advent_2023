package main

import (
	_ "embed"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//go:embed input
var data string

func init() {
	data = strings.TrimRight(data, "\n")
	if len(data) == 0 {
		panic("No input file")
	}
}

// Digit mapping for part 2.
var mapping = map[string]int{
	"one":   1,
	"two":   2,
	"three": 3,
	"four":  4,
	"five":  5,
	"six":   6,
	"seven": 7,
	"eight": 8,
	"nine":  9,
}

// Summation for parts 1 & 2
func sum(input []int) int {
	var sum int
	for _, v := range input {
		sum += v
	}
	return sum
}

func main() {

	pre, ans, post := time.Now(), sum(parsePart1(data)), time.Now()
	fmt.Printf("Part1 answer: %d, in %s\n", ans, post.Sub(pre))

	pre, ans, post = time.Now(), sum(parsePart2(data)), time.Now()
	fmt.Printf("Part2 answer: %d, in %s\n", ans, post.Sub(pre))

}

// Parse normally, just extract digits
func parsePart1(input string) []int {
	ans := []int{}
	pattern := regexp.MustCompile(`\d`)
	for _, line := range strings.Split(input, "\n") {
		matches := pattern.FindAllString(line, -1)
		if len(matches) == 0 {
			continue
		}
		i, _ := strconv.Atoi(matches[0] + matches[len(matches)-1])
		ans = append(ans, i)
	}
	return ans
}

// Parsing part two is harder.  The problem is with overlapping regular expressions!
//
//	"oneight" should return [1 8], not [1 1]. but normally regualr expressions in go dont allow that
//	so you parse out "one" then are left with "ight". So my approach is to just get each map, and roll
//	down the start index to the match - 1, so match "one" from "oneight" then search "eight" as its only
//	the trailing letter than can overlap.
func parsePart2(input string) []int {

	// This just does the integer converstions
	convert := func(in string) int {
		if v, ok := mapping[in]; ok {
			return v
		}
		i, _ := strconv.Atoi(in) // Cant fail
		return i
	}

	ans := []int{}
	digitRegex := regexp.MustCompile(`\d|one|two|three|four|five|six|seven|eight|nine`)
	// Just run the above on each line!
	for _, line := range strings.Split(input, "\n") {
		fmatch := digitRegex.FindString(line)
		lmatch := fmatch
		lastIdx := 0
		for {
			match := digitRegex.FindString(line[lastIdx:])
			if match == "" {
				break
			}

			lastIdx += len(match)
			lmatch = match
			if len(match) > 1 {
				lastIdx -= 1
			}
		}
		i, _ := strconv.Atoi(fmt.Sprintf("%d%d", convert(fmatch), convert(lmatch)))
		ans = append(ans, i)
	}
	return ans
}
