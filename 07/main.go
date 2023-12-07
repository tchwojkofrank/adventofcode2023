package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"sort"
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

type Hand struct {
	deal []int
	hand []int
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

var ValueMap1 = map[rune]int{
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'T': 10,
	'J': 11,
	'Q': 12,
	'K': 13,
	'A': 14,
}

var ValueMap2 = map[rune]int{
	'J': 1,
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'T': 10,
	'Q': 11,
	'K': 12,
	'A': 13,
}

var ValueMap = ValueMap1

func (h Hand) is5ofKind() bool {
	return h.hand[0] == h.hand[4]
}

func (h Hand) is4ofKind() bool {
	return h.hand[0] == h.hand[3] || h.hand[1] == h.hand[4]
}

func (h Hand) isFullHouse() bool {
	return h.hand[0] == h.hand[2] && h.hand[3] == h.hand[4] ||
		h.hand[0] == h.hand[1] && h.hand[2] == h.hand[4]
}

func (h Hand) is3ofKind() bool {
	return h.hand[0] == h.hand[2] || h.hand[1] == h.hand[3] || h.hand[2] == h.hand[4]
}

func (h Hand) is2Pairs() bool {
	return h.hand[0] == h.hand[1] && h.hand[2] == h.hand[3] ||
		h.hand[0] == h.hand[1] && h.hand[3] == h.hand[4] ||
		h.hand[1] == h.hand[2] && h.hand[3] == h.hand[4]
}

func (h Hand) isPair() bool {
	return h.hand[0] == h.hand[1] || h.hand[1] == h.hand[2] ||
		h.hand[2] == h.hand[3] || h.hand[3] == h.hand[4]
}

func (h Hand) handValue() int {
	if h.is5ofKind() {
		return 6
	}
	if h.is4ofKind() {
		return 5
	}
	if h.isFullHouse() {
		return 4
	}
	if h.is3ofKind() {
		return 3
	}
	if h.is2Pairs() {
		return 2
	}
	if h.isPair() {
		return 1
	}
	return 0
}

func (h Hand) jokerCount() int {
	i := 0
	for ; i < 5 && h.hand[i] == 1; i++ {
	}
	return i
}

func (h Hand) is5ofKind2(j int) bool {
	switch j {
	case 0:
		return h.is5ofKind()
	case 1:
		return h.is4ofKind()
	case 2:
		return h.isFullHouse()
	case 3:
		return h.isFullHouse()
	case 4:
		return true
	case 5:
		return true
	}
	return false
}

func (h Hand) is4ofKind2(j int) bool {
	switch j {
	case 0:
		return h.is4ofKind()
	case 1:
		return h.is3ofKind()
	case 2:
		return h.is2Pairs()
	case 3:
		return true
	case 4:
		return true
	case 5:
		return true
	}
	return false
}

func (h Hand) isFullHouse2(j int) bool {
	switch j {
	case 0:
		return h.isFullHouse()
	case 1:
		return h.is3ofKind() || h.is2Pairs()
	case 2:
		return h.is2Pairs()
	case 3:
		return h.isFullHouse()
	case 4:
		return true
	case 5:
		return true
	}
	return false
}

func (h Hand) is3ofKind2(j int) bool {
	switch j {
	case 0:
		return h.is3ofKind()
	case 1:
		return h.isPair()
	case 2:
		return true
	case 3:
		return true
	case 4:
		return true
	case 5:
		return true
	}
	return false
}

func (h Hand) is2Pairs2(j int) bool {
	switch j {
	case 0:
		return h.is2Pairs()
	case 1:
		return h.isPair()
	case 2:
		return true
	case 3:
		return true
	case 4:
		return true
	case 5:
		return true
	}
	return false
}

func (h Hand) isPair2(j int) bool {
	switch j {
	case 0:
		return h.isPair()
	case 1:
		return true
	case 2:
		return true
	case 3:
		return true
	case 4:
		return true
	case 5:
		return true
	}
	return false
}

func (h Hand) handValue2() int {
	jokerCount := h.jokerCount()
	if h.is5ofKind2(jokerCount) {
		return 6
	}
	if h.is4ofKind2(jokerCount) {
		return 5
	}
	if h.isFullHouse2(jokerCount) {
		return 4
	}
	if h.is3ofKind2(jokerCount) {
		return 3
	}
	if h.is2Pairs2(jokerCount) {
		return 2
	}
	if h.isPair2(jokerCount) {
		return 1
	}
	return 0
}

func getHandAndBid(line string) (Hand, int) {
	deal := make([]int, 5)
	hand := make([]int, 5)
	fields := strings.Split(line, " ")
	for i, r := range fields[0] {
		deal[i] = ValueMap[r]
		hand[i] = ValueMap[r]
	}
	sort.IntSlice(hand).Sort()
	bid, _ := strconv.Atoi(fields[1])
	return Hand{deal, hand}, bid
}

type HandInfo struct {
	hand  Hand
	value int
	bid   int
}

func run(input string) string {
	handInfos := make([]HandInfo, 0)
	for _, line := range strings.Split(input, "\n") {
		hand, bid := getHandAndBid(line)
		value := hand.handValue()
		handInfos = append(handInfos, HandInfo{hand, value, bid})
	}
	sort.Slice(handInfos, func(i, j int) bool {
		if handInfos[i].value < handInfos[j].value {
			return true
		}
		if handInfos[i].value > handInfos[j].value {
			return false
		}
		for n, v := range handInfos[j].hand.deal {
			if handInfos[i].hand.deal[n] < v {
				return true
			}
			if handInfos[i].hand.deal[n] > v {
				return false
			}
		}
		return false
	})
	slices.Reverse(handInfos)
	total := 0
	highestValue := len(handInfos)
	for i, h := range handInfos {
		value := h.bid * (highestValue - i)
		total += value
	}

	fmt.Printf("Total: %d\n", total)

	return fmt.Sprintf("%d", total)
}

func run2(input string) string {
	ValueMap = ValueMap2
	handInfos := make([]HandInfo, 0)
	for _, line := range strings.Split(input, "\n") {
		hand, bid := getHandAndBid(line)
		value := hand.handValue2()
		handInfos = append(handInfos, HandInfo{hand, value, bid})
	}
	sort.Slice(handInfos, func(i, j int) bool {
		if handInfos[i].value < handInfos[j].value {
			return true
		}
		if handInfos[i].value > handInfos[j].value {
			return false
		}
		for n, v := range handInfos[j].hand.deal {
			if handInfos[i].hand.deal[n] < v {
				return true
			}
			if handInfos[i].hand.deal[n] > v {
				return false
			}
		}
		return false
	})
	slices.Reverse(handInfos)
	total := 0
	highestValue := len(handInfos)
	for i, h := range handInfos {
		value := h.bid * (highestValue - i)
		total += value
	}

	fmt.Printf("Total: %d\n", total)

	return fmt.Sprintf("%d", total)
}
