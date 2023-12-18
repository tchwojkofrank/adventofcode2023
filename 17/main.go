package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"chwojkofrank.com/astar"
	"chwojkofrank.com/cursor"
	"chwojkofrank.com/dijkstra"
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
	None  = 0
	Up    = 1
	Right = 2
	Down  = 3
	Left  = 4
	Any   = 5
)

type Point struct {
	x             int
	y             int
	lastDirection int
	count         int
}

var maxWidth int
var maxHeight int
var cityMap [][]int
var minHeatLoss int

func (p Point) DNeighbors() ([]dijkstra.Node, []int) {
	neighbors := make([]dijkstra.Node, 0, 4)
	newDirection := None
	newCount := 0
	costs := make([]int, 0, 4)
	// check Up
	if p.y > 0 && p.lastDirection != Down &&
		(p.lastDirection != Up || p.count < 3) {
		newDirection = Up
		newCount = 1
		if p.lastDirection == Up {
			newCount = p.count + 1
		}
		neighbors = append(neighbors, dijkstra.Node(Point{p.x, p.y - 1, newDirection, newCount}))
		costs = append(costs, cityMap[p.y-1][p.x])
	}
	// check Down
	if p.y < maxHeight-1 && p.lastDirection != Up &&
		(p.lastDirection != Down || p.count < 3) {
		newDirection = Down
		newCount = 1
		if p.lastDirection == Down {
			newCount = p.count + 1
		}
		// check if we are at the end
		if p.x == maxWidth-1 && p.y == maxHeight-2 {
			neighbors = append(neighbors, dijkstra.Node(Point{p.x, p.y - 1, Any, 0}))
		} else {
			neighbors = append(neighbors, dijkstra.Node(Point{p.x, p.y + 1, newDirection, newCount}))
		}
		costs = append(costs, cityMap[p.y+1][p.x])
	}
	// check Left
	if p.x > 0 && p.lastDirection != Right &&
		(p.lastDirection != Left || p.count < 3) {
		newDirection = Left
		newCount = 1
		if p.lastDirection == Left {
			newCount = p.count + 1
		}
		neighbors = append(neighbors, dijkstra.Node(Point{p.x - 1, p.y, newDirection, newCount}))
		costs = append(costs, cityMap[p.y][p.x-1])
	}
	// check Right
	if p.x < maxWidth-1 && p.lastDirection != Left &&
		(p.lastDirection != Right || p.count < 3) {
		newDirection = Right
		newCount = 1
		if p.lastDirection == Right {
			newCount = p.count + 1
		}
		// check if we are at the end
		if p.x == maxWidth-2 && p.y == maxHeight-1 {
			neighbors = append(neighbors, dijkstra.Node(Point{p.x + 1, p.y, Any, 0}))
		} else {
			neighbors = append(neighbors, dijkstra.Node(Point{p.x + 1, p.y, newDirection, newCount}))
		}
		costs = append(costs, cityMap[p.y][p.x+1])
	}

	return neighbors, costs
}

func (p Point) Neighbors() []astar.Node {
	neighbors := make([]astar.Node, 0, 4)
	newDirection := None
	var newCount int

	// constraints:
	// 1. we can't go back the way we came
	// 2. we can't go straight more than 3 times in a row
	// 3. we can't go off the edge of the map

	// We can go left if:
	// 1. we are not on the left edge
	// 2. we haven't gone right last
	// 3. we haven't gone left 3 times in a row
	if p.x > 0 &&
		p.lastDirection != Right &&
		(p.lastDirection != Left || p.count < 3) {
		newDirection = Left
		if p.lastDirection == Left {
			newCount = p.count + 1
		} else {
			newCount = 1
		}
		// add the Left neighbor with the new count
		neighbors = append(neighbors, astar.Node(Point{p.x - 1, p.y, newDirection, newCount}))
	}

	// We can go right if:
	// 1. we are not on the right edge
	// 2. we haven't gone left last
	// 3. we haven't gone right 3 times in a row
	if p.x < maxWidth-1 &&
		p.lastDirection != Left &&
		(p.lastDirection != Right || p.count < 3) {
		newDirection = Right
		if p.lastDirection == Right {
			newCount = p.count + 1
		} else {
			newCount = 1
		}
		// check if we're adding the end node
		if p.x == maxWidth-2 && p.y == maxHeight-1 {
			// if this is the end node, then we don't need to keep track of the direction or count
			neighbors = append(neighbors, astar.Node(Point{p.x + 1, p.y, Any, 0}))
		} else {
			neighbors = append(neighbors, astar.Node(Point{p.x + 1, p.y, newDirection, newCount}))
		}
	}

	// We can go up if:
	// 1. we are not on the top edge
	// 2. we haven't gone down last
	// 3. we haven't gone up 3 times in a row
	if p.y > 0 &&
		p.lastDirection != Down &&
		(p.lastDirection != Up || p.count < 3) {
		newDirection = Up
		if p.lastDirection == Up {
			newCount = p.count + 1
		} else {
			newCount = 1
		}
		neighbors = append(neighbors, astar.Node(Point{p.x, p.y - 1, newDirection, newCount}))
	}

	// We can go down if:
	// 1. we are not on the bottom edge
	// 2. we haven't gone up last
	// 3. we haven't gone down 3 times in a row
	if p.y < maxHeight-1 && p.lastDirection != Up &&
		(p.lastDirection != Down || p.count < 3) {
		newDirection = Down
		if p.lastDirection == Down {
			newCount = p.count + 1
		} else {
			newCount = 1
		}
		if p.x == maxWidth-1 && p.y == maxHeight-2 {
			neighbors = append(neighbors, astar.Node(Point{p.x, p.y + 1, Any, 0}))
		} else {
			neighbors = append(neighbors, astar.Node(Point{p.x, p.y + 1, newDirection, newCount}))
		}
	}

	return neighbors
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
	return fmt.Sprintf("%d,%d-%d%d", p.x, p.y, p.lastDirection, p.count)
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
	startBlock := Point{0, 0, None, 0}
	endBlock := Point{maxWidth - 1, maxHeight - 1, Any, 0}
	path := astar.Astar(startBlock, endBlock)
	totalCost := 0
	for i := 1; i < len(path); i++ {
		point := path[i].(Point)
		totalCost += cityMap[point.y][point.x]
		fmt.Printf("%v\t%6d\n", point, totalCost)
	}

	fmt.Printf("A* Total cost: %d\n", totalCost)

	return fmt.Sprintf("%d", totalCost)
}

type Point2 Point

func (p Point2) Name() string {
	return Point(p).Name()
}

func (p Point2) Neighbors() []astar.Node {
	neighbors := make([]astar.Node, 0, 4)
	newDirection := None
	var newCount int

	// constraints:
	// 1. we can't go back the way we came
	// 2. we must go straight at least 4 times
	// 2A. we must turn after 10 straight moves
	// 3. we can't go off the edge of the map

	// We can go left if:
	// 1. we are not on the left edge
	// 2. we haven't gone right last
	// 3. we weren't going left but we have gone straight 4 times in a row
	// 4. we haven't gone left 10 times in a row
	if p.x > 0 &&
		p.lastDirection != Right &&
		((p.lastDirection == Left && p.count < 10) ||
			(p.lastDirection != Left && p.count >= 4)) {
		newDirection = Left
		if p.lastDirection == Left {
			newCount = p.count + 1
		} else {
			newCount = 1
		}
		// add the Left neighbor with the new count
		neighbors = append(neighbors, astar.Node(Point2{p.x - 1, p.y, newDirection, newCount}))
	}

	// We can go right if:
	// 0. we are just starting
	// 1. we are not on the right edge
	// 2. we haven't gone left last
	// 3. we have gone straight 4 times in a row
	// 4. we haven't gone right 10 times in a row
	if (p.lastDirection == None) || (p.x < maxWidth-1 &&
		p.lastDirection != Left &&
		((p.lastDirection == Right && p.count < 10) ||
			(p.lastDirection != Right && p.count >= 4))) {
		newDirection = Right
		if p.lastDirection == Right {
			newCount = p.count + 1
		} else {
			newCount = 1
		}
		// check if we're adding the end node
		// we must have gone straight at least 4 times
		if p.x == maxWidth-2 && p.y == maxHeight-1 && p.count >= 4 {
			// if this is the end node, then we don't need to keep track of the direction or count
			neighbors = append(neighbors, astar.Node(Point2{p.x + 1, p.y, Any, 0}))
		} else {
			neighbors = append(neighbors, astar.Node(Point2{p.x + 1, p.y, newDirection, newCount}))
		}
	}

	// We can go up if:
	// 1. we are not on the top edge
	// 2. we haven't gone down last
	// 3. we have gone straight 4 times in a row
	// 4. we haven't gone up 10 times in a row
	if p.y > 0 &&
		p.lastDirection != Down &&
		((p.lastDirection == Up && p.count < 10) ||
			(p.lastDirection != Up && p.count >= 4)) {
		newDirection = Up
		if p.lastDirection == Up {
			newCount = p.count + 1
		} else {
			newCount = 1
		}
		neighbors = append(neighbors, astar.Node(Point2{p.x, p.y - 1, newDirection, newCount}))
	}

	// We can go down if:
	// 0. we are just starting
	// 1. we are not on the bottom edge
	// 2. we haven't gone up last
	// 3. we have gone straight 4 times in a row
	// 4. we haven't gone down 10 times in a row
	if (p.lastDirection == None) || (p.y < maxHeight-1 && p.lastDirection != Up &&
		((p.lastDirection == Down && p.count < 10) ||
			(p.lastDirection != Down && p.count >= 4))) {
		newDirection = Down
		if p.lastDirection == Down {
			newCount = p.count + 1
		} else {
			newCount = 1
		}
		// check if we are at the end
		// we must have gone straight at least 4 times
		if p.x == maxWidth-1 && p.y == maxHeight-2 && p.count >= 4 {
			neighbors = append(neighbors, astar.Node(Point2{p.x, p.y + 1, Any, 0}))
		} else {
			neighbors = append(neighbors, astar.Node(Point2{p.x, p.y + 1, newDirection, newCount}))
		}
	}

	return neighbors
}

func (p Point2) Cost(name string) int {
	return Point(p).Cost(name)
}

func (p Point2) Heuristic(target string) int {
	return Point(p).Heuristic(target)
}

func run2(input string) string {
	city := strings.Split(input, "\n")
	getCityMap(city)
	startBlock := Point2{0, 0, None, 0}
	endBlock := Point2{maxWidth - 1, maxHeight - 1, Any, 0}
	path := astar.Astar(startBlock, endBlock)
	totalCost := 0
	cursor.Clear()
	for i := 1; i < len(path); i++ {
		point := path[i].(Point2)
		totalCost += cityMap[point.y][point.x]
	}

	fmt.Printf("A* Total cost: %d\n", totalCost)

	return fmt.Sprintf("%d", totalCost)
}
