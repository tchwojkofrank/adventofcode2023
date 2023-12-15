package main

import (
	"fmt"
	"hash/maphash"
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
	run(text)
	end := time.Now()
	fmt.Printf("Running time: %v\n", end.Sub(start))
	start = time.Now()
	run2(text)
	end = time.Now()
	fmt.Printf("Running time: %v\n", end.Sub(start))
}

type Point struct {
	x int
	y int
}

type StoneMap struct {
	m      map[Point]rune
	width  int
	height int
}

func readStoneMap(input string) StoneMap {
	lines := strings.Split(input, "\n")
	slices.Reverse(lines)
	var stoneMap StoneMap
	stoneMap.m = make(map[Point]rune)
	stoneMap.width = len(lines[0])
	stoneMap.height = len(lines)
	for y, line := range lines {
		for x, stone := range line {
			if stone == '#' || stone == 'O' {
				stoneMap.m[Point{x, y}] = stone
			}
		}
	}
	return stoneMap
}

func (stoneMap *StoneMap) canRollInto(p Point) bool {
	if p.x < 0 || p.x >= stoneMap.width || p.y < 0 || p.y >= stoneMap.height {
		return false
	}
	_, ok := stoneMap.m[p]
	return !ok
}

func (stoneMap *StoneMap) tiltRowN(y int) {
	for x := 0; x < stoneMap.width; x++ {
		stone := stoneMap.m[Point{x, y}]
		if stone == 'O' {
			for z := y + 1; z < stoneMap.height && stoneMap.canRollInto(Point{x, z}); z++ {
				delete(stoneMap.m, Point{x, z - 1})
				stoneMap.m[Point{x, z}] = 'O'
			}
		}
	}
}

func (stoneMap *StoneMap) tiltN() {
	for y := stoneMap.height - 1; y >= 0; y-- {
		stoneMap.tiltRowN(y)
	}
}

func (stoneMap *StoneMap) tiltColW(x int) {
	for y := 0; y < stoneMap.height; y++ {
		stone := stoneMap.m[Point{x, y}]
		if stone == 'O' {
			for z := x - 1; z >= 0 && stoneMap.canRollInto(Point{z, y}); z-- {
				delete(stoneMap.m, Point{z + 1, y})
				stoneMap.m[Point{z, y}] = 'O'
			}
		}
	}
}

func (stoneMap *StoneMap) tiltW() {
	for x := 0; x < stoneMap.width; x++ {
		stoneMap.tiltColW(x)
	}
}

func (stoneMap *StoneMap) tiltRowS(y int) {
	for x := 0; x < stoneMap.width; x++ {
		stone := stoneMap.m[Point{x, y}]
		if stone == 'O' {
			for z := y - 1; z >= 0 && stoneMap.canRollInto(Point{x, z}); z-- {
				delete(stoneMap.m, Point{x, z + 1})
				stoneMap.m[Point{x, z}] = 'O'
			}
		}
	}
}

func (stoneMap *StoneMap) tiltS() {
	for y := 0; y < stoneMap.height; y++ {
		stoneMap.tiltRowS(y)
	}
}

func (stoneMap *StoneMap) tiltColE(x int) {
	for y := 0; y < stoneMap.height; y++ {
		stone := stoneMap.m[Point{x, y}]
		if stone == 'O' {
			for z := x + 1; z < stoneMap.width && stoneMap.canRollInto(Point{z, y}); z++ {
				delete(stoneMap.m, Point{z - 1, y})
				stoneMap.m[Point{z, y}] = 'O'
			}
		}
	}
}

func (stoneMap *StoneMap) tiltE() {
	for x := stoneMap.width - 1; x >= 0; x-- {
		stoneMap.tiltColE(x)
	}
}

func (stoneMap *StoneMap) tiltCycle() {
	stoneMap.tiltN()
	stoneMap.tiltW()
	stoneMap.tiltS()
	stoneMap.tiltE()
}

func (stoneMap *StoneMap) tiltCycleRepeat(count int) {
	for i := 0; i < count; i++ {
		stoneMap.tiltCycle()
	}
}

func (stoneMap *StoneMap) load() int {
	load := 0
	for k, v := range stoneMap.m {
		if v == 'O' {
			load += k.y + 1
		}
	}
	return load
}

func (stoneMap StoneMap) String() string {
	result := ""
	for y := stoneMap.height - 1; y >= 0; y-- {
		for x := 0; x < stoneMap.width; x++ {
			r, ok := stoneMap.m[Point{x, y}]
			if !ok {
				r = '.'
			}
			result += string(r)
		}
		trailer := fmt.Sprintf("\t%v\n", y+1)
		result += trailer
	}
	return result
}

func run(input string) string {
	stoneMap := readStoneMap(input)
	stoneMap.tiltN()
	load := stoneMap.load()
	fmt.Printf("Load: %d\n", load)

	return fmt.Sprintf("%d", load)
}

const CycleCount = 1000000000

type BucketID struct {
	index int
	start int
	end   int
}

type BucketDependencies map[BucketID][]BucketID
type BucketValues map[BucketID]int
type VBuckets [][]BucketID
type HBuckets [][]BucketID
type VBucketValues BucketValues
type HBucketValues BucketValues
type VBucketDependencies BucketDependencies
type HBucketDependencies BucketDependencies

func getVBucketsFromInput(lines []string) (VBuckets, VBucketValues) {
	width := len(lines[0])
	vBuckets := make(VBuckets, width)
	vBucketValues := make(VBucketValues)
	for col := 0; col < width; col++ {
		vBuckets[col] = make([]BucketID, 0)
		startedInterval := false
		ballCount := 0
		bucketID := BucketID{index: col}
		for row := 0; row < len(lines); row++ {
			switch lines[row][col] {
			case '#':
				if startedInterval {
					bucketID.end = row - 1
					startedInterval = false
					vBuckets[col] = append(vBuckets[col], bucketID)
					vBucketValues[bucketID] = ballCount
					ballCount = 0
				}
			default:
				if lines[row][col] == 'O' {
					ballCount++
				}
				if !startedInterval {
					bucketID.start = row
					startedInterval = true
				}
			}
		}
		if startedInterval {
			bucketID.end = len(lines) - 1
			vBuckets[col] = append(vBuckets[col], bucketID)
			vBucketValues[bucketID] = ballCount
		}
	}
	return vBuckets, vBucketValues
}

func getHBucketsFromInput(lines []string) (HBuckets, HBucketValues) {
	height := len(lines)
	hBuckets := make(HBuckets, height)
	hBucketValues := make(HBucketValues)
	for row := 0; row < height; row++ {
		hBuckets[row] = make([]BucketID, 0)
		startedInterval := false
		ballCount := 0
		bucketID := BucketID{index: row}
		for col := 0; col < len(lines[0]); col++ {
			switch lines[row][col] {
			case '#':
				if startedInterval {
					bucketID.end = col - 1
					startedInterval = false
					hBuckets[row] = append(hBuckets[row], bucketID)
					hBucketValues[bucketID] = ballCount
					ballCount = 0
				}
			default:
				if lines[row][col] == 'O' {
					ballCount++
				}
				if !startedInterval {
					bucketID.start = col
					startedInterval = true
				}
			}
		}
		if startedInterval {
			bucketID.end = len(lines[0]) - 1
			hBuckets[row] = append(hBuckets[row], bucketID)
			hBucketValues[bucketID] = ballCount
		}
	}
	return hBuckets, hBucketValues
}

type BucketInfo struct {
	vBuckets      VBuckets
	vBucketValues VBucketValues
	hBuckets      HBuckets
	hBucketValues HBucketValues
}

func setBucketDependencies(bucketInfo BucketInfo) (VBucketDependencies, HBucketDependencies) {
	vBD := make(VBucketDependencies)
	hBD := make(HBucketDependencies)
	// start by setting up vertical dependencies on horizontal buckets
	// a vertical bucket depends on all horizontal buckets that overlap with it
	for col := 0; col < len(bucketInfo.vBuckets); col++ {
		for _, vBucket := range bucketInfo.vBuckets[col] {
			for row := vBucket.start; row <= vBucket.end; row++ {
				for _, hBucket := range bucketInfo.hBuckets[row] {
					if hBucket.start <= col && hBucket.end >= col {
						vBD[vBucket] = append(vBD[vBucket], hBucket)
					}
				}
			}
		}
	}
	// now set up horizontal dependencies on vertical buckets
	// a horizontal bucket depends on all vertical buckets that overlap with it
	for row := 0; row < len(bucketInfo.hBuckets); row++ {
		for _, hBucket := range bucketInfo.hBuckets[row] {
			for col := hBucket.start; col <= hBucket.end; col++ {
				for _, vBucket := range bucketInfo.vBuckets[col] {
					if vBucket.start <= row && vBucket.end >= row {
						hBD[hBucket] = append(hBD[hBucket], vBucket)
					}
				}
			}
		}
	}

	return vBD, hBD
}

// Tilt everything to the north
// For each horizontal bucket, check if any balls roll into it from the south
// For each column in the horizontal bucket, a ball rolls into that slot if
// the dependent vertical buckets in that column has enough balls in it to reach the row
func tiltNorth(bucketInfo *BucketInfo, hBucketDependencies HBucketDependencies) {
	for row := 0; row < len(bucketInfo.hBuckets); row++ {
		for _, hBucket := range bucketInfo.hBuckets[row] {
			bucketInfo.hBucketValues[hBucket] = 0
			for _, dependency := range hBucketDependencies[hBucket] {
				requiredBallCount := row - dependency.start + 1
				if requiredBallCount <= bucketInfo.vBucketValues[dependency] {
					bucketInfo.hBucketValues[hBucket] += 1
				}
			}
		}
	}
}

func tiltWest(bucketInfo *BucketInfo, vBucketDependencies VBucketDependencies) {
	for col := 0; col < len(bucketInfo.vBuckets); col++ {
		for _, vBucket := range bucketInfo.vBuckets[col] {
			bucketInfo.vBucketValues[vBucket] = 0
			for _, dependency := range vBucketDependencies[vBucket] {
				requiredBallCount := col - dependency.start + 1
				if requiredBallCount <= bucketInfo.hBucketValues[dependency] {
					bucketInfo.vBucketValues[vBucket] += 1
				}
			}
		}
	}
}

func tiltSouth(bucketInfo *BucketInfo, hBucketDependencies HBucketDependencies) {
	for row := len(bucketInfo.hBuckets) - 1; row >= 0; row-- {
		for _, hBucket := range bucketInfo.hBuckets[row] {
			bucketInfo.hBucketValues[hBucket] = 0
			for _, dependency := range hBucketDependencies[hBucket] {
				requiredBallCount := dependency.end - row + 1
				if requiredBallCount <= bucketInfo.vBucketValues[dependency] {
					bucketInfo.hBucketValues[hBucket] += 1
				}
			}
		}
	}
}

func tiltEast(bucketInfo *BucketInfo, vBucketDependencies VBucketDependencies) {
	for col := len(bucketInfo.vBuckets) - 1; col >= 0; col-- {
		for _, vBucket := range bucketInfo.vBuckets[col] {
			bucketInfo.vBucketValues[vBucket] = 0
			for _, dependency := range vBucketDependencies[vBucket] {
				requiredBallCount := dependency.end - col + 1
				if requiredBallCount <= bucketInfo.hBucketValues[dependency] {
					bucketInfo.vBucketValues[vBucket] += 1
				}
			}
		}
	}
}

func (bucketInfo *BucketInfo) tiltCycle(hBucketDependencies HBucketDependencies, vBucketDependencies VBucketDependencies) {
	tiltNorth(bucketInfo, hBucketDependencies)
	tiltWest(bucketInfo, vBucketDependencies)
	tiltSouth(bucketInfo, hBucketDependencies)
	tiltEast(bucketInfo, vBucketDependencies)
}

func (bi BucketInfo) String() string {
	result := ""
	hBs := bi.hBuckets
	hV := bi.hBucketValues
	result += "Horizontal buckets:\n"
	for _, hBucket := range hBs {
		for _, b := range hBucket {
			result += fmt.Sprintf("H%v: %d\n", b, hV[b])
		}
	}

	fmt.Println()

	vBs := bi.vBuckets
	vV := bi.vBucketValues
	result += "Vertical buckets:\n"
	for _, vBucket := range vBs {
		for _, b := range vBucket {
			result += fmt.Sprintf("V%v: %d\n", b, vV[b])
		}
	}

	fmt.Print("\n")

	return result
}

func (bi BucketInfo) load() int {
	load := 0
	for hID, hv := range bi.hBucketValues {
		row := hID.index
		value := len(bi.hBuckets) - row
		load += value * hv
	}
	return load
}

func run2(input string) string {
	load := 0

	lines := strings.Split(input, "\n")
	bucketInfo := BucketInfo{}
	bucketInfo.vBuckets, bucketInfo.vBucketValues = getVBucketsFromInput(lines)
	bucketInfo.hBuckets, bucketInfo.hBucketValues = getHBucketsFromInput(lines)
	vBucketDependencies, hBucketDependencies := setBucketDependencies(bucketInfo)

	var hash maphash.Hash
	const cacheSize = 100000
	vcache := make([]uint64, 0, cacheSize)
	hcache := make([]uint64, 0, cacheSize)

	for i := 0; i < CycleCount; i++ {
		bucketInfo.tiltCycle(hBucketDependencies, vBucketDependencies)
		// check if we have a cycle
		hash.Reset()
		hash.Write([]byte(fmt.Sprintf("%v", bucketInfo.vBucketValues)))
		vhash := hash.Sum64()
		hash.Reset()
		hash.Write([]byte(fmt.Sprintf("%v", bucketInfo.hBucketValues)))
		hhash := hash.Sum64()
		for j := 0; j < len(vcache); j++ {
			if vcache[j] == vhash && hcache[j] == hhash {
				// calculate the remaining cycles
				remainingCycles := (CycleCount - i) % (i - j)
				i = CycleCount - remainingCycles
				break
			}
		}
		if (i % cacheSize) >= len(vcache) {
			vcache = append(vcache, vhash)
			hcache = append(hcache, hhash)
		} else {
			vcache[i%cacheSize] = vhash
			hcache[i%cacheSize] = hhash
		}
	}
	load = bucketInfo.load()
	fmt.Printf("Load: %d\n", load)

	return fmt.Sprintf("%d", load)
}
