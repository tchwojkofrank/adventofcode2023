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

func getReadings(line string) []int {
	readingStrings := strings.Split(line, " ")
	readings := make([]int, len(readingStrings))
	for i, readingString := range readingStrings {
		readings[i], _ = strconv.Atoi(readingString)
	}
	return readings
}

func getDiffs(readings []int) []int {
	diffs := make([]int, len(readings)-1)
	for i := 0; i < len(readings)-1; i++ {
		diffs[i] = readings[i+1] - readings[i]
	}
	return diffs
}

func isZero(diffs []int) bool {
	for _, diff := range diffs {
		if diff != 0 {
			return false
		}
	}
	return true
}

func getPrediction(readings []int) int {
	diffs := getDiffs(readings)
	if isZero(diffs) {
		return 0
	}
	prediction := getPrediction(diffs) + diffs[len(diffs)-1]
	return prediction
}

func getPrediction2(readings []int) int {
	diffs := getDiffs(readings)
	if isZero(diffs) {
		return 0
	}
	prediction := diffs[0] - getPrediction2(diffs)
	return prediction
}

func run(input string) string {
	lines := strings.Split(input, "\n")
	totalPrediction := 0
	for _, line := range lines {
		readings := getReadings(line)
		prediction := getPrediction(readings) + readings[len(readings)-1]
		totalPrediction += prediction
	}
	fmt.Printf("Total prediction: %v\n", totalPrediction)
	return fmt.Sprintf("%v", totalPrediction)
}

func run2(input string) string {
	lines := strings.Split(input, "\n")
	totalPrediction := 0
	for _, line := range lines {
		readings := getReadings(line)
		prediction := readings[0] - getPrediction2(readings)
		totalPrediction += prediction
	}
	fmt.Printf("Total prediction: %v\n", totalPrediction)
	return fmt.Sprintf("%v", totalPrediction)
}
