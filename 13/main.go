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

func checkSymmetry(midpoint int, section []string, targetDiff int) bool {
	up := midpoint
	down := midpoint + 1
	totalDiff := 0
	for up >= 0 && down < len(section) {
		stringDiff := findDiff(section[up], section[down])
		totalDiff += stringDiff
		if totalDiff > targetDiff {
			return false
		}
		up--
		down++
	}
	return totalDiff == targetDiff
}

func findSymmetry(section []string, targetDiff int) (int, bool) {
	values := make([]int, 0)
	for y := 0; y < len(section)-1; y++ {
		if checkSymmetry(y, section, targetDiff) {
			values = append(values, y)
		}
	}
	if len(values) == 0 {
		return -1, false
	}
	sum := 0
	for _, value := range values {
		sum += value
	}
	return sum, true
}

// we only care about the first diff
func findDiff(s1 string, s2 string) int {
	diff := 0
	for i := 0; i < len(s1) && diff < 2; i++ {
		if s1[i] != s2[i] {
			diff++
		}
	}
	return diff
}

func findVerticalSymmetry(section []string, targetDiff int) (int, bool) {
	return findSymmetry(section, targetDiff)
}

func transpose(section []string) []string {
	transposed := make([]string, len(section[0]))
	for x := 0; x < len(section[0]); x++ {
		for y := 0; y < len(section); y++ {
			transposed[x] += string(section[y][x])
		}
	}
	return transposed
}

func findHorizontalSymmetry(section []string, targetDiff int) (int, bool) {
	transposedSection := transpose(section)
	return findSymmetry(transposedSection, targetDiff)
}

func run(input string) string {
	sections := strings.Split(input, "\n\n")
	totalScore := 0
	for _, section := range sections {
		verticalScore := 0
		horizontalScore := 0
		lines := strings.Split(section, "\n")
		v, vFound := findVerticalSymmetry(lines, 0)
		if vFound {
			verticalScore = v + 1
		}
		h, hFound := findHorizontalSymmetry(lines, 0)
		if hFound {
			horizontalScore = h + 1
		}
		totalScore += horizontalScore + 100*verticalScore
	}
	fmt.Printf("Total score: %d\n", totalScore)
	return fmt.Sprintf("%d", totalScore)
}

func run2(input string) string {
	sections := strings.Split(input, "\n\n")
	totalScore := 0
	for _, section := range sections {
		verticalScore := 0
		horizontalScore := 0
		lines := strings.Split(section, "\n")
		v, vFound := findVerticalSymmetry(lines, 1)
		if vFound {
			verticalScore = v + 1
		}
		h, hFound := findHorizontalSymmetry(lines, 1)
		if hFound {
			horizontalScore = h + 1
		}
		totalScore += horizontalScore + 100*verticalScore
	}
	fmt.Printf("Total score: %d\n", totalScore)
	return fmt.Sprintf("%d", totalScore)
}
