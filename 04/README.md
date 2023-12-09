# [Day 4: Scratchcards](https://adventofcode.com/2023/day/4)

## Part 1

### The problem

Each line of the input is a set of drawn numbers and a scratchcard.
If n = the number of scratchcard numbers that match drawn numbers,
the scratchcard is worth 2^n points. (The first match is worth one
point, and the points double for each other match.) The answer is
the sum of the points.

### The solution

I used a set for the drawn numbers, and a slice for the scratchcard
numbers. This let me easily check if any given number was drawn.
Set the initial score to zero. Then a loop through the slice of scratchcard
numbers, checking each one against the drawn numbers, setting the value
to 1 for the first match, and doubling each one after. Slightly more
efficient (depending on how the compiler optimized it) would be to
just count the matching numbers and set the value to `1 << count-1`

If S is the number of scratchcards, and N is the count of numbers on each
scratchcard, then the running time is `O(S*N)`. (For each scratchcard check
each number.) Storage is also `O(S*N)`.

## Part 2

### The twist

It turns out the value of the scratchcards was misunderstood. The
number of matching numbers is the number of copies of other scratchcards
you win. So if scratchcard 3 had two matching numbers, you would win
copies of scratchcards 4 and 5. The goal is to figure out how many
scratchcards you would win.

### The solution

I created a slice that stores the number of copies of each scratchcard,
and initialized the entire slice to 1. For each scratchcard,
find the number of matches on the card (`matches`) and increment the number of copies
of the next `matches` cards by the number of copies of this card.
Finally sum the number of copies stored in the slice.

The running time is still `O(S*N)` since we're still processing each unique card
once.

