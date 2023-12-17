package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"chwojkofrank.com/cursor"
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

type Direction struct {
	dx int
	dy int
}

type Vector struct {
	p Point
	d Direction
}

type Beam []Vector

func mapInput(input string) (map[Point]rune, int, int) {
	grid := make(map[Point]rune)
	lines := strings.Split(input, "\n")
	height := len(lines)
	width := len(lines[0])
	for y, line := range lines {
		for x, c := range line {
			grid[Point{x, y}] = c
		}
	}
	return grid, width, height
}

func nextVector(v Vector, r rune) (Vector, *Vector) {
	newVector := v
	switch r {
	case '|':
		if v.d.dx == 0 {
			newVector.p.x += newVector.d.dx
			newVector.p.y += newVector.d.dy
			return newVector, nil
		} else {
			teeVector := Vector{newVector.p, Direction{0, 1}}
			newVector.d = Direction{0, -1}
			newVector.p.x += newVector.d.dx
			newVector.p.y += newVector.d.dy
			return newVector, &teeVector
		}
	case '-':
		if v.d.dy == 0 {
			newVector.p.x += newVector.d.dx
			newVector.p.y += newVector.d.dy
			return newVector, nil
		} else {
			teeVector := Vector{newVector.p, Direction{1, 0}}
			newVector.d = Direction{-1, 0}
			newVector.p.x += newVector.d.dx
			newVector.p.y += newVector.d.dy
			return newVector, &teeVector
		}
	case '\\':
		if v.d.dx == 0 {
			if v.d.dy == 1 {
				newVector.d = Direction{1, 0}
			} else {
				newVector.d = Direction{-1, 0}
			}
		} else {
			if v.d.dx == 1 {
				newVector.d = Direction{0, 1}
			} else {
				newVector.d = Direction{0, -1}
			}
		}
		newVector.p.x += newVector.d.dx
		newVector.p.y += newVector.d.dy
		return newVector, nil
	case '/':
		if v.d.dx == 0 {
			if v.d.dy == 1 {
				newVector.d = Direction{-1, 0}
			} else {
				newVector.d = Direction{1, 0}
			}
		} else {
			if v.d.dx == 1 {
				newVector.d = Direction{0, -1}
			} else {
				newVector.d = Direction{0, 1}
			}
		}
		newVector.p.x += newVector.d.dx
		newVector.p.y += newVector.d.dy
		return newVector, nil
	default:
		newVector.p.x += newVector.d.dx
		newVector.p.y += newVector.d.dy
		return newVector, nil
	}
}

const (
	Up    = 1
	Left  = 2
	Down  = 4
	Right = 8
)

func vectorToDirection(v Vector) int {
	if v.d.dx == 0 {
		if v.d.dy == 1 {
			return Down
		} else {
			return Up
		}
	} else {
		if v.d.dx == 1 {
			return Right
		} else {
			return Left
		}
	}
}

func directionToRune(d int) rune {
	switch d {
	case 1:
		return '╹'
	case 2:
		return '╸'
	case 3:
		return '┛'
	case 4:
		return '╻'
	case 5:
		return '┃'
	case 6:
		return '┓'
	case 7:
		return '┫'
	case 8:
		return '╺'
	case 9:
		return '┗'
	case 10:
		return '━'
	case 11:
		return '┻'
	case 12:
		return '┏'
	case 13:
		return '┣'
	case 14:
		return '┳'
	case 15:
		return '╋'
	default:
		return '.'
	}
}

func printVector(v Vector, mappedPoints map[Point]int) {
	cursor.Position(v.p.x+5, v.p.y+5)
	fmt.Printf("%c", directionToRune(mappedPoints[v.p]))
}

func clearBoard() {
	cursor.Position(0, 0)
	cursor.Clear()
}

func traceBeams(grid map[Point]rune, vectorsToTrace []Vector, energizedPoints *map[Point]bool, width int, height int) {
	vector := vectorsToTrace[0]
	mappedVectors := make(map[Vector]struct{})
	mappedPoints := make(map[Point]int)
	(*energizedPoints)[vector.p] = true
	mappedVectors[vector] = struct{}{}
	mappedPoints[vector.p] = Left
	for len(vectorsToTrace) > 0 {
		newVector, newTee := nextVector(vector, grid[vector.p])
		if newTee != nil {
			if _, ok := mappedVectors[*newTee]; !ok {
				vectorsToTrace = append(vectorsToTrace, *newTee)
				mappedPoints[newTee.p] = vectorToDirection(*newTee)
				printVector(*newTee, mappedPoints)
			} else {
				mappedPoints[newTee.p] |= vectorToDirection(*newTee)
				printVector(*newTee, mappedPoints)
			}
		}
		_, mapped := mappedVectors[newVector]
		if newVector.p.x >= 0 && newVector.p.y >= 0 && newVector.p.x < width && newVector.p.y < height && !mapped {
			(*energizedPoints)[newVector.p] = true
			mappedVectors[newVector] = struct{}{}
			mappedPoints[newVector.p] = vectorToDirection(newVector)
			printVector(newVector, mappedPoints)
			vector = newVector
		} else {
			vectorsToTrace = vectorsToTrace[1:]
			if len(vectorsToTrace) > 0 {
				vector = vectorsToTrace[0]
			}
			if !(newVector.p.x >= 0 && newVector.p.y >= 0 && newVector.p.x < width && newVector.p.y < height) && mapped {
				mappedPoints[newVector.p] |= vectorToDirection(newVector)
				printVector(newVector, mappedPoints)
			}
		}
	}
}

func run(input string) string {

	clearBoard()
	grid, width, height := mapInput(input)
	newVector := Vector{Point{0, 0}, Direction{1, 0}}
	energizedPoints := make(map[Point]bool)
	vectorsToTrace := make([]Vector, 1)
	vectorsToTrace[0] = newVector

	traceBeams(grid, vectorsToTrace, &energizedPoints, width, height)
	//count energized points
	count := 0
	for _, v := range energizedPoints {
		if v {
			count++
		}
	}

	fmt.Printf("Energized points: %d\n", count)

	return fmt.Sprintf("%d", count)
}

func getEnergy(grid map[Point]rune, width int, height int, startV Vector) int {
	newVector := startV
	energizedPoints := make(map[Point]bool)
	vectorsToTrace := make([]Vector, 1)
	vectorsToTrace[0] = newVector

	traceBeams(grid, vectorsToTrace, &energizedPoints, width, height)
	//count eneregized points
	count := 0
	for _, v := range energizedPoints {
		if v {
			count++
		}
	}
	return count
}

func run2(input string) string {
	grid, width, height := mapInput(input)

	maxEnergy := 0
	// check all top and bottom row possibilities
	for x := 0; x < width; x++ {
		energy := getEnergy(grid, width, height, Vector{Point{x, 0}, Direction{0, 1}})
		if energy > maxEnergy {
			maxEnergy = energy
		}
		energy = getEnergy(grid, width, height, Vector{Point{x, height - 1}, Direction{0, -1}})
		if energy > maxEnergy {
			maxEnergy = energy
		}
	}

	// check all left and right column possibilities
	for y := 0; y < height; y++ {
		energy := getEnergy(grid, width, height, Vector{Point{0, y}, Direction{1, 0}})
		if energy > maxEnergy {
			maxEnergy = energy
		}
		energy = getEnergy(grid, width, height, Vector{Point{width - 1, y}, Direction{-1, 0}})
		if energy > maxEnergy {
			maxEnergy = energy
		}
	}
	fmt.Printf("Maximum Energized points: %d\n", maxEnergy)

	return fmt.Sprintf("%d", maxEnergy)
}
