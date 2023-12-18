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

type Point struct {
	x, y int
}

func followInstruction(instruction string, current *Point, trench *map[Point]string, min *Point, max *Point) {
	fields := strings.Split(instruction, " ")
	direction := fields[0]
	count, _ := strconv.Atoi(fields[1])
	dp := Point{0, 0}
	rgb := fields[2]
	switch direction[0] {
	case 'R':
		dp.x = 1
	case 'D':
		dp.y = 1
	case 'U':
		dp.y = -1
	case 'L':
		dp.x = -1
	}
	for i := 0; i < count; i++ {
		current.x += dp.x
		current.y += dp.y
		(*trench)[*current] = rgb
	}
	if current.x < min.x {
		min.x = current.x
	}
	if current.y < min.y {
		min.y = current.y
	}
	if current.x > max.x {
		max.x = current.x
	}
	if current.y > max.y {
		max.y = current.y
	}
}

func floodFill(trench map[Point]string, min Point, max Point) (int, map[Point]bool) {
	visited := make(map[Point]bool)
	start := Point{(max.x + min.x) / 2, (max.y + min.y) / 2}
	visited[start] = true
	queue := make([]Point, 0)
	queue = append(queue, start)
	area := 0
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		area++
		for _, dp := range []Point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			next := Point{current.x + dp.x, current.y + dp.y}
			if next.x < min.x || next.x > max.x || next.y < min.y || next.y > max.y {
				continue
			}
			if visited[next] {
				continue
			}
			_, ok := trench[next]
			if ok {
				continue
			}
			visited[next] = true
			queue = append(queue, next)
		}
	}
	return area, visited
}

func run(input string) string {
	instructions := strings.Split(input, "\n")
	trench := make(map[Point]string)
	current := Point{0, 0}
	min := Point{0, 0}
	max := Point{0, 0}
	trench[current] = "#000000"
	for _, instruction := range instructions {
		followInstruction(instruction, &current, &trench, &min, &max)
	}
	_, visited := floodFill(trench, min, max)
	totalArea := 0
	for y := min.y; y <= max.y; y++ {
		for x := min.x; x <= max.x; x++ {
			if visited[Point{x, y}] {
				fmt.Printf("X")
				totalArea++
			} else if trench[Point{x, y}] != "" {
				fmt.Printf("#")
				totalArea++
			} else {
				fmt.Printf(" ")
			}
		}
		fmt.Printf("\n")
	}
	fmt.Printf("Area: %d\n", totalArea)
	return strconv.Itoa(totalArea)
}

type Instruction struct {
	direction string
	count     int
}

type Segment struct {
	start, end Point
}

func getInstructions(instructionStrings []string) []Instruction {
	instructions := make([]Instruction, len(instructionStrings))
	for i, instructionString := range instructionStrings {
		fields := strings.Split(instructionString, " ")
		rgb := fields[2][2 : len(fields[2])-1]
		hex := rgb[0 : len(rgb)-1]
		direction := rgb[len(rgb)-1:]
		count, _ := strconv.ParseUint(hex, 16, 64)
		instructions[i] = Instruction{direction, int(count)}
	}
	return instructions
}

// 0 = right, 1 = down, 2 = left, 3 = up
// returns all the segments defined by the instructions, and the min and max points
func getSegments(instructions []Instruction) ([]Segment, Segment) {
	minPoint := Point{1000000, 1000000}
	maxPoint := Point{-1000000, -1000000}
	segments := make([]Segment, len(instructions))
	current := Point{0, 0}
	segment := Segment{current, current}
	for i, instruction := range instructions {
		switch instruction.direction {
		case "0":
			segment.end.x += instruction.count
		case "1":
			segment.end.y += instruction.count
		case "2":
			segment.end.x -= instruction.count
		case "3":
			segment.end.y -= instruction.count
		}
		if segment.end.x < minPoint.x {
			minPoint.x = segment.end.x
		}
		if segment.end.y < minPoint.y {
			minPoint.y = segment.end.y
		}
		if segment.end.x > maxPoint.x {
			maxPoint.x = segment.end.x
		}
		if segment.end.y > maxPoint.y {
			maxPoint.y = segment.end.y
		}
		segments[i] = segment
		segment.start = segment.end
	}
	return segments, Segment{minPoint, maxPoint}
}

func rectArea(segment Segment, minmax Segment) int {
	// calculating the area of the rectangle between the segment and
	// the line defined by y = minmax.start.y
	return (segment.start.x - segment.end.x) * (segment.end.y - minmax.start.y)
}

func isUp(segment Segment) bool {
	return segment.end.y > segment.start.y
}

func isLeft(segment Segment) bool {
	return segment.end.x < segment.start.x
}

func run2(input string) string {
	instructionStrings := strings.Split(input, "\n")
	instructions := getInstructions(instructionStrings)
	segments, minMaxSegment := getSegments(instructions)

	// We're going to add all rectangles that are defined by
	// going from right to left to the minimum y value, and subtract
	// all rectangles going from left to right to the minimum y value.
	// This will give us the area of the interior region of the path
	// defined by the segments.
	// Because of pixels, we also need to add the length of all upward
	// and leftward segments.
	// There is also always one corner pixel that is not otherwise counted.
	totalArea := 1
	for _, segment := range segments {
		totalArea += rectArea(segment, minMaxSegment)
		// because of pixels, add the length of upward and leftward segments
		if isUp(segment) {
			totalArea += segment.end.y - segment.start.y
		}
		if isLeft(segment) {
			totalArea += segment.start.x - segment.end.x
		}
	}
	fmt.Printf("Area: %d\n", totalArea)
	return strconv.Itoa(totalArea)
}
