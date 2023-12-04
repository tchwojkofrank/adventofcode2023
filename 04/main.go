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

// getCardID returns the id of a card
func getCardID(line string) int {
	line = strings.Trim(line[5:], " ")
	cardID, _ := strconv.Atoi(line)
	return cardID
}

func getValues(line string) map[int]struct{} {
	numbers := make(map[int]struct{}, (len(line)+1)/3)
	for i := 0; i < len(line); i += 3 {
		r1 := rune(line[i])
		r2 := rune(line[i+1])
		number := 0
		if r1 != ' ' {
			number = int(r1-'0')*10 + int(r2-'0')
		} else {
			number = int(r2 - '0')
		}
		numbers[number] = struct{}{}
	}
	return numbers
}

// getCardValue returns the id and value of a card
func getCardIdValue(line string) (int, int) {
	value := 0
	cardStrings := strings.Split(line, ":")
	cardID := getCardID(cardStrings[0])
	cardSplit := strings.Split(cardStrings[1], "|")
	winners := getValues(cardSplit[0][1 : len(cardSplit[0])-1])
	myNumbers := getValues(cardSplit[1][1:])
	for number := range myNumbers {
		if _, ok := winners[number]; ok {
			if value == 0 {
				value = 1
			} else {
				value = value * 2
			}
		}
	}

	return cardID, value
}

func getCardIdValue2(line string) (int, int) {
	value := 0
	cardStrings := strings.Split(line, ":")
	cardID := getCardID(cardStrings[0])
	cardSplit := strings.Split(cardStrings[1], "|")
	winners := getValues(cardSplit[0][1 : len(cardSplit[0])-1])
	myNumbers := getValues(cardSplit[1][1:])
	for number := range myNumbers {
		if _, ok := winners[number]; ok {
			if value == 0 {
				value = 1
			} else {
				value = value + 1
			}
		}
	}

	return cardID, value
}

func run(input string) string {
	lines := strings.Split(input, "\n")
	sum := 0
	for _, line := range lines {
		_, cardValue := getCardIdValue(line)
		sum += cardValue
	}
	fmt.Printf("Sum: %v\n", sum)
	return fmt.Sprintf("%v", sum)
}

func run2(input string) string {
	lines := strings.Split(input, "\n")
	cardCopyCount := make(map[int]int, len(lines))
	for i := 1; i <= len(lines); i++ {
		cardCopyCount[i] = 1
	}
	for index, line := range lines {
		cardID := index + 1
		_, cardValue := getCardIdValue2(line)
		for i := cardID + 1; i <= cardID+cardValue && i <= len(lines); i++ {
			cardCopyCount[i] += cardCopyCount[cardID]
		}
	}
	sum := 0
	for _, count := range cardCopyCount {
		sum += count
	}

	fmt.Printf("Sum: %v\n", sum)
	return fmt.Sprintf("%v", sum)
}
