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

func run2(input string) string {
	pattern = regexp.MustCompile(`([a-z])([<>])([0-9]+):([a-zAR]+)`)
	sections := strings.Split(input, "\n\n")
	workflowStrings := strings.Split(sections[0], "\n")
	workflows := make(map[string]Workflow)
	for _, workflowString := range workflowStrings {
		workflow := getWorkflow(workflowString)
		workflows[workflow.name] = workflow
	}
	return ""
}
