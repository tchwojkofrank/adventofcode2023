package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func readInput(fname string) string {
	content, err := os.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
	}

	// Convert []byte to string
	return string(content)
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Printf("Enter input file name.\n")
		return
	}
	params := os.Args[1]
	inputName := strings.Split(params, " ")[0]
	text := readInput(inputName)
	start := time.Now()
	run(text)
	end := time.Now()
	fmt.Printf("Running time: %v\n", end.Sub(start))
	start = time.Now()
	run2(text)
	end = time.Now()
	fmt.Printf("Running time: %v\n", end.Sub(start))
}

type Pull struct {
	red   int
	green int
	blue  int
}

var CubeCounts = Pull{red: 12, green: 13, blue: 14}

type Game struct {
	pulls []Pull
	ID    int
}

func isPullPossible(pull Pull) bool {
	if pull.red > CubeCounts.red {
		return false
	}
	if pull.green > CubeCounts.green {
		return false
	}
	if pull.blue > CubeCounts.blue {
		return false
	}

	return true
}

func getGameID(idString string) int {
	var id int
	fmt.Sscanf(idString, "Game %d", &id)
	return id
}

// pullString is of the form "%d <color>," repeated up to three times once for each color
// I should have used a map instead of a struct
func getPull(pullString string) Pull {
	pull := Pull{red: 0, green: 0, blue: 0}
	pullString = strings.TrimSpace(pullString)
	parts := strings.Split(pullString, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		fields := strings.Split(part, " ")
		value, _ := strconv.Atoi(fields[0])
		switch fields[1] {
		case "red":
			pull.red = value
		case "green":
			pull.green = value
		case "blue":
			pull.blue = value
		}
	}
	return pull
}

func parseGameString(gameString string) Game {
	var game Game
	game.pulls = make([]Pull, 0)
	parts := strings.Split(gameString, ":")
	game.ID = getGameID(parts[0])
	parts[1] = strings.TrimSpace(parts[1])
	pulls := strings.Split(parts[1], ";")
	for _, pullString := range pulls {
		pull := getPull(pullString)
		game.pulls = append(game.pulls, pull)
	}

	return game
}

func run(input string) string {
	gameStrings := strings.Split(input, "\n")
	games := make([]Game, len(gameStrings))
	for i, gameString := range gameStrings {
		games[i] = parseGameString(gameString)
	}
	possiblePulls := make([]int, 0)
	sum := 0
	for _, game := range games {
		possible := true
		for _, pull := range game.pulls {
			if !isPullPossible(pull) {
				possible = false
				break
			}
		}
		if possible {
			possiblePulls = append(possiblePulls, game.ID)
			sum += game.ID
		}
	}
	fmt.Printf("Possible games: %v\n", possiblePulls)
	fmt.Printf("Sum: %v\n", sum)

	return strconv.Itoa(sum)
}

func run2(input string) string {
	gameStrings := strings.Split(input, "\n")
	games := make([]Game, len(gameStrings))
	for i, gameString := range gameStrings {
		games[i] = parseGameString(gameString)
	}

	minCubes := make([]Pull, len(games))
	powers := make([]int, len(games))
	sum := 0
	for i, game := range games {
		minCubes[i] = Pull{red: 0, green: 0, blue: 0}
		for _, pull := range game.pulls {
			if pull.red > minCubes[i].red {
				minCubes[i].red = pull.red
			}
			if pull.green > minCubes[i].green {
				minCubes[i].green = pull.green
			}
			if pull.blue > minCubes[i].blue {
				minCubes[i].blue = pull.blue
			}
		}
		powers[i] = minCubes[i].red * minCubes[i].green * minCubes[i].blue
		sum += powers[i]
	}

	fmt.Printf("Powers: %v\n", powers)
	fmt.Printf("Sum: %v\n", sum)

	return strconv.Itoa(sum)

}
