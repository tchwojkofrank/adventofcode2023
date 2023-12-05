package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"sort"
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

func getSeeds(input string) []int {
	input = strings.TrimPrefix(input, "seeds: ")
	input = strings.TrimSpace(input)
	seedStrings := strings.Split(input, " ")
	seeds := make([]int, len(seedStrings))
	for i, seedString := range seedStrings {
		seed, _ := strconv.Atoi(seedString)
		seeds[i] = seed
	}
	return seeds
}

type Range struct {
	start int
	end   int
}

type RangeTransform struct {
	source Range
	offset int
}

func getSeeds2(input string) []Range {
	input = strings.TrimPrefix(input, "seeds: ")
	input = strings.TrimSpace(input)
	seedStrings := strings.Split(input, " ")
	seeds := make([]Range, len(seedStrings)/2)
	for i := 0; i < len(seedStrings); i += 2 {
		startString := seedStrings[i]
		lengthString := seedStrings[i+1]
		start, _ := strconv.Atoi(startString)
		length, _ := strconv.Atoi(lengthString)
		seeds[i/2] = Range{start, start + length - 1}
	}
	return seeds
}

func getRangeTransform(line string) RangeTransform {
	line = strings.TrimSpace(line)
	rangeStrings := strings.Split(line, " ")
	startValue, _ := strconv.Atoi(rangeStrings[0])
	startKey, _ := strconv.Atoi(rangeStrings[1])
	length, _ := strconv.Atoi(rangeStrings[2])
	return RangeTransform{Range{startKey, startKey + length - 1}, startValue - startKey}
}

func mapSourceToDest(rangeMaps []RangeTransform, source int) int {
	for _, rangeMap := range rangeMaps {
		if source >= rangeMap.source.start && source <= rangeMap.source.end {
			return source + rangeMap.offset
		}
	}
	return source
}

type mapNameType struct {
	key   string
	value string
}

func createMapFromSection(section string, rangeMap map[string][]RangeTransform, mapMap map[string]string) (map[string][]RangeTransform, map[string]string) {
	lines := strings.Split(section, "\n")
	mapNameString := strings.TrimSuffix(lines[0], " map:")
	mapNameStrings := strings.Split(mapNameString, "-to-")
	mapName := mapNameType{mapNameStrings[0], mapNameStrings[1]}
	lines = lines[1:]
	ranges := make([]RangeTransform, 0)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		ranges = append(ranges, getRangeTransform(line))
	}
	// sort the ranges by key
	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].source.start < ranges[j].source.start
	})
	rangeMap[mapName.key] = ranges
	mapMap[mapName.key] = mapName.value
	return rangeMap, mapMap
}

func mapSeedToLocation(seed int, rangeMap map[string][]RangeTransform, mapMap map[string]string) int {
	sourceMapName := "seed"
	destMapName := mapMap[sourceMapName]
	source := seed
	var dest int
	for sourceMapName != "location" {
		dest = mapSourceToDest(rangeMap[sourceMapName], source)
		sourceMapName = destMapName
		destMapName = mapMap[sourceMapName]
		source = dest
	}

	return dest
}

func (r *Range) intersects(a Range) bool {
	return (a.start >= r.start && a.start <= r.end) ||
		(a.end >= r.start && a.end <= r.end)
}

func (r *Range) contains(a Range) bool {
	return a.start >= r.start && a.end <= r.end
}

func (rt *RangeTransform) transform(a Range) Range {
	if rt.source.contains(a) {
		return Range{a.start + rt.offset, a.end + rt.offset}
	}
	log.Fatal("Range a not contained in the transform range")
	return Range{0, 0}
}

func (r *Range) split(a Range) []Range {
	if r.contains(a) {
		return []Range{a}
	}
	if a.contains(*r) {
		r1 := Range{a.start, r.start - 1}
		r2 := Range{r.start, r.end}
		r3 := Range{r.end + 1, a.end}
		return []Range{r1, r2, r3}
	}
	if r.intersects(a) {
		if a.start >= r.start {
			r1 := Range{a.start, r.end}
			r2 := Range{r.end + 1, a.end}
			return []Range{r1, r2}
		} else {
			r1 := Range{a.start, r.start - 1}
			r2 := Range{r.start, a.end}
			return []Range{r1, r2}
		}
	} else {
		return []Range{a}
	}
}

func normalize(rangeMaps []RangeTransform, sources []Range) []Range {
	destRanges := make([]Range, 0)
	sourceRanges := sources
	foundSplit := true
	for len(sourceRanges) > 0 {
		source := sourceRanges[0]
		foundSplit = false
		// if the range is contained in a rangeMap, or does not intersect any rangeMap, add it to the destRanges
		for _, rangeMap := range rangeMaps {
			if rangeMap.source.contains(source) {
				destRanges = append(destRanges, source)
				sourceRanges = sourceRanges[1:]
				foundSplit = true
				break
			}
			if rangeMap.source.intersects(source) && !rangeMap.source.contains(source) {
				sourceRanges = append(sourceRanges, rangeMap.source.split(source)...)
				sourceRanges = sourceRanges[1:]
				foundSplit = true
			}
		}
		if !foundSplit {
			destRanges = append(destRanges, source)
			sourceRanges = sourceRanges[1:]
		}
	}
	return destRanges
}

func mapSourceRangeToDestRanges(sourceRanges []Range, rangeMaps []RangeTransform) []Range {
	// the rangeMaps are sorted by startKey (probably not needed)
	// first split the source ranges into ranges that are completely contained in a single rangeMap
	normalizedSourceRanges := normalize(rangeMaps, sourceRanges)
	transformedRanges := make([]Range, 0)

	// transform all the ranges
	for _, nr := range normalizedSourceRanges {
		found := false
		for _, rangeMap := range rangeMaps {
			if rangeMap.source.contains(nr) {
				transformedRanges = append(transformedRanges, rangeMap.transform(nr))
				found = true
				break
			}
		}
		if !found {
			transformedRanges = append(transformedRanges, nr)
		}
	}
	// sort the transformed ranges by start value
	sort.Slice(transformedRanges, func(i, j int) bool {
		return transformedRanges[i].start < transformedRanges[j].start
	})

	return transformedRanges
}

func run(input string) string {
	sections := strings.Split(input, "\n\n")
	seeds := getSeeds(sections[0])
	fmt.Printf("Seeds: %v\n", seeds)
	transformMap := make(map[string][]RangeTransform)
	mapMap := make(map[string]string)

	for _, section := range sections[1:] {
		transformMap, mapMap = createMapFromSection(section, transformMap, mapMap)
	}

	seedLocationMap := make(map[int]int)
	minLocation := math.MaxInt
	for _, seed := range seeds {
		seedLocationMap[seed] = mapSeedToLocation(seed, transformMap, mapMap)
		if seedLocationMap[seed] < minLocation {
			minLocation = seedLocationMap[seed]
		}
	}

	fmt.Printf("Min location: %v\n", minLocation)
	return fmt.Sprintf("%v", minLocation)
}

func run2(input string) string {
	sections := strings.Split(input, "\n\n")

	// seed ranges are pairs of integers. The first is the offset, the second is the length
	seedRanges := getSeeds2(sections[0])
	fmt.Printf("Seeds: %v\n", seedRanges)
	transformMap := make(map[string][]RangeTransform)
	mapMap := make(map[string]string)

	for _, section := range sections[1:] {
		transformMap, mapMap = createMapFromSection(section, transformMap, mapMap)
	}

	// seedLocationMap := make(map[int]int)
	sourceMap := "seed"
	destMap := mapMap[sourceMap]
	var destRanges []Range
	sourceRanges := seedRanges
	sort.Slice(sourceRanges, func(i, j int) bool {
		return seedRanges[i].start < seedRanges[j].start
	})
	for sourceMap != "location" {
		destRanges = mapSourceRangeToDestRanges(sourceRanges, transformMap[sourceMap])
		sourceMap = destMap
		destMap = mapMap[sourceMap]
		sourceRanges = destRanges
		sort.Slice(sourceRanges, func(i, j int) bool {
			return destRanges[i].start < destRanges[j].start
		})
		destRanges = make([]Range, 0)
	}

	fmt.Printf("Min location: %v\n", sourceRanges[0].start)
	return fmt.Sprintf("%v", sourceRanges[0].start)
}
