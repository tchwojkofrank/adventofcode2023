package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
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

// let's change out the characters to make it easier to work with
func replaceChars(line string) string {
	line = strings.ReplaceAll(line, "#", "X")
	line = strings.ReplaceAll(line, ".", "o")
	line = strings.ReplaceAll(line, "?", "S")
	return line
}

func getContiguousCounts(line string) []int {
	contiguousCountStrings := strings.Split(line, ",")
	contiguousCounts := make([]int, len(contiguousCountStrings))
	for i, contiguousCountString := range contiguousCountStrings {
		contiguousCounts[i], _ = strconv.Atoi(contiguousCountString)
	}
	return contiguousCounts
}

func makePatern(contiguousCounts []int) string {
	pattern := "^([oS]*)"
	for i, contiguousCount := range contiguousCounts {
		pattern += "("
		pattern += strings.Repeat("[XS]", contiguousCount)
		pattern += ")"
		if i < len(contiguousCounts)-1 {
			pattern += "([oS]+)"
		}
	}
	pattern += "([oS]*)$"
	return pattern
}

func countArrangements(record string, pattern *regexp.Regexp) int {
	// find the index of the first wildcard
	firstWildcardIndex := strings.Index(record, "S")
	// if there are no wildcards left, then we can just check the pattern
	if firstWildcardIndex == -1 {
		result := pattern.FindString(record)
		if result != "" {
			return 1
		} else {
			return 0
		}
	}
	// otherwise, we need to check all possible combinations of the first wildcard
	count := 0
	newRecordX := record[:firstWildcardIndex] + "X" + record[firstWildcardIndex+1:]
	resultX := pattern.FindString(newRecordX)
	if resultX != "" {
		count += countArrangements(newRecordX, pattern)
	}
	newRecordO := record[:firstWildcardIndex] + "o" + record[firstWildcardIndex+1:]
	resultO := pattern.FindString(newRecordO)
	if resultO != "" {
		count += countArrangements(newRecordO, pattern)
	}
	return count
}

func run(input string) string {
	lines := strings.Split(input, "\n")
	totalArrangements := 0
	for _, line := range lines {
		fields := strings.Split(line, " ")
		conditionString := fields[0]
		conditionString = replaceChars(conditionString)
		contiguousCounts := getContiguousCounts(fields[1])
		patternString := makePatern(contiguousCounts)
		pattern := regexp.MustCompile(patternString)
		arrangementCount := countArrangements(conditionString, pattern)
		totalArrangements += arrangementCount
	}
	fmt.Printf("Total arrangements: %d\n", totalArrangements)
	return fmt.Sprintf("%d", totalArrangements)
}

func unfold(line string, separator string) string {
	return line + separator + line + separator + line + separator + line + separator + line
}

func run2(input string) string {
	lines := strings.Split(input, "\n")
	totalArrangements := 0
	for i, line := range lines {
		fmt.Printf("Processing line %d\n", i)
		fields := strings.Split(line, " ")
		conditionString := unfold(fields[0], "?")
		conditionString = replaceChars(conditionString)
		contiguousCounts := getContiguousCounts(unfold(fields[1], ","))
		patternString := makePatern(contiguousCounts)
		pattern := regexp.MustCompile(patternString)
		arrangementCount := countArrangements(conditionString, pattern)
		totalArrangements += arrangementCount
	}
	fmt.Printf("Total arrangements: %d\n", totalArrangements)
	return fmt.Sprintf("%d", totalArrangements)
}
