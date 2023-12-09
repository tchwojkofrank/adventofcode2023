# [Day 2: Cube Conundrum](https://adventofcode.com/2023/day/2)

## Part 1

### The problem
Each line of the input records multiple rounds of a game: the number and color
of cubes that have been drawn from a bag with red, blue, and green cubes.
The goal is to determine if the game was possible given that the game starts with
12 red, 13 blue, and 14 green cubes in the bag.

### The solution
The solution is straightforward, determining if the number of drawn cubes of
a given color is less than the original number of cubes in a bag.
A loop to process each game, and a loop for each pull and color to compare.
If N is the number of games, then this is O(N).

## Part 2

### The twist
Instead of assuming the number of each color cubes in the bag at the beginning of
each game, the goal is to find the minimum number of cubes possible in the bag to
start each game.

### The solution
Structurally the solution is the same: loop over each round, and save the maximum
value for each color.
Again this is O(N)
