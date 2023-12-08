package main

import (
	"fmt"
	"log"
	"os"
	"slices"
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
	// run(text)
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

func getStartNodes(nodes map[string]Node) []string {
	startNodes := make([]string, 0)
	for k, _ := range nodes {
		if strings.HasSuffix(k, "A") {
			startNodes = append(startNodes, k)
		}
	}
	return startNodes
}

func isEndStateLengths(counts []int) bool {
	for _, count := range counts {
		if count == 0 {
			return false
		}
	}
	return true
}

func isEndState(nodes []string) bool {
	for _, node := range nodes {
		if !isEndNode(node) {
			return false
		}
	}
	return true
}

func isEndNode(nodeName string) bool {
	return strings.HasSuffix(nodeName, "Z")
}

func traverseParallel(directions string, nodes map[string]Node, startNodes []string) int {
	currentNodes := slices.Clone(startNodes)
	fmt.Printf("Start nodes: %v\n", startNodes)
	i := 0
	count := 0
	for !isEndState(currentNodes) {
		for n := 0; n < len(currentNodes); n++ {
			// once we know how
			node := nodes[currentNodes[n]]
			direction := directions[i]
			if direction == 'L' {
				currentNodes[n] = node.left
			} else {
				currentNodes[n] = node.right
			}
			// if isEndNode(currentNodes[n]) && counts[n] == 0 {
			// 	counts[n] = count
			// }
			// if count > 0 && currentNodes[n] == startNodes[n] {
			// 	cycleLengths[n] = count
			// }
		}
		i = (i + 1) % len(directions)
		count++
		endCount := 0
		for _, node := range currentNodes {
			if isEndNode(node) {
				endCount++
			}
		}
		if endCount > 0 {
			fmt.Printf("%v Current nodes: %v\n", count, currentNodes)
		}
	}
	return count
}

// returns the ending string and the number of steps to get there from the start
// and the last node visited if the whole instruction set were run from the start
func traverseNode(directions string, nodes map[string]Node, start string) (string, int, string, int) {
	endNodeName := ""
	endCount := -1
	lastNodeName := ""
	lastCount := 0
	current := start
	i := 0
	for ; i < len(directions); i++ {
		node := nodes[current]
		direction := directions[i]
		if direction == 'L' {
			current = node.left
		} else {
			current = node.right
		}
		if isEndNode(current) && endCount == -1 {
			endNodeName = current
			endCount = i + 1
		}
	}
	lastNodeName = current
	lastCount = i
	return endNodeName, endCount, lastNodeName, lastCount
}

func run(input string) string {
	sections := strings.Split(input, "\n\n")
	directions := sections[0]
	nodes := getNodes(sections[1])
	stepCount := traverse(directions, nodes, "AAA", "ZZZ")
	fmt.Printf("Step count: %v\n", stepCount)
	return fmt.Sprintf("%v", stepCount)
}

func run2(input string) string {
	sections := strings.Split(input, "\n\n")
	directions := sections[0]
	fmt.Printf("Directions: %v\n", directions)
	nodes := getNodes(sections[1])
	// startNodes := getStartNodes(nodes)
	stepCount := 0
	for k, _ := range nodes {
		endNodeName, endCount, lastNodeName, lastCount := traverseNode(directions, nodes, k)
		fmt.Printf("Start node: %v, end node: %v, end count: %v, last node: %v, last count: %v\n", k, endNodeName, endCount, lastNodeName, lastCount)
	}
	// stepCount = traverseParallel(directions, nodes, startNodes)

	return fmt.Sprintf("%v", stepCount)
}
