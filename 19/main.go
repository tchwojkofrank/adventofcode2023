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

type Condition struct {
	attribute string
	operator  string
	value     int
}

type Case struct {
	condition Condition
	next      string
}

type Workflow struct {
	name        string
	cases       []Case
	defaultCase string
}

var pattern *regexp.Regexp

func getWorkflow(s string) Workflow {
	workflow := Workflow{}
	var rulesString string
	workflow.name, rulesString, _ = strings.Cut(s, "{")
	rulesString = strings.Trim(rulesString, "}")
	ruleStrings := strings.Split(rulesString, ",")
	workflow.cases = make([]Case, len(ruleStrings)-1)
	for i := 0; i < len(ruleStrings)-1; i++ {
		ruleString := ruleStrings[i]
		parts := pattern.FindStringSubmatch(ruleString)
		condition := Condition{parts[1], parts[2], 0}
		condition.value, _ = strconv.Atoi(parts[3])
		workflow.cases[i] = Case{condition, parts[4]}
	}
	workflow.defaultCase = ruleStrings[len(ruleStrings)-1]
	return workflow
}

type Object struct {
	x      int
	m      int
	a      int
	s      int
	rating int
}

func getObject(s string) Object {
	object := Object{}
	s = strings.TrimPrefix(s, "{")
	s = strings.TrimSuffix(s, "}")
	parts := strings.Split(s, ",")
	object.x, _ = strconv.Atoi(parts[0][2:])
	object.m, _ = strconv.Atoi(parts[1][2:])
	object.a, _ = strconv.Atoi(parts[2][2:])
	object.s, _ = strconv.Atoi(parts[3][2:])
	object.rating = object.x + object.m + object.a + object.s
	return object
}

func isAccepted(first string, workflows map[string]Workflow, object Object) bool {
	workflow := workflows[first]
	decision := false
	for !decision {
		result := ""
		for _, c := range workflow.cases {
			var attributeValue int
			match := false
			switch c.condition.attribute {
			case "x":
				attributeValue = object.x
			case "m":
				attributeValue = object.m
			case "a":
				attributeValue = object.a
			case "s":
				attributeValue = object.s
			}
			switch c.condition.operator {
			case "<":
				if attributeValue < c.condition.value {
					result = c.next
					match = true
				}
			case ">":
				if attributeValue > c.condition.value {
					result = c.next
					match = true
				}
			}
			if match {
				break
			}
		}
		if result == "" {
			result = workflow.defaultCase
		}
		switch result {
		case "R":
			return false
		case "A":
			return true
		default:
			workflow = workflows[result]
		}
	}
	return false
}

func run(input string) string {
	pattern = regexp.MustCompile(`([a-z])([<>])([0-9]+):([a-zAR]+)`)
	sections := strings.Split(input, "\n\n")
	workflowStrings := strings.Split(sections[0], "\n")
	workflows := make(map[string]Workflow)
	for _, workflowString := range workflowStrings {
		workflow := getWorkflow(workflowString)
		workflows[workflow.name] = workflow
	}
	objectStrings := strings.Split(sections[1], "\n")
	objects := make([]Object, len(objectStrings))
	for i, objectString := range objectStrings {
		objects[i] = getObject(objectString)
	}
	totalRating := 0
	for _, object := range objects {
		if isAccepted("in", workflows, object) {
			totalRating += object.rating
		}
	}

	fmt.Printf("Total rating: %d\n", totalRating)
	return fmt.Sprintf("%d", totalRating)
}

type Range struct {
	min int
	max int
}

type RatingRange struct {
	x Range
	m Range
	a Range
	s Range
}

func updateRatingRange(c Case, ratingRange RatingRange) RatingRange {
	switch c.condition.attribute {
	case "x":
		switch c.condition.operator {
		case "<":
			ratingRange.x.max = min(c.condition.value-1, ratingRange.x.max)
		case ">":
			ratingRange.x.min = max(c.condition.value+1, ratingRange.x.min)
		}
	case "m":
		switch c.condition.operator {
		case "<":
			ratingRange.m.max = min(c.condition.value-1, ratingRange.m.max)
		case ">":
			ratingRange.m.min = max(c.condition.value+1, ratingRange.m.min)
		}
	case "a":
		switch c.condition.operator {
		case "<":
			ratingRange.a.max = min(c.condition.value-1, ratingRange.a.max)
		case ">":
			ratingRange.a.min = max(c.condition.value+1, ratingRange.a.min)
		}
	case "s":
		switch c.condition.operator {
		case "<":
			ratingRange.s.max = min(c.condition.value-1, ratingRange.s.max)
		case ">":
			ratingRange.s.min = max(c.condition.value+1, ratingRange.s.min)
		}
	}
	return ratingRange
}

func getRatingRanges(currentWF string, workflows map[string]Workflow, ratingRange RatingRange) []RatingRange {
	ratingRanges := make([]RatingRange, 0)
	for _, c := range workflows[currentWF].cases {
		newRange := updateRatingRange(c, ratingRange)
		if c.next == "A" {
			ratingRanges = append(ratingRanges, newRange)
		} else if c.next != "R" {
			ratingRanges = append(ratingRanges, getRatingRanges(c.next, workflows, newRange)...)
		}
	}
	if workflows[currentWF].defaultCase == "A" {
		ratingRanges = append(ratingRanges, ratingRange)
	} else if workflows[currentWF].defaultCase != "R" {
		ratingRanges = append(ratingRanges, getRatingRanges(workflows[currentWF].defaultCase, workflows, ratingRange)...)
	}
	return ratingRanges
}

type cube4d struct {
	sides [4]Range
}

func (c cube4d) volume() int {
	return (c.sides[0].max - c.sides[0].min + 1) * (c.sides[1].max - c.sides[1].min + 1) * (c.sides[2].max - c.sides[2].min + 1) * (c.sides[3].max - c.sides[3].min + 1)
}

func (c cube4d) intersects(c2 cube4d) bool {
	for i := 0; i < 4; i++ {
		if c.sides[i].min > c2.sides[i].max || c.sides[i].max < c2.sides[i].min {
			return false
		}
	}
	return true
}

func (c cube4d) intersection(c2 cube4d) cube4d {
	c3 := cube4d{}
	for i := 0; i < 4; i++ {
		c3.sides[i].min = max(c.sides[i].min, c2.sides[i].min)
		c3.sides[i].max = min(c.sides[i].max, c2.sides[i].max)
	}
	return c3
}

// return true if c2 is completely contained in c1
func (c cube4d) contains(c2 cube4d) bool {
	for i := 0; i < 4; i++ {
		if c2.sides[i].min < c.sides[i].min || c2.sides[i].max > c.sides[i].max {
			return false
		}
	}
	return true
}

func ratingsRangeToCube(ratingRange RatingRange) cube4d {
	return cube4d{[4]Range{ratingRange.x, ratingRange.m, ratingRange.a, ratingRange.s}}
}

// break a cube into smaller cubes that do not intersect the other cubes
func breakCube(cube cube4d, newCube cube4d) []cube4d {
	newCubes := make([]cube4d, 0)
	newCubes = append(newCubes, cube)
	if !cube.intersects(newCube) {
		newCubes = append(newCubes, newCube)
		return newCubes
	}
	for i := 0; i < 4; i++ {
		if cube.sides[i].min < newCube.sides[i].min {
			newCube1 := cube
			newCube1.sides[i].max = newCube.sides[i].min - 1
			newCubes = append(newCubes, breakCube(newCube1, newCube)...)
		}
		if cube.sides[i].max > newCube.sides[i].max {
			newCube1 := cube
			newCube1.sides[i].min = newCube.sides[i].max + 1
			newCubes = append(newCubes, breakCube(newCube1, newCube)...)
		}
	}

	return newCubes
}

// break all cubes into smaller cubes that do not intersect the other cubes
func breakCubes(cubes []cube4d) []cube4d {
	newCubes := make([]cube4d, 0)
	for i := 0; i < len(cubes); i++ {
		cube1 := cubes[i]
		for j := i + 1; j < len(cubes); j++ {
			cube2 := cubes[j]
			newCubes = append(newCubes, breakCube(cube1, cube2)...)
		}
	}
	return newCubes
}

// remove any cubes that are completely contained in other cubes
func reduceCubes(cubes []cube4d) []cube4d {
	newCubes := make([]cube4d, 0)
	for i := 0; i < len(cubes); i++ {
		contained := false
		for j := 0; j < len(cubes); j++ {
			if i != j && cubes[j].contains(cubes[i]) {
				contained = true
				break
			}
		}
		if !contained {
			newCubes = append(newCubes, cubes[i])
		}
	}
	return newCubes
}

func getTotalVolume(cubes []cube4d) int {
	totalVolume := 0
	cubes = reduceCubes(cubes)
	cubes = breakCubes(cubes)
	for _, cube := range cubes {
		totalVolume += cube.volume()
	}
	return totalVolume
}

func run2(input string) string {
	pattern = regexp.MustCompile(`([a-z])([<>])([0-9]+):([a-zAR]+)`)
	sections := strings.Split(input, "\n\n")
	workflowStrings := strings.Split(sections[0], "\n")
	workflows := make(map[string]Workflow)
	for _, workflowString := range workflowStrings {
		workflow := getWorkflow(workflowString)
		workflows[workflow.name] = workflow
	}
	ratingRanges := getRatingRanges("in", workflows, RatingRange{Range{1, 4000}, Range{1, 4000}, Range{1, 4000}, Range{1, 4000}})
	cubes := make([]cube4d, len(ratingRanges))
	for i, ratingRange := range ratingRanges {
		cubes[i] = ratingsRangeToCube(ratingRange)
	}
	cubes = reduceCubes(cubes)
	totalVolume := getTotalVolume(cubes)
	fmt.Printf("Total volume: %d\n", totalVolume)
	return fmt.Sprintf("%d", totalVolume)
}
