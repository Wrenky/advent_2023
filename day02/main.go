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
	// Strip trailing newline
	data = strings.TrimRight(data, "\n")
	if len(data) == 0 {
		panic("No input file")
	}
}

func main() {
	pre, games, post := time.Now(), parseInput(data), time.Now()
	fmt.Printf("Data parsed in %s\n", post.Sub(pre))

	pre, ans, post := time.Now(), p1(games), time.Now()
	fmt.Printf("Part1 answer: %d, in %s\n", ans, post.Sub(pre))

	pre, ans, post = time.Now(), p2(games), time.Now()
	fmt.Printf("Part2 answer: %d, in %s\n", ans, post.Sub(pre))

}

// Handle the nasty parsing, put everything into []Game
func parseInput(input string) []Game {
	games := []Game{}
	gamePattern := regexp.MustCompile(`Game\s(\d+):(.*)`)
	setPattern := regexp.MustCompile(`(\d+)\s(blue|red|green)`)
	for _, line := range strings.Split(input, "\n") {
		gameMatch := gamePattern.FindStringSubmatch(line)
		gameId, sets := gameMatch[1], gameMatch[2]
		i, _ := strconv.Atoi(gameId)
		newGame := Game{Id: i}
		for _, set := range strings.Split(sets, ";") {
			batches := setPattern.FindAllStringSubmatch(set, -1)
			for _, batch := range batches {
				amount, _ := strconv.Atoi(batch[1])
				newGame.Sets = append(newGame.Sets, Set{Amount: amount, Color: batch[2]})
			}
		}
		games = append(games, newGame)
	}
	return games
}

type Game struct {
	Id   int
	Sets []Set
}
type Set struct {
	Amount int
	Color  string
}

func p1(games []Game) int {
	answer := 0
	colorMax := map[string]int{"red": 12, "green": 13, "blue": 14}
	for _, game := range games {
		valid := true
		for _, set := range game.Sets {
			if set.Amount > colorMax[set.Color] {
				valid = false
			}
		}
		if valid {
			answer += game.Id
		}
	}
	return answer
}
func p2(games []Game) int {
	answer := 0
	for _, game := range games {
		lowestRequired := map[string]int{"red": -1, "green": -1, "blue": -1}
		for _, set := range game.Sets {
			lowestRequired[set.Color] = max(set.Amount, lowestRequired[set.Color])
		}
		power := 1
		for _, v := range lowestRequired {
			power *= v
		}
		answer += power
	}
	return answer
}
