package astar

import (
	"math"
	"sort"
)

// We can run the A* algorithm on a graph of nodes, if:
// 1. we can get a name for the node to hash by
// 2. get all neighbors of a node
// 3. get the cost (or weight) of the edges
// 4. have a heuristic function to guess at the cost to a target node
type Node interface {
	name() string
	neighbors() []*Node
	cost(string) int
	heuristic(target string) int
}

type Graph struct {
	root  *Node
	graph []*Node
}

func getPath(prev map[string]*Node, current *Node) []*Node {
	path := make([]*Node, 1)
	path[0] = current
	for p, ok := prev[(*current).name()]; ok; {
		current = p
		path = append([]*Node{current}, path...)
	}
	return path
}

type NodeScore struct {
	n          *Node
	guessScore int
}

// astar finds a path from start to send.
func astar(start *Node, end *Node) []*Node {
	// The set of nodes we have visisted, and want to expand from
	seenNodes := make([]*Node, 1)
	seenNodes[0] = start

	// A map from the name of a node to the previous node on the best path
	prev := make(map[string]*Node)

	// The cheapest score found so far for a node, given its name
	cheapestScore := make(map[string]int)
	cheapestScore[(*start).name()] = 0

	// The current best guess for a score from the start node to the node with the given name
	guessScore := make(map[string]int)
	guessScore[(*start).name()] = (*start).heuristic((*start).name())

	// while we have nodes to expand
	for len(seenNodes) > 0 {

		// Keep seenNodes sorted by lowest score, pick the first one
		current := seenNodes[0]

		// if we found the path, stop
		if (*current).name() == (*end).name() {
			return getPath(prev, current)
		}

		// we're expanding from this node, so remove it from the list
		seenNodes = seenNodes[1:]

		neighbors := (*current).neighbors()
		for _, n := range neighbors {

			// Our guess is the best score we have plus the cost of traversing the edge
			guess := cheapestScore[(*current).name()] + (*current).cost((*n).name())

			// If we found a better guess
			if guess < cheapestScore[(*n).name()] {
				prev[(*n).name()] = current
				cheapestScore[(*n).name()] = guess
				guessScore[(*n).name()] = guess + (*start).heuristic((*n).name())

				found := false
				for i := 0; i < len(seenNodes) && !found; i++ {
					if (*(seenNodes[i])).name() == (*n).name() {
						found = true
					}
				}

				// if the neighbor isn't already in the set of nodes to expand, add it
				if !found {
					seenNodes = append(seenNodes, n)
					// sort the list by our best guess
					sort.Slice(seenNodes, func(i, j int) bool {
						iScore, iOK := guessScore[(*(seenNodes[i])).name()]
						if !iOK {
							iScore = math.MaxInt
						}
						jScore, jOK := guessScore[(*(seenNodes[j])).name()]
						if !jOK {
							jScore = math.MaxInt
						}
						return iScore < jScore
					})
				}

			}
		}
	}

	// Open set is empty but goal was never reached
	return nil
}
