package main

import (
	"fmt"
	"log"
	"os"
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

func expandVertically(data map[Point]struct{}, height int, width int, increment int) (map[Point]struct{}, int) {
	newMap := make(map[Point]struct{})
	y2 := 0
	for y := 0; y < height; y++ {
		foundGalaxy := false
		for x := 0; x < width; x++ {
			if _, ok := data[Point{x, y}]; ok {
				newMap[Point{x, y2}] = struct{}{}
				foundGalaxy = true
			}
		}
		if !foundGalaxy {
			y2 = y2 + increment
		}
		y2++
	}
	return newMap, y2
}

func expandHorizontally(data map[Point]struct{}, height int, width int, increment int) (map[Point]struct{}, int) {
	newMap := make(map[Point]struct{})
	x2 := 0
	for x := 0; x < width; x++ {
		foundGalaxy := false
		for y := 0; y < height; y++ {
			if _, ok := data[Point{x, y}]; ok {
				newMap[Point{x2, y}] = struct{}{}
				foundGalaxy = true
			}
		}
		if !foundGalaxy {
			x2 = x2 + increment
		}
		x2++
	}
	return newMap, x2
}

func expandMap(data map[Point]struct{}, height int, width int, increment int) (map[Point]struct{}, int, int) {
	verticalExpansion, newHeight := expandVertically(data, height, width, increment)
	expandedMap, newWidth := expandHorizontally(verticalExpansion, newHeight, width, increment)
	return expandedMap, newHeight, newWidth
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
	height := len(lines)
	width := len(lines[0])
	preExpansion := mapData(lines)
	currentMap, _, _ := expandMap(preExpansion, height, width, 1)
	distanceSum := sumAllDistances(currentMap)
	fmt.Printf("Distance sum: %d\n", distanceSum)
	return fmt.Sprintf("%d", distanceSum)
}

func run2(input string) string {
	lines := strings.Split(input, "\n")
	height := len(lines)
	width := len(lines[0])
	preExpansion := mapData(lines)
	currentMap, _, _ := expandMap(preExpansion, height, width, 1000000-1)
	distanceSum := sumAllDistances(currentMap)
	fmt.Printf("Distance sum: %d\n", distanceSum)
	return fmt.Sprintf("%d", distanceSum)
}
