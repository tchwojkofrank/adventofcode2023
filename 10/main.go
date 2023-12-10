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

func makeMaps(input string) (NeighborMap, Map, Point) {
	m := make(Map)
	neighborMap := make(NeighborMap)
	lines := strings.Split(input, "\n")
	var animal Point
	for y := 0; y < len(lines); y++ {
		for x := 0; x < len(lines[y]); x++ {
			neighborMap[Point{x, y}] = make([]Point, 0, 4)
		}
	}
	width := len(lines[0])
	for y, line := range lines {
		for x, char := range line {
			p := Point{x, y}
			m[p] = char
			switch char {
			case '|':
				if y > 0 {
					neighborMap[p] = append(neighborMap[p], Point{x, y - 1})
				}
				if y < len(lines)-1 {
					neighborMap[p] = append(neighborMap[p], Point{x, y + 1})
				}
			case '-':
				if x > 0 {
					neighborMap[p] = append(neighborMap[p], Point{x - 1, y})
				}
				if x < len(line)-1 {
					neighborMap[p] = append(neighborMap[p], Point{x + 1, y})
				}
			case 'L':
				if y > 0 {
					neighborMap[p] = append(neighborMap[p], Point{x, y - 1})
				}
				if x < len(line)-1 {
					neighborMap[p] = append(neighborMap[p], Point{x + 1, y})
				}
			case 'J':
				if y > 0 {
					neighborMap[p] = append(neighborMap[p], Point{x, y - 1})
				}
				if x > 0 {
					neighborMap[p] = append(neighborMap[p], Point{x - 1, y})
				}
			case '7':
				if y < len(lines)-1 {
					neighborMap[p] = append(neighborMap[p], Point{x, y + 1})
				}
				if x > 0 {
					neighborMap[p] = append(neighborMap[p], Point{x - 1, y})
				}
			case 'F':
				if y < len(lines)-1 {
					neighborMap[p] = append(neighborMap[p], Point{x, y + 1})
				}
				if x < len(line)-1 {
					neighborMap[p] = append(neighborMap[p], Point{x + 1, y})
				}
			case '.':
			case 'S':
				animal = p
			}
		}
	}
	// update neighbors for the animal
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			if animal.x+dx < 0 || animal.x+dx >= width || animal.y+dy < 0 || animal.y+dy >= len(lines) {
				continue
			}
			np := Point{animal.x + dx, animal.y + dy}
			for _, nn := range neighborMap[np] {
				if nn == animal {
					neighborMap[animal] = append(neighborMap[animal], np)
				}
			}
		}
	}

	return neighborMap, m, animal
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

func run(input string) string {
	neighborMap, _, animal := makeMaps(input)
	_, maxPoint, maxDistance := makeDistances(neighborMap, animal)
	fmt.Printf("Max distance: %d\n", maxDistance)
	fmt.Printf("Max point: %v\n", maxPoint)
	return fmt.Sprintf("%d", maxDistance)
}

func run2(input string) string {
	return ""
}
