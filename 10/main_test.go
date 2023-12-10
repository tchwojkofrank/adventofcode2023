package main

import (
	"fmt"
	"testing"
	"time"
)

func TestPuzzle1(t *testing.T) {
	start := time.Now()
	text := readInput("test")
	answer := run(text)
	end := time.Now()
	fmt.Printf("Running time: %v\n", end.Sub(start))
	testAnswer := readInput("testAnswer")
	if answer != testAnswer {
		t.Errorf("Expected %v, got %v\n", testAnswer, answer)
	}
}

func TestPuzzle2(t *testing.T) {
	start := time.Now()
	text := readInput("test2")
	answer := run2(text)
	end := time.Now()
	fmt.Printf("Running time: %v\n", end.Sub(start))
	testAnswer := readInput("testAnswer2")
	if answer != testAnswer {
		t.Errorf("Expected %v, got %v\n", testAnswer, answer)
	}
}
