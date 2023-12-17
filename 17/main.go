package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"chwojkofrank.com/astar"
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

const (
	None = iota
	Up
	Right
	Down
	Left
	Any
)

type Point struct {
	x    int
	y    int
	prev [3]int
}

var maxWidth int
var maxHeight int
var cityMap [][]int
var minHeatLoss int

func (p Point) Neighbors() []astar.Node {
	neighbors := make([]astar.Node, 0, 4)

	anyPrev := [3]int{Any, Any, Any}

	prev := [3]int{p.prev[1], p.prev[2], None}
	if p.x > 0 && p.prev[2] != Right &&
		(p.prev[0] != Left || p.prev[1] != Left || p.prev[2] != Left) {
		prev[2] = Left
		neighbors = append(neighbors, astar.Node(Point{p.x - 1, p.y, prev}))
	}
	if p.x < maxWidth-1 && p.prev[2] != Left &&
		(p.prev[0] != Right || p.prev[1] != Right || p.prev[2] != Right) {
		prev[2] = Right
		if p.x == maxWidth-2 && p.y == maxHeight-1 {
			neighbors = append(neighbors, astar.Node(Point{p.x + 1, p.y, anyPrev}))
		} else {
			neighbors = append(neighbors, astar.Node(Point{p.x + 1, p.y, prev}))
		}
	}
	if p.y > 0 && p.prev[2] != Down &&
		(p.prev[0] != Up || p.prev[1] != Up || p.prev[2] != Up) {
		prev[2] = Up
		if p.x == maxWidth-1 && p.y == maxHeight-2 {
			neighbors = append(neighbors, astar.Node(Point{p.x, p.y - 1, anyPrev}))
		} else {
			neighbors = append(neighbors, astar.Node(Point{p.x, p.y - 1, prev}))
		}
	}
	if p.y < maxHeight-1 && p.prev[2] != Up &&
		(p.prev[0] != Down || p.prev[1] != Down || p.prev[2] != Down) {
		prev[2] = Down
		neighbors = append(neighbors, astar.Node(Point{p.x, p.y + 1, prev}))
	}

	return neighbors
}

func (p Point) getDirection(to Point) int {
	if p.x == to.x {
		if p.y > to.y {
			return Up
		} else {
			return Down
		}
	} else {
		if p.x > to.x {
			return Left
		} else {
			return Right
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

func (p Point) Cost(toName string) int {
	to := Point{}
	fmt.Sscanf(toName, "%d,%d", &to.x, &to.y)
	// dir := p.getDirection(to)
	dx := abs(p.x - to.x)
	dy := abs(p.y - to.y)
	if dx*dy != 0 || dx+dy != 1 {
		return 10000
	}

	return cityMap[to.y][to.x]
}

func (p Point) Heuristic(toName string) int {
	to := Point{}
	fmt.Sscanf(toName, "%d,%d", &to.x, &to.y)
	return (abs(p.x-to.x) + abs(p.y-to.y)) * minHeatLoss
}

func (p Point) Name() string {
	return fmt.Sprintf("%d,%d-%d%d%d", p.x, p.y, p.prev[0], p.prev[1], p.prev[2])
}

func getCityMap(city []string) {
	maxHeight = len(city)
	maxWidth = len(city[0])
	cityMap = make([][]int, maxHeight)
	minHeatLoss = 10000
	for y := 0; y < maxHeight; y++ {
		cityMap[y] = make([]int, maxWidth)
		for x := 0; x < maxWidth; x++ {
			cityMap[y][x], _ = strconv.Atoi(string(city[y][x]))
			if cityMap[y][x] < minHeatLoss {
				minHeatLoss = cityMap[y][x]
			}
		}
	}
}

func (p Point) String() string {
	return fmt.Sprintf("[%d, %d] %d\n", p.x, p.y, cityMap[p.y][p.x])
}

func run(input string) string {
	city := strings.Split(input, "\n")
	getCityMap(city)
	startBlock := Point{0, 0, [3]int{None, None, None}}
	endBlock := Point{maxWidth - 1, maxHeight - 1, [3]int{Any, Any, Any}}
	path := astar.Astar(startBlock, endBlock)
	fmt.Printf("Path:\n%v\n", path)
	totalCost := 0
	for i := 1; i < len(path); i++ {
		point := path[i].(Point)
		totalCost += cityMap[point.y][point.x]
	}
	fmt.Printf("Total cost: %d\n", totalCost)
	return fmt.Sprintf("%d", totalCost)
}

func run2(input string) string {
	return ""
}
