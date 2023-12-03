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

type Value struct {
	start Point
	end   Point
	value int
}

var objects = make(map[Point]string)
var values = make([]Value, 0)

func getThings(line string) []string {
	values := strings.Split(line, ".")
	things := make([]string, 0)
	for _, value := range values {
		if len(value) > 0 {
			things = append(things, value)
		}
	}
	return things
}

func updateMaps(line string, p Point) {
	value := 0
	start := p
	end := p
	for i, r := range line {
		if r >= '0' && r <= '9' {
			value = value*10 + int(r-'0')
			end = Point{p.x + i, p.y}
		} else {
			// add value to values map if not zero
			if value > 0 {
				values = append(values, Value{start, end, value})
				value = 0
			}
			// add this object to objects map if it is not a '.'
			if r != '.' {
				objects[Point{p.x + i, p.y}] = string(r)
			}
			start = Point{p.x + i + 1, p.y}
			end = start
		}
	}
	if value > 0 {
		values = append(values, Value{start, end, value})
		value = 0
	}
}

func checkValue(v Value) bool {
	// if there is an object adjecent to the value, return true
	for x := v.start.x - 1; x <= v.end.x+1; x++ {
		for y := v.start.y - 1; y <= v.end.y+1; y++ {
			_, ok := objects[Point{x, y}]
			if ok {
				return true
			}
		}
	}
	return false
}

func run(input string) string {
	lines := strings.Split(input, "\n")
	for y, line := range lines {
		updateMaps(line, Point{0, y})
	}
	sum := 0
	everything := 0
	for _, v := range values {
		everything += v.value
		if checkValue(v) {
			sum += v.value
		}
	}
	fmt.Printf("Sum: %d\n", sum)
	return fmt.Sprintf("%d", sum)
}

func getObjects(lines []string) map[Point]string {
	objects := make(map[Point]string)
	for y, line := range lines {
		for x, r := range line {
			if r != '.' && (r < '0' || r > '9') {
				objects[Point{x, y}] = string(r)
			}
		}
	}
	return objects
}

func isAdjecent(p1 Point, p2 Point, objects map[Point]string) bool {
	for y := p1.y - 1; y <= p2.y+1; y++ {
		for x := p1.x - 1; x <= p2.x+1; x++ {
			_, ok := objects[Point{x, y}]
			if ok {
				return true
			}
		}
	}
	return false
}

func getAdjacentValues(p Point) []Value {
	adjacents := make([]Value, 0)
	// get all values adjecent to p
	for _, v := range values {
		if p.x >= v.start.x-1 && p.x <= v.end.x+1 && p.y >= v.start.y-1 && p.y <= v.end.y+1 {
			adjacents = append(adjacents, v)
		}
	}
	return adjacents
}

func run2(input string) string {
	sum := 0
	// check objects for gears
	for p, o := range objects {
		if o == "*" {
			adjacents := getAdjacentValues(p)
			if len(adjacents) == 2 {
				sum += adjacents[0].value * adjacents[1].value
			}
		}
	}
	fmt.Printf("Sum: %d\n", sum)
	return fmt.Sprintf("%d", sum)
}
