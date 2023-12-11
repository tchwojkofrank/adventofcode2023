package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"golang.org/x/exp/maps"
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

type Point struct {
	x int
	y int
}

func mapData(lines []string) map[Point]struct{} {
	data := make(map[Point]struct{})
	for y, line := range lines {
		for x, char := range line {
			if char == '#' {
				data[Point{x, y}] = struct{}{}
			}
		}
	}
	return data
}

func expandVertically(data map[Point]struct{}, increment int) map[Point]struct{} {
	newMap := make(map[Point]struct{})
	galaxies := maps.Keys(data)
	sort.Slice(galaxies, func(i, j int) bool {
		return galaxies[i].y < galaxies[j].y
	})
	prevY := 0
	y2 := 0
	for _, galaxy := range galaxies {
		if galaxy.y == prevY || galaxy.y == prevY+1 {
			y2 = y2 + (galaxy.y - prevY)
			prevY = galaxy.y
			newMap[Point{galaxy.x, y2}] = struct{}{}
		} else {
			y2 = y2 + increment*(galaxy.y-prevY-1) + (galaxy.y - prevY)
			prevY = galaxy.y
			newMap[Point{galaxy.x, y2}] = struct{}{}
		}
	}
	return newMap
}

func expandHorizontally(data map[Point]struct{}, increment int) map[Point]struct{} {
	newMap := make(map[Point]struct{})
	galaxies := maps.Keys(data)
	sort.Slice(galaxies, func(i, j int) bool {
		return galaxies[i].x < galaxies[j].x
	})
	x2 := 0
	prevX := 0
	for _, galaxy := range galaxies {
		if galaxy.x == prevX || galaxy.x == prevX+1 {
			x2 = x2 + (galaxy.x - prevX)
			prevX = galaxy.x
			newMap[Point{x2, galaxy.y}] = struct{}{}
		} else {
			x2 = x2 + increment*(galaxy.x-prevX-1) + (galaxy.x - prevX)
			prevX = galaxy.x
			newMap[Point{x2, galaxy.y}] = struct{}{}
		}
	}
	return newMap
}

func expandMap(data map[Point]struct{}, increment int) map[Point]struct{} {
	verticalExpansion := expandVertically(data, increment)
	expandedMap := expandHorizontally(verticalExpansion, increment)
	return expandedMap
}

func (p Point) isLess(q Point) bool {
	if p.y < q.y {
		return true
	}
	if p.y == q.y && p.x < q.x {
		return true
	}
	return false
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func sumAllDistances(data map[Point]struct{}) int {
	// sum all the distances between every pair of galaxies
	sum := 0
	for galaxy1 := range data {
		for galaxy2 := range data {
			if galaxy1 != galaxy2 && galaxy1.isLess(galaxy2) {
				// find distance between galaxy1 and galaxy2
				distance := abs(galaxy1.x-galaxy2.x) + abs(galaxy1.y-galaxy2.y)
				sum += distance
			}
		}
	}
	return sum
}

func run(input string) string {
	lines := strings.Split(input, "\n")
	preExpansion := mapData(lines)
	currentMap := expandMap(preExpansion, 1)
	// printMap(currentMap)
	distanceSum := sumAllDistances(currentMap)
	fmt.Printf("Distance sum: %d\n", distanceSum)
	return fmt.Sprintf("%d", distanceSum)
}

const ExpansionRate = 1000000

func run2(input string) string {
	lines := strings.Split(input, "\n")
	preExpansion := mapData(lines)
	currentMap := expandMap(preExpansion, ExpansionRate-1)
	// printMap(currentMap)
	distanceSum := sumAllDistances(currentMap)
	fmt.Printf("Distance sum: %d\n", distanceSum)
	return fmt.Sprintf("%d", distanceSum)
}
