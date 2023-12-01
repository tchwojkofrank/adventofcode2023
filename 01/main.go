package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const (
	one   = "one"
	two   = "two"
	three = "three"
	four  = "four"
	five  = "five"
	six   = "six"
	seven = "seven"
	eight = "eight"
	nine  = "nine"
	zero  = "zero"
)

var numberMap = map[string]int{one: 1, two: 2, three: 3, four: 4, five: 5, six: 6, seven: 7, eight: 8, nine: 9, zero: 0}

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
	start := time.Now()
	text := readInput(inputName)
	run(text)
	end := time.Now()
	fmt.Printf("Running time: %v\n", end.Sub(start))
}

func getFirstDigitFromLine(line string) int {
	for i, c := range line {
		if c >= '0' && c <= '9' {
			return int(c - '0')
		}
		for number, value := range numberMap {
			if strings.HasPrefix(line[i:], number) {
				return value
			}
		}
	}
	return -1
}

func getLastDigitFromLine(line string) int {
	for i := len(line) - 1; i >= 0; i-- {
		c := line[i]
		if c >= '0' && c <= '9' {
			return int(c - '0')
		}
		for number, value := range numberMap {
			if strings.HasPrefix(line[i:], number) {
				return value
			}
		}
	}
	return -1
}

func getLineValue(line string) int {
	first := getFirstDigitFromLine(line)
	last := getLastDigitFromLine(line)
	fmt.Printf("First: %d, Last: %d\n", first, last)
	return first*10 + last
}

func getSumFromLines(lines []string) int {
	sum := 0
	for _, line := range lines {
		sum += getLineValue(line)
	}
	return sum
}

func run(input string) {
	lines := strings.Split(input, "\n")
	linesValue := getSumFromLines(lines)
	fmt.Printf("Sum of all lines: %d\n", linesValue)
}
