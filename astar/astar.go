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
	Name() string
	Neighbors() []Node
	Cost(string) int
	Heuristic(target string) int
}

type Graph struct {
	root  Node
	graph []Node
}

func getPath(prev map[string]Node, current Node) []Node {
	path := make([]Node, 1)
	path[0] = current
	for p, ok := prev[(current).Name()]; ok; p, ok = prev[(current).Name()] {
		current = p
		path = append([]Node{current}, path...)
	}
	return path
}

type NodeScore struct {
	n          Node
	guessScore int
}

// astar finds a path from start to send.
func Astar(start Node, end Node) []Node {
	// The set of nodes we have visisted, and want to expand from
	seenNodes := make([]Node, 1)
	seenNodes[0] = start

	// A map from the name of a node to the previous node on the best path
	prev := make(map[string]Node)

	// The cheapest score found so far for a node, given its name
	cheapestScore := make(map[string]int)
	cheapestScore[(start).Name()] = 0

	// The current best guess for a score from the start node to the node with the given name
	guessScore := make(map[string]int)
	guessScore[(start).Name()] = (start).Heuristic((start).Name())

	// clock := 0
	// now := time.Now()
	// while we have nodes to expand
	for len(seenNodes) > 0 {

		// nextTime := time.Now()
		// Keep seenNodes sorted by lowest score, pick the first one
		current := seenNodes[0]
		// cursor.Position(0, 0)
		// if clock == 0 {
		// 	cursor.Clear()
		// 	fmt.Printf("Seen Nodes: %10d\n", len(seenNodes))
		// 	fmt.Printf("Guess size: %10d\n", len(guessScore))
		// 	fmt.Printf("Time: %v\n", nextTime.Sub(now))
		// 	now = nextTime
		// }
		// clock++
		// if clock == 1000 {
		// 	clock = 0
		// }

		// if we found the path, stop
		if (current).Name() == (end).Name() {
			return getPath(prev, current)
		}

		// we're expanding from this node, so remove it from the list
		seenNodes = seenNodes[1:]

		neighbors := current.Neighbors()
		for _, n := range neighbors {

			// Our guess is the best score we have plus the cost of traversing the edge
			guess := cheapestScore[current.Name()] + current.Cost((n).Name())

			// If we found a better guess
			cScore, ok := cheapestScore[n.Name()]
			if !ok || guess < cScore {
				nName := n.Name()
				prev[nName] = current
				cheapestScore[nName] = guess
				guessScore[nName] = guess + (start).Heuristic(end.Name())

				found := false
				for i := 0; i < len(seenNodes) && !found; i++ {
					if (seenNodes[i]).Name() == nName {
						found = true
					}
				}

				// if the neighbor isn't already in the set of nodes to expand, add it
				if !found {
					seenNodes = append(seenNodes, n)
					// sort the list by our best guess
					sort.Slice(seenNodes, func(i, j int) bool {
						iScore, iOK := guessScore[(seenNodes[i]).Name()]
						if !iOK {
							iScore = math.MaxInt / 2
						}
						jScore, jOK := guessScore[(seenNodes[j]).Name()]
						if !jOK {
							jScore = math.MaxInt / 2
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
