package main

import (
	"fmt"
	"log"
	"os"
	"slices"
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

type StoneMap struct {
	m      map[Point]rune
	width  int
	height int
}

func readStoneMap(input string) StoneMap {
	lines := strings.Split(input, "\n")
	slices.Reverse(lines)
	var stoneMap StoneMap
	stoneMap.m = make(map[Point]rune)
	stoneMap.width = len(lines[0])
	stoneMap.height = len(lines)
	for y, line := range lines {
		for x, stone := range line {
			if stone == '#' || stone == 'O' {
				stoneMap.m[Point{x, y}] = stone
			}
		}
	}
	return stoneMap
}

func (stoneMap *StoneMap) canRollInto(p Point) bool {
	if p.x < 0 || p.x >= stoneMap.width || p.y < 0 || p.y >= stoneMap.height {
		return false
	}
	_, ok := stoneMap.m[p]
	return !ok
}

func (stoneMap *StoneMap) tiltRow(y int) {
	for x := 0; x < stoneMap.width; x++ {
		stone := stoneMap.m[Point{x, y}]
		if stone == 'O' {
			for z := y + 1; z < stoneMap.height && stoneMap.canRollInto(Point{x, z}); z++ {
				delete(stoneMap.m, Point{x, z - 1})
				stoneMap.m[Point{x, z}] = 'O'
			}
		}
	}
}

func (stoneMap *StoneMap) tilt() {
	for y := stoneMap.height - 1; y >= 0; y-- {
		stoneMap.tiltRow(y)
	}
}

func (stoneMap *StoneMap) load() int {
	load := 0
	for k, v := range stoneMap.m {
		if v == 'O' {
			load += k.y + 1
		}
	}
	return load
}

func (stoneMap StoneMap) String() string {
	result := ""
	for y := stoneMap.height - 1; y >= 0; y-- {
		for x := 0; x < stoneMap.width; x++ {
			r, ok := stoneMap.m[Point{x, y}]
			if !ok {
				r = '.'
			}
			result += string(r)
		}
		trailer := fmt.Sprintf("\t%v\n", y+1)
		result += trailer
	}
	return result
}

func run(input string) string {
	stoneMap := readStoneMap(input)
	stoneMap.tilt()
	load := stoneMap.load()
	fmt.Printf("%v\n", stoneMap)
	fmt.Printf("Load: %d\n", load)

	return fmt.Sprintf("%d", load)
}

func run2(input string) string {
	return ""
}
