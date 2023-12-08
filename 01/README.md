# Day 1: Trebuchet?!

[Link to day 1](https://adventofcode.com/2023/day/1)

## Part 1

### The problem
The input is a set of strings, and the goal is to find the first
numeral character and last numeral character in the string, and
combine them into a two digit number.

### The solution
For the first numeral, I scanned the string from the beginning until
I found a character in the range from '0' to '9'.

For the second numeral, I scanned backwards from the last character
in the string until I found a numeral.

## Part 2

### The twist
It turns out the digits may also be spelled out. The same algorithm
as part 1 works, we just need to check for words as well as numerals.
