# Day 3: Gear Ratios

[Link to day 3](https://adventofcode.com/2023/day/3)

## Part 1

### The problem
The input is a 2d map of numbers and parts (characters other than digits and `.`).
The goal is to determine which numbers are adjacent to a part and sum them.

### The solution
My solution has some room for optimization, both in time and storage.
I first scanned the map recording the location of the values and parts.
The numbers are stored in a slice of structures that save their value along with their
start and end position (2D points)
The parts are stored as a set. In Go, the easiest implementation of a set is a map where
the value associated with the key is an empty struct
I then iterate through the list of numbers, checking all neighboring positions to the
start and end positions of the number for a part.
So, if the grid is M by N, and there are V values, and P parts, the running time of this
is O(M+N)+O(V)*O(P).
O(M+N) to scan the map
O(V) to loop through the values
O(P) to check the set for each value

I could have inspected the map directly instead of querying the set of parts for a slight
optimization.

## Part 2

### The twist
It turns out parts represented by `*` are gears if they have exactly two adjacent numbers, and we
take the sum of each gear's product of the two values to get the answer.

### The solution
Since I already had the list of values and objects from part 1, I looped through the set of
parts, looking for parts labeled with `*`. Each time I found one, I looped through all values to
determine if they were adjacent, by comparing their start and end points against the gear's point.
If there were exactly two, I multipled their values, and added them to the sum.

Again this is not optimal, and took an additional O(P) + O(G)*O(V) time, where G is the number
of gears.

It does have the benefit of having few edge cases to consider. Because of the usage of the map,
I also didn't have to worry about index values outside the range of the map.
