package main

import (
	"advent/helpers"
	_ "embed"
	"fmt"
	"regexp"
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

	// Get seeds, generate the maps and make a chain
	seeds := getSeeds(data)
	s2s := mapGenerator(data, regexp.MustCompile(`seed-to-soil\smap:\n((?:\d+\s+\d+\s+\d+\n)+)`))
	s2f := mapGenerator(data, regexp.MustCompile(`soil-to-fertilizer\smap:\n((?:\d+\s+\d+\s+\d+\n)+)`))
	f2w := mapGenerator(data, regexp.MustCompile(`fertilizer-to-water\smap:\n((?:\d+\s+\d+\s+\d+\n)+)`))
	w2l := mapGenerator(data, regexp.MustCompile(`water-to-light\smap:\n((?:\d+\s+\d+\s+\d+\n)+)`))
	l2t := mapGenerator(data, regexp.MustCompile(`light-to-temperature\smap:\n((?:\d+\s+\d+\s+\d+\n)+)`))
	t2h := mapGenerator(data, regexp.MustCompile(`temperature-to-humidity\smap:\n((?:\d+\s+\d+\s+\d+\n)+)`))
	h2l := mapGenerator(data, regexp.MustCompile(`humidity-to-location\smap:\n((?:\d+\s+\d+\s+\d+\n)+)`))
	chain := func(seed int) int {
		return h2l(t2h(l2t(w2l(f2w(s2f(s2s(seed)))))))
	}
	pre := time.Now()
	p1 := lo.Map(seeds, func(seed int, _ int) int {
		return chain(seed)
	})
	post := time.Now()
	fmt.Printf("part1 answer: %d in %s\n", lo.Min(p1), post.Sub(pre))

	pre = time.Now()
	ans := lo.Min(lop.Map(lo.Chunk(seeds, 2), func(r []int, _ int) int {
		answers := []int{}
		for i := r[0]; i < r[0]+r[1]; i++ {
			answers = append(answers, chain(i))
		}
		return lo.Min(answers)
	}))
	post = time.Now()
	fmt.Printf("part2 answer: %d in %s\n", ans, post.Sub(pre))
}

// Gather seed values
func getSeeds(input string) []int {
	pattern := regexp.MustCompile(`seeds:\s(.*?)\n`)
	m := pattern.FindStringSubmatch(input)
	return lo.FilterMap(strings.Split(m[1], " "), func(digit string, _ int) (int, bool) {
		return helpers.Atoi(digit), (digit != "")
	})
}

// Generate our mapper functions from the input. Return a int->int function
func mapGenerator(input string, pattern *regexp.Regexp) func(int) int {
	m := pattern.FindStringSubmatch(input)
	cleaned := lo.FilterMap(strings.Split(m[1], "\n"), func(line string, _ int) ([]int, bool) {
		if line == "" {
			return []int{}, false
		}
		return lo.FilterMap(strings.Split(line, " "), func(digit string, _ int) (int, bool) {
			return helpers.Atoi(digit), (digit != "")
		}), true
	})
	return func(i int) int {
		for _, v := range cleaned {
			dest, source, rang := v[0], v[1], v[2]
			if (i >= source) && (i < (source + rang)) {
				return dest + (i - source)
			}
		}
		return i
	}
}
