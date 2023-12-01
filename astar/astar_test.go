package astar

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

/*
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
*/

const (
	maxX = 9
	maxY = 9
)

type TestNode struct {
	x, y int
}

func (n *TestNode) name() string {
	return fmt.Sprintf("%d,%d", n.x, n.y)
}

func (n *TestNode) neighbors() []Node {
	if n.x == 0 && n.y == 0 {
		return []Node{
			&TestNode{1, 0},
			&TestNode{0, 1},
		}
	}
	if n.x == 0 {
		return []Node{
			&TestNode{0, n.y + 1},
			&TestNode{0, n.y - 1},
			&TestNode{1, n.y},
		}
	}
	if n.y == 0 {
		return []Node{
			&TestNode{n.x + 1, 0},
			&TestNode{n.x - 1, 0},
			&TestNode{n.x, 1},
		}
	}
	if n.x == maxX && n.y == maxY {
		return []Node{
			&TestNode{maxX - 1, maxY},
			&TestNode{maxX, maxY - 1},
		}
	}
	if n.x == maxX {
		return []Node{
			&TestNode{maxX, n.y + 1},
			&TestNode{maxX, n.y - 1},
			&TestNode{maxX - 1, n.y},
		}
	}
	if n.y == maxY {
		return []Node{
			&TestNode{n.x + 1, maxY},
			&TestNode{n.x - 1, maxY},
			&TestNode{n.x, maxY - 1},
		}
	}

	return []Node{
		&TestNode{n.x + 1, n.y},
		&TestNode{n.x - 1, n.y},
		&TestNode{n.x, n.y + 1},
		&TestNode{n.x, n.y - 1},
	}
}

func getXYfromName(name string) (int, int) {
	// name is of the form "x,y"
	// so we want to split on the comma
	// and convert the strings to ints
	xy := strings.Split(name, ",")
	x, _ := strconv.Atoi(xy[0])
	y, _ := strconv.Atoi(xy[1])
	return x, y
}

// this is what the nodes look like
/*
  0 1 2 3 4 5 6 7 8 9
0 . . . . . . . # . .
1 . . . . . . . # . .
2 . . . . . . . # . .
3 .	. . . . . . # . .
4 . . . . . . . # . .
5 . . . . . . . . . .
6 . . . . . . . # . .
7 . . . . . . . # . .
8 # # # # # # # # . .
9 . . . . . . . . . .
*/

func (n *TestNode) cost(target string) int {
	x, y := getXYfromName(target)
	startX, startY := getXYfromName(n.name())
	// make a wall from 7,0 to 7,4
	// and from 7,6 to 7,8
	// and from 0,8 to 6,8
	if (x == 7 && y >= 0 && y <= 4) || (x == 7 && y >= 6 && y <= 8) || (y == 8 && x >= 0 && x <= 6) {
		return 10000
	}
	return abs(x-startX) + abs(y-startY)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (n *TestNode) heuristic(target string) int {
	x, y := getXYfromName(target)
	startX, startY := getXYfromName(n.name())
	return abs(x-startX) + abs(y-startY)
}

func (n *TestNode) String() string {
	return fmt.Sprintf("%d,%d", n.x, n.y)
}

// TestAstar tests the A* algorithm
func TestAstar(t *testing.T) {
	// Create a graph
	// 1. Create nodes
	// 2. Connect nodes
	// 3. Create a graph
	// 4. Run A* on the graph
	// 5. Check that the path is correct

	// Create nodes from 0,0 to maxX,maxY
	testNodes := make([]Node, 0)
	var start Node
	var end Node
	for x := 0; x <= maxX; x++ {
		for y := 0; y <= maxY; y++ {
			testNode := &TestNode{x, y}
			if x == 0 && y == 0 {
				start = testNode
			}
			if x == maxX && y == maxY {
				end = testNode
			}
			testNodes = append(testNodes, &TestNode{x, y})
		}
	}

	// Connect nodes
	for _, n := range testNodes {
		(n).neighbors()
	}

	// Run A* on the graph
	path := astar(start, end)

	fmt.Printf("Path: %v\n", path)

	// Check that the path is correct
	if len(path) != maxX+maxY+1 {
		t.Errorf("Path length is incorrect: %d", len(path))
	}

}
