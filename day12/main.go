package main

import (
	"advent/helpers"
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
	records := parseInput(data)
	pre := time.Now()
	p1 := lo.Sum(lo.Map(records, func(r record, i int) int { return Find(r) }))
	post := time.Now()
	fmt.Printf("Part1: %d in %s\n", p1, post.Sub(pre))

	pre = time.Now()
	expanded := lo.Map(records, func(r record, _ int) record { return expandRecord(r, 5) })
	p2 := lo.Sum(lop.Map(expanded, func(r record, i int) int { return Find(r) }))
	post = time.Now()
	fmt.Printf("Part2: %d in %s\n", p2, post.Sub(pre))
}

func expandRecord(r record, amount int) record {
	newLine := r.line
	for i := 0; i < amount-1; i++ {
		newLine = append(newLine, '?')
		newLine = append(newLine, r.line...)
	}
	newSprings := r.springs
	for i := 0; i < amount-1; i++ {
		newSprings = append(newSprings, r.springs...)
	}
	return record{
		line:    newLine,
		springs: newSprings,
	}
}

type record struct {
	line    []rune
	springs []int
}

// Memoization
type Memo map[Key]int
type Key struct {
	i          int
	si         int
	currLength int
}

// Wrap find to make the main less scary
func Find(r record) int {
	m := Memo(make(map[Key]int))
	return find(r.line, r.springs, &m, 0, 0, 0)
}

// line: The line we are on
// springs: the set of springs we need to find matches on
// m: Memoize the values
// i: poisition in line
// si: position in springs
// currLength: length of the current block we are in
func find(line []rune, springs []int, m *Memo, i int, si int, currLength int) int {

	// Memoize results
	k := Key{i: i, si: si, currLength: currLength}
	if v, ok := (*m)[k]; ok {
		return v
	}

	if i == len(line) {
		// Done with the line!
		if (si == len(springs) && currLength == 0) ||
			(si == len(springs)-1 && (springs[si] == currLength)) {
			// Done with all the springs, and not in a current spring set, OR
			// Done with all the springs, and the current spring set == currLeng
			return 1
		}
		// Invalid, so return
		return 0
	}

	res := 0
	for _, v := range []rune{'.', '#'} {
		if line[i] == v || line[i] == '?' {
			// We are matched, or a ?
			switch {
			case v == '.' && currLength == 0:
				// Its a dot, and we dont have a current block (dot to dot):
				// Continue down the line
				res += find(line, springs, m, i+1, si, 0)
			case v == '.' && currLength > 0 && si < len(springs) && springs[si] == currLength:
				// Its a dot, but we had a running block and it perfectly matches our spring length
				// move up the springs, resent blockLength
				res += find(line, springs, m, i+1, si+1, 0)
			case v == '#':
				// Its another block. Continue down the line
				res += find(line, springs, m, i+1, si, currLength+1)
			}
		}
	}

	(*m)[k] = res
	return res
}

// This needs to change to match your actual input
func parseInput(input string) []record {
	return lo.Map(strings.Split(input, "\n"), func(line string, _ int) record {
		l := strings.Split(line, " ")
		return record{
			line: []rune(l[0]),
			springs: lo.Map(strings.Split(l[1], ","), func(d string, _ int) int {
				return helpers.Atoi(d)
			}),
		}
	})
}
