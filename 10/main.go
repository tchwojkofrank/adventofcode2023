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

type NeighborMap map[Point][]Point
type Map map[Point]rune
type Distances map[Point]int

const (
	Up = iota
	Right
	Down
	Left
)

func direction(delta Point) int {
	if delta.x == 0 && delta.y == -1 {
		return Up
	}
	if delta.x == 1 && delta.y == 0 {
		return Right
	}
	if delta.x == 0 && delta.y == 1 {
		return Down
	}
	if delta.x == -1 && delta.y == 0 {
		return Left
	}
	return -1
}

func directionToPoint(direction int) Point {
	switch direction {
	case Up:
		return Point{0, -1}
	case Right:
		return Point{1, 0}
	case Down:
		return Point{0, 1}
	case Left:
		return Point{-1, 0}
	}
	return Point{0, 0}
}

// returns true if the two points are connected
// assumes a delta of 1
func isConnected(m Map, p Point, delta Point) bool {
	p2 := Point{p.x + delta.x, p.y + delta.y}
	switch direction(delta) {
	case Up:
		switch m[p] {
		case 'L', '|', 'J':
			if m[p2] == 'F' || m[p2] == '|' || m[p2] == '7' || m[p2] == 'S' {
				return true
			} else {
				return false
			}
		default:
			return false
		}
	case Right:
		switch m[p] {
		case 'F', '-', 'L':
			if m[p2] == 'J' || m[p2] == '-' || m[p2] == '7' || m[p2] == 'S' {
				return true
			} else {
				return false
			}
		default:
			return false
		}
	case Down:
		switch m[p] {
		case '7', '|', 'F':
			if m[p2] == 'L' || m[p2] == '|' || m[p2] == 'J' || m[p2] == 'S' {
				return true
			} else {
				return false
			}
		default:
			return false
		}
	case Left:
		switch m[p] {
		case 'J', '-', '7':
			if m[p2] == 'F' || m[p2] == '-' || m[p2] == 'L' || m[p2] == 'S' {
				return true
			} else {
				return false
			}
		default:
			return false
		}
	}
	return false
}

func makeMaps(input string) (NeighborMap, Map, Point, rune) {
	m := make(Map)
	neighborMap := make(NeighborMap)
	lines := strings.Split(input, "\n")
	var animal Point
	for y := 0; y < len(lines); y++ {
		for x := 0; x < len(lines[y]); x++ {
			neighborMap[Point{x, y}] = make([]Point, 0, 4)
		}
	}
	var animalRune rune
	for y, line := range lines {
		for x, char := range line {
			m[Point{x, y}] = char
		}
	}
	for y, line := range lines {
		for x, char := range line {
			p := Point{x, y}
			// m[p] = char
			switch char {
			case '|':
				if y > 0 {
					delta := Point{0, -1}
					if isConnected(m, p, delta) {
						neighborMap[p] = append(neighborMap[p], Point{x, y - 1})
					}
				}
				if y < len(lines)-1 {
					delta := Point{0, 1}
					if isConnected(m, p, delta) {
						neighborMap[p] = append(neighborMap[p], Point{x, y + 1})
					}
				}
			case '-':
				if x > 0 {
					delta := Point{-1, 0}
					if isConnected(m, p, delta) {
						neighborMap[p] = append(neighborMap[p], Point{x - 1, y})
					}
				}
				if x < len(line)-1 {
					delta := Point{1, 0}
					if isConnected(m, p, delta) {
						neighborMap[p] = append(neighborMap[p], Point{x + 1, y})
					}
				}
			case 'L':
				if y > 0 {
					delta := Point{0, -1}
					if isConnected(m, p, delta) {
						neighborMap[p] = append(neighborMap[p], Point{x, y - 1})
					}
				}
				if x < len(line)-1 {
					delta := Point{1, 0}
					if isConnected(m, p, delta) {
						neighborMap[p] = append(neighborMap[p], Point{x + 1, y})
					}
				}
			case 'J':
				if y > 0 {
					delta := Point{0, -1}
					if isConnected(m, p, delta) {
						neighborMap[p] = append(neighborMap[p], Point{x, y - 1})
					}
				}
				if x > 0 {
					delta := Point{-1, 0}
					if isConnected(m, p, delta) {
						neighborMap[p] = append(neighborMap[p], Point{x - 1, y})
					}
				}
			case '7':
				if y < len(lines)-1 {
					delta := Point{0, 1}
					if isConnected(m, p, delta) {
						neighborMap[p] = append(neighborMap[p], Point{x, y + 1})
					}
				}
				if x > 0 {
					delta := Point{-1, 0}
					if isConnected(m, p, delta) {
						neighborMap[p] = append(neighborMap[p], Point{x - 1, y})
					}
				}
			case 'F':
				if y < len(lines)-1 {
					delta := Point{0, 1}
					if isConnected(m, p, delta) {
						neighborMap[p] = append(neighborMap[p], Point{x, y + 1})
					}
				}
				if x < len(line)-1 {
					delta := Point{1, 0}
					if isConnected(m, p, delta) {
						neighborMap[p] = append(neighborMap[p], Point{x + 1, y})
					}
				}
			case '.':
			case 'S':
				animal = p
				connectedDirections := make(map[int]bool)
				for dy := -1; dy <= 1; dy++ {
					for dx := -1; dx <= 1; dx++ {
						delta := Point{dx, dy}
						if dx*dy != 0 || (dx == 0 && dy == 0) {
							continue
						}
						dp := Point{p.x + dx, p.y + dy}
						direction := direction(delta)
						switch direction {
						case Up:
							// then the m[dp] should be facing down
							if m[dp] == 'F' || m[dp] == '|' || m[dp] == '7' {
								connectedDirections[direction] = true
							}
						case Right:
							// then the m[dp] should be facing left
							if m[dp] == 'J' || m[dp] == '-' || m[dp] == '7' {
								connectedDirections[direction] = true
							}
						case Down:
							// then the m[dp] should be facing up
							if m[dp] == 'L' || m[dp] == '|' || m[dp] == 'J' {
								connectedDirections[direction] = true
							}
						case Left:
							// then the m[dp] should be facing right
							if m[dp] == 'F' || m[dp] == '-' || m[dp] == 'L' {
								connectedDirections[direction] = true
							}
						}
					}
				}
				if len(connectedDirections) == 0 {
					log.Fatal("No connected directions")
				} else if len(connectedDirections) != 2 {
					log.Fatal("Animal is not on a loop")
				}
				if connectedDirections[Up] && connectedDirections[Left] {
					animalRune = 'J'
					upDelta := directionToPoint(Up)
					leftDelta := directionToPoint(Left)
					upPoint := Point{p.x + upDelta.x, p.y + upDelta.y}
					leftPoint := Point{p.x + leftDelta.x, p.y + leftDelta.y}
					neighborMap[p] = append(neighborMap[p], upPoint)
					neighborMap[p] = append(neighborMap[p], leftPoint)
				} else if connectedDirections[Up] && connectedDirections[Right] {
					animalRune = 'L'
					upDelta := directionToPoint(Up)
					rightDelta := directionToPoint(Right)
					upPoint := Point{p.x + upDelta.x, p.y + upDelta.y}
					rightPoint := Point{p.x + rightDelta.x, p.y + rightDelta.y}
					neighborMap[p] = append(neighborMap[p], upPoint)
					neighborMap[p] = append(neighborMap[p], rightPoint)
				} else if connectedDirections[Down] && connectedDirections[Left] {
					animalRune = '7'
					downDelta := directionToPoint(Down)
					leftDelta := directionToPoint(Left)
					downPoint := Point{p.x + downDelta.x, p.y + downDelta.y}
					leftPoint := Point{p.x + leftDelta.x, p.y + leftDelta.y}
					neighborMap[p] = append(neighborMap[p], downPoint)
					neighborMap[p] = append(neighborMap[p], leftPoint)
				} else if connectedDirections[Down] && connectedDirections[Right] {
					animalRune = 'F'
					downDelta := directionToPoint(Down)
					rightDelta := directionToPoint(Right)
					downPoint := Point{p.x + downDelta.x, p.y + downDelta.y}
					rightPoint := Point{p.x + rightDelta.x, p.y + rightDelta.y}
					neighborMap[p] = append(neighborMap[p], downPoint)
					neighborMap[p] = append(neighborMap[p], rightPoint)
				} else if connectedDirections[Up] && connectedDirections[Down] {
					animalRune = '|'
					upDelta := directionToPoint(Up)
					downDelta := directionToPoint(Down)
					upPoint := Point{p.x + upDelta.x, p.y + upDelta.y}
					downPoint := Point{p.x + downDelta.x, p.y + downDelta.y}
					neighborMap[p] = append(neighborMap[p], upPoint)
					neighborMap[p] = append(neighborMap[p], downPoint)
				} else if connectedDirections[Left] && connectedDirections[Right] {
					animalRune = '-'
					leftDelta := directionToPoint(Left)
					rightDelta := directionToPoint(Right)
					leftPoint := Point{p.x + leftDelta.x, p.y + leftDelta.y}
					rightPoint := Point{p.x + rightDelta.x, p.y + rightDelta.y}
					neighborMap[p] = append(neighborMap[p], leftPoint)
					neighborMap[p] = append(neighborMap[p], rightPoint)
				} else {
					log.Fatal("Animal is not on a loop")
				}
			}
		}
	}

	return neighborMap, m, animal, animalRune
}

func makeDistances(neighborMap NeighborMap, animal Point) (Distances, Point, int) {
	maxDistance := 0
	maxPoint := Point{-1, -1}
	distances := make(Distances)
	distances[animal] = 0
	queue := make([]Point, 0, 100)
	queue = append(queue, animal)
	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		for _, neighbor := range neighborMap[p] {
			if _, ok := distances[neighbor]; !ok {
				distances[neighbor] = distances[p] + 1
				if distances[neighbor] > maxDistance {
					maxDistance = distances[neighbor]
					maxPoint = neighbor
				}
				queue = append(queue, neighbor)
			}
		}
	}

	return distances, maxPoint, maxDistance
}

func findLoop(neighborMap NeighborMap, animal Point) map[Point]struct{} {
	visited := make(map[Point]struct{})
	visited[animal] = struct{}{}
	next := neighborMap[animal]
	// traverse neighbors until we can't find any more
	for len(next) > 0 {
		p := next[0]
		next = next[1:]
		if _, ok := visited[p]; ok {
			continue
		}
		visited[p] = struct{}{}
		for _, neighbor := range neighborMap[p] {
			if _, ok := visited[neighbor]; !ok {
				next = append(next, neighbor)
			}
		}
	}
	return visited
}

func isInside(m Map, p Point, height int, width int, loop map[Point]struct{}) bool {
	// draw a line from {-1,y} to the point p
	// count the number of times the line crosses the boundary
	// if the number is odd, the point is inside
	// if the number is even, the point is outside
	// if the number is 0, the point is on the boundary
	count := 0
	for x := 0; x <= p.x; x++ {
		checkPoint := Point{x, p.y}
		_, ok := loop[checkPoint]
		if !ok {
			continue
		}
		r, _ := m[checkPoint]
		switch r {
		case 'L', '|', 'J':
			count++
		}
	}
	if count%2 == 1 {
		return true
	}
	return false
}

// prints the map and returns the number of points inside the loop
func printMap(height int, width int, m Map, neighbors NeighborMap, animal Point) int {
	count := 0
	loop := findLoop(neighbors, animal)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			p := Point{x, y}
			if _, ok := loop[p]; ok {
				s := " "
				switch m[p] {
				case 'L':
					s = "┗"
				case 'J':
					s = "┛"
				case '7':
					s = "┓"
				case 'F':
					s = "┏"
				case '|':
					s = "┃"
				case '-':
					s = "━"
				}
				fmt.Print(s)
			} else {
				if isInside(m, p, height, width, loop) {
					fmt.Print("█")
					count++
				} else {
					fmt.Print(string(m[p]))
				}
			}
		}
		fmt.Printf("\n")
	}
	return count
}

func run(input string) string {
	neighborMap, _, animal, _ := makeMaps(input)
	_, maxPoint, maxDistance := makeDistances(neighborMap, animal)
	fmt.Printf("Max distance: %d\n", maxDistance)
	fmt.Printf("Max point: %v\n", maxPoint)

	return fmt.Sprintf("%d", maxDistance)
}

func run2(input string) string {
	neighborMap, m, animal, r := makeMaps(input)
	m[animal] = r

	lines := strings.Split(input, "\n")
	height := len(lines)
	width := len(lines[0])

	count := printMap(height, width, m, neighborMap, animal)

	fmt.Printf("Count: %d\n", count)
	return fmt.Sprintf("%d", count)
}
