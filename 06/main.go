package main

import (
	"fmt"
	"log"
	"math"
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

func getTimes(line string) []float64 {
	line = strings.TrimPrefix(line, "Time:")
	line = strings.TrimSpace(line)
	timeStrings := strings.Fields(line)
	times := make([]float64, len(timeStrings))
	for i, timeString := range timeStrings {
		t, _ := strconv.Atoi(timeString)
		times[i] = float64(t)
	}
	return times
}

func getDistances(line string) []float64 {
	line = strings.TrimPrefix(line, "Distance:")
	line = strings.TrimSpace(line)
	distanceStrings := strings.Fields(line)
	distances := make([]float64, len(distanceStrings))
	for i, distanceString := range distanceStrings {
		d, _ := strconv.Atoi(distanceString)
		distances[i] = float64(d)
	}
	return distances
}

func run(input string) string {
	lines := strings.Split(input, "\n")
	times := getTimes(lines[0])
	var x1, x2 float64
	distances := getDistances(lines[1])
	total := 1
	for i, t := range times {
		d := distances[i]
		x1 = t/2 + math.Sqrt(t*t-4*d)/-2 + 0.01
		x2 = t/2 - math.Sqrt(t*t-4*d)/-2 - 0.01

		minTime := int(math.Ceil(x1) + 0.1)
		maxTime := int(math.Floor(x2) + 0.1)

		winCount := maxTime - minTime + 1
		total = total * winCount
	}

	fmt.Printf("Total: %d\n", total)
	return fmt.Sprintf("%d", total)
}

func run2(input string) string {
	return ""
}
