package main

import (
	"advent/helpers"
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
	data = strings.TrimRight(data, "\n")
	if len(data) == 0 {
		panic("No input file")
	}
}

func main() {

	pre := time.Now()
	ans := lo.Sum(calibrationValues(data))
	post := time.Now()
	fmt.Printf("Part1 answer: %d, in %s\n", ans, post.Sub(pre))

	pre = time.Now()
	ans = lo.Sum(alphaNumCalibrationValues(data))
	post = time.Now()
	fmt.Printf("Part2 answer: %d, in %s\n", ans, post.Sub(pre))
}

// Find the part 1 calibration values
func calibrationValues(input string) []int {
	pattern := regexp.MustCompile(`\d`)
	return lo.Map(strings.Split(input, "\n"), func(line string, _ int) int {
		matches := pattern.FindAllString(line, -1)
		return helpers.Atoi(matches[0] + matches[len(matches)-1])
	})
}

// Find the part 2 calibration values
func alphaNumCalibrationValues(input string) []int {
	digitRegex := regexp.MustCompile(`\d|one|two|three|four|five|six|seven|eight|nine`)
	return lo.Map(strings.Split(input, "\n"), func(line string, _ int) int {
		digits := findAllMatches([]int{}, line, digitRegex)
		return helpers.Atoi(fmt.Sprintf("%d%d", digits[0], digits[len(digits)-1]))
	})
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

func convert(in string) int {
	if v, ok := mapping[in]; ok {
		return v
	}
	return helpers.Atoi(in)
}

// Parsing part two is harder.  The problem is with overlapping regular expressions!
// "oneight" should return [1 8], not [1 1].  This means a regexp alone wont work as
// we need a positive/negative lookback which isnt supported in go- so we find a value,
// move backwards by 1 and rescan forward.
func findAllMatches(agg []int, line string, r *regexp.Regexp) []int {
	match := r.FindString(line)
	var lastIdx int
	if match == "" {
		return agg
	}
	lastIdx += len(match)
	if len(match) > 1 {
		lastIdx -= 1
	}
	return findAllMatches(append(agg, convert(match)), line[lastIdx:], r)
}
