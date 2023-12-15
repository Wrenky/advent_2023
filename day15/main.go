package main

import (
	"advent/helpers"
	_ "embed"
	"fmt"
	"regexp"
	"strings"

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
	input := parseInput(data)
	p1 := lo.Sum(lo.Map(input, func(v []rune, _ int) int {
		return runHASH(v)
	}))
	fmt.Printf("Part1: %d\n", p1)

	setPattern := regexp.MustCompile(`(.*?)(=|-)(\d)?`)
	b := newBox()
	for _, v := range input {
		matches := setPattern.FindAllStringSubmatch(string(v), -1)
		label, op, fl := matches[0][1], ([]rune(matches[0][2]))[0], matches[0][3]
		switch op {
		case '=':
			b.insert(lens{label, helpers.Atoi(fl)})
		case '-':
			b.remove(lens{label, 0})
		}
	}

	sum := 0
	for bn, c := range b.box {
		for sn, l := range c.content {
			sum += (bn + 1) * (sn + 1) * l.fl
		}
	}
	fmt.Printf("Part2: %d\n", sum)
}

// ----------------------------------------------------------
type boxes struct {
	box map[int]contents
}

type contents struct {
	labels  map[string]int
	content []lens
}

type lens struct {
	label string
	fl    int
}

func newBox() boxes {
	b := boxes{}
	b.box = make(map[int](contents))
	return b
}

func (b *boxes) remove(l lens) {
	index := runHASH([]rune(l.label))
	if c, ok := b.box[index]; ok {
		if _, ok := c.labels[l.label]; ok {
			// Its here! Remove it.
			var i int
			for i = 0; i < len(c.content); i++ {
				if c.content[i].label != l.label {
					// Found it.
					continue
				}
				s := b.box[index]
				s.content = remove(s.content, i)
				b.box[index] = s
				break
			}
			delete(b.box[index].labels, l.label)
		}
		if len(b.box[index].content) == 0 {
			delete(b.box, index)
		}
	}
}

func remove(slice []lens, s int) []lens {
	return append(slice[:s], slice[s+1:]...)
}

func (b *boxes) insert(l lens) {
	index := runHASH([]rune(l.label))
	// Find box
	if c, ok := b.box[index]; ok {
		// Box exists!
		// Check if label exists
		if _, ok := c.labels[l.label]; ok {
			// Label exists! Find and update it
			for i := 0; i < len(c.content); i++ {
				if c.content[i].label == l.label {
					// Update it to my lens value
					c.content[i].fl = l.fl
					bc := b.box[index]
					bc.content = c.content
					b.box[index] = bc
					break
				}
			}
		} else {
			// No label here, add to end
			c.content = append(c.content, l)
			bc := b.box[index]
			bc.content = c.content
			bc.labels[l.label] = 0
			b.box[index] = bc

		}
	} else {
		// new box!
		c := contents{}
		c.labels = make(map[string]int)
		c.labels[l.label] = index
		c.content = append([]lens{}, l)
		b.box[index] = c
	}
}

//----------------------------------------------------------

func runHASH(in []rune) int {
	cv := 0
	for _, r := range in {
		cv = HASH(cv, r)
	}
	return cv
}
func HASH(start int, c rune) int {
	return (((start + int(c)) * 17) % 256)
}

// This needs to change to match your actual input
func parseInput(input string) [][]rune {
	return lo.Map(strings.Split(input, ","), func(line string, _ int) []rune {
		return []rune(line)
	})
}
