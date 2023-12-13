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

func isContiguous(record string, index int, contiguousCount int) bool {
	i := index
	for ; i < len(record) && i < index+contiguousCount; i++ {
		if record[i] != 'X' && record[i] != 'S' {
			return false
		}
	}
	// make sure the next character is not an X or it's at the end of the string
	if (i < len(record) && record[i] != 'X') || i == len(record) {
		return true
	}
	return false
}

// return the indices of the first position of all possible matches for the first contiguousCount
func findFirstMatches(record string, contiguousCount int) []int {
	i := 0

	// find the first X in the record
	firstX := strings.Index(record, "X")
	if firstX == -1 {
		firstX = len(record)
	}
	// the last position to try is the first X
	end := min(firstX, len(record)-contiguousCount)

	matches := make([]int, 0)

	for i <= end {
		if isContiguous(record, i, contiguousCount) {
			matches = append(matches, i)
		}
		i++
	}

	return matches
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

type CalcCache map[string]int

var myCache CalcCache = make(CalcCache)

func makeCacheKey(record string, contiguousCounts []int) string {
	return fmt.Sprintf("%s-%v", record, contiguousCounts)
}

func countArrangements(record string, contiguousCounts []int, depth int) int {
	// fmt.Print(strings.Repeat("\t", depth))
	// fmt.Printf("record: %s, contiguousCounts: %v\n", record, contiguousCounts)
	if len(contiguousCounts) == 0 {
		return 0
	}
	cachedResult, ok := myCache[makeCacheKey(record, contiguousCounts)]
	if ok {
		return cachedResult
	}
	count := 0
	// iterate over each possible placement of the first contiguousCount, and recurse
	contiguousCount := contiguousCounts[0]
	firstMatches := findFirstMatches(record, contiguousCount)
	if len(firstMatches) == 0 {
		return 0
	}
	i := 0
	for _, i = range firstMatches {
		newRecord := record[i:]
		// fmt.Print(strings.Repeat("\t", depth))
		// fmt.Printf("\tnewRecord: %s, counts: %v\n", newRecord, contiguousCounts)
		if len(firstMatches) == 0 {
			break
		}
		if len(contiguousCounts) == 1 {
			if contiguousCounts[0]+1 < len(newRecord) {
				newRecord = newRecord[contiguousCounts[0]+1:]
				if strings.Count(newRecord, "X") == 0 {
					count++
				}
			} else {
				count++
			}
		} else {
			if contiguousCounts[0]+1 < len(newRecord) {
				count += countArrangements(newRecord[contiguousCounts[0]+1:], contiguousCounts[1:], depth+1)
			}
		}
		// fmt.Print(strings.Repeat("\t", depth))
		// fmt.Printf("count: %d\n", count)
	}

	myCache[makeCacheKey(record, contiguousCounts)] = count
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
		// patternString := makePatern(contiguousCounts)
		// pattern := regexp.MustCompile(patternString)
		myCache = make(CalcCache)
		arrangementCount := countArrangements(conditionString, contiguousCounts, 0)
		// fmt.Printf("Condition: %s\tCounts: %v\tArrangements: %d\n", conditionString, contiguousCounts, arrangementCount)
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
	for _, line := range lines {
		fields := strings.Split(line, " ")
		conditionString := unfold(fields[0], "?")
		conditionString = replaceChars(conditionString)
		contiguousCounts := getContiguousCounts(unfold(fields[1], ","))
		// patternString := makePatern(contiguousCounts)
		// pattern := regexp.MustCompile(patternString)
		myCache = make(CalcCache)
		arrangementCount := countArrangements(conditionString, contiguousCounts, 0)
		totalArrangements += arrangementCount
	}
	fmt.Printf("Total arrangements: %d\n", totalArrangements)
	return fmt.Sprintf("%d", totalArrangements)
}
