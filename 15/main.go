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

func hash(input string) int {
	value := 0
	for _, r := range input {
		value += int(r)
		value *= 17
		value = value % 256
	}
	return value
}

func run(input string) string {
	instructions := strings.Split(input, ",")
	sum := 0
	for _, instruction := range instructions {
		value := hash(instruction)
		sum += value
	}
	fmt.Printf("Sum: %v\n", sum)
	return fmt.Sprintf("%v", sum)
}

type Lens struct {
	label string
	lens  int
}

func run2(input string) string {
	instructions := strings.Split(input, ",")
	boxes := make([][]Lens, 256)
	for i := 0; i < 256; i++ {
		boxes[i] = make([]Lens, 0)
	}

	for _, instruction := range instructions {
		if strings.Contains(instruction, "-") {
			label := strings.TrimSuffix(instruction, "-")
			boxIndex := hash(label)
			for i, lens := range boxes[boxIndex] {
				if lens.label == label {
					boxes[boxIndex] = append(boxes[boxIndex][:i], boxes[boxIndex][i+1:]...)
					break
				}
			}
		} else if strings.Contains(instruction, "=") {
			fields := strings.Split(instruction, "=")
			label := fields[0]
			boxIndex := hash(label)
			focalLength, _ := strconv.Atoi(fields[1])
			var i int
			var lens Lens
			for i = 0; i < len(boxes[boxIndex]); i++ {
				lens = boxes[boxIndex][i]
				if lens.label == label {
					boxes[boxIndex][i].lens = focalLength
					break
				}
			}
			if i == len(boxes[boxIndex]) {
				boxes[boxIndex] = append(boxes[boxIndex], Lens{label: label, lens: focalLength})
			}
		}
	}
	sum := 0
	for i, box := range boxes {
		for slot, lens := range box {
			focalPower := (i + 1) * (slot + 1) * lens.lens
			sum += focalPower
		}
	}
	fmt.Printf("Sum: %v\n", sum)
	return fmt.Sprintf("%v", sum)
}
