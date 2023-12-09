# [Day 5: If You Give A Seed A Fertilizer!(https://adventofcode.com/2023/day/5)

## Part 1

### The problem

We are given a list of seeds (represented as numbers) and a list of maps. Each
map is defined by a set of transforms defined by destination values, source values,
and counts. So a map with values `D S C` would map the range `[S,S+C-1]` to `[D,D+C-1]`.
If a source value doesn't fall within any of the transform's ranges, then the value is
left alone.

Each map also describes a `from` and `to` name. So `seed` to `soil` or `fertilizer`
to `water`. The goal is to start with the `seed` value and map it to a `location`
value, and find the lowest `location` value.

### The solution

I created a `Range`, which has a `start` and `end`, and a `RangeTransform` which is
a `Range` paired with an `offset`. Now I can build each map as a slice of `RangeTransform`.
Starting with the seed value, I apply each map successively, until I find the `location`.
I keep track of the lowest `location` value, and updated it if I find a new smaller `location`.
I didn't need to, but I kept a slice of all the seed's locations.

Applying each map requires finding which transform's `Range`, if any, applies, and add
that transform's offset to the source value. This search is `O(T)` where `T` is the number
of transforms in a map. So the total time to apply all the maps to the seed is `O(M*T)`
where `M` is the number of maps. Finding the locations for all the seeds is `O(S*M*T)` where
`S` is the number of seeds. This could be improved to `O(S*M*log(T))` if we sorted the
transforms and used a binary search.

## Part 2

### The twist

It turns out the list of seed values is actually a list of seed value ranges, with
pairs of seed values and counts, and we need to find the lowest location value for any of
these seed values.

### The solution

This can be brute forced by simply applying part 1's solution to each seed value, but for
the input `S` is very large. This could probably be calculated in the scale of minutes,
but we can do better.

Instead of mapping every seed, we can break down the ranges of seed values into smaller
ranges that fit within the map's transform ranges. Now we can map just the start and end
of each range to get a new set of ranges for the next map. Once we have the final list of
ranges, we can just use the lowest start value of the final range.

This has a running time of roughly `O(R*M*T)` where `R` is the number of seed ranges.
There's a bit of added complexity in that the transforms can split the ranges, but it's
not relevant in practice for the inputs for this problem. The number of ranges stays
manageable.

