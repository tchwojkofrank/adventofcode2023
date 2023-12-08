package main

import (
	"fmt"
	"log"
	"os"
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

type Node struct {
	name  string
	left  string
	right string
}

func getNode(line string) Node {
	parts := strings.Split(line, " = ")
	destinationString := strings.Trim(parts[1], "()")
	destinations := strings.Split(destinationString, ", ")
	return Node{
		name:  parts[0],
		left:  destinations[0],
		right: destinations[1],
	}
}

func getNodes(input string) map[string]Node {
	nodes := make(map[string]Node)
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		node := getNode(line)
		nodes[node.name] = node
	}
	return nodes
}

func traverse(directions string, nodes map[string]Node, start string, end string) int {
	current := start
	i := 0
	count := 0
	for strings.Compare(current, end) != 0 {
		node := nodes[current]
		direction := directions[i]
		if direction == 'L' {
			current = node.left
		} else {
			current = node.right
		}
		i = (i + 1) % len(directions)
		count++
	}
	return count
}

func getNodesBySuffix(nodes map[string]Node, suffix string) []string {
	startNodes := make([]string, 0)
	for k, _ := range nodes {
		if strings.HasSuffix(k, suffix) {
			startNodes = append(startNodes, k)
		}
	}
	return startNodes
}

func isEndNode(nodeName string) bool {
	return strings.HasSuffix(nodeName, "Z")
}

func traverseNode(directions string, nodes map[string]Node, start string) (string, int) {
	endNodeName := ""
	endCount := -1
	current := start
	i := 0
	count := 0
	for !isEndNode(current) || count == 0 {
		node := nodes[current]
		direction := directions[i]
		if direction == 'L' {
			current = node.left
		} else {
			current = node.right
		}
		i = (i + 1) % len(directions)
		count = count + 1
	}
	endNodeName = current
	endCount = count
	return endNodeName, endCount
}

func run(input string) string {
	sections := strings.Split(input, "\n\n")
	directions := sections[0]
	nodes := getNodes(sections[1])
	stepCount := traverse(directions, nodes, "AAA", "ZZZ")
	fmt.Printf("Step count: %v\n", stepCount)
	return fmt.Sprintf("%v", stepCount)
}

type NodeResult struct {
	name  string
	count int
}

func gcd(a int, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func lcm(a int, b int) int {
	return (a * b) / gcd(a, b)
}

func lcmSlice(nums []int) int {
	result := nums[0]
	for i := 1; i < len(nums); i++ {
		result = lcm(result, nums[i])
	}
	return result
}

func run2(input string) string {
	sections := strings.Split(input, "\n\n")
	directions := sections[0]
	fmt.Printf("Directions: %v\n", directions)
	nodes := getNodes(sections[1])
	startNodes := getNodesBySuffix(nodes, "A")
	mapResults := make(map[string]NodeResult)
	counts := make([]int, len(startNodes))
	for i, start := range startNodes {
		endNodeName, endCount := traverseNode(directions, nodes, start)
		mapResults[start] = NodeResult{
			name:  endNodeName,
			count: endCount,
		}
		counts[i] = endCount
	}

	lcm := lcmSlice(counts)

	fmt.Printf("Least common multiple: %v\n", lcm)

	return fmt.Sprintf("%v", lcm)
}
