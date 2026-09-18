# The test recipe

The user writes a solution alone, with no compiler and no LLM, then runs the tests. The tests are the only judge. So they must catch a wrong answer on any input the constraints allow, and they must say clearly what went wrong.

Every test compares the solution with a plain reference, or with an answer that the construction of the input already gives. A test that repeats the logic of the solution proves nothing.

## The two helpers

Write these first. Every test uses them.

### brute

The slow and obvious answer. Two nested loops are fine. It must be so plain that a reader sees it is right.

```go
// brute gives the answer of the task in the slow and plain way: it looks at
// every number and gives back the index of the one that is target. The tests
// compare the solution with it.
func brute(nums []int, target int) int {
	for i, v := range nums {
		if v == target {
			return i
		}
	}

	return -1
}
```

### check

It calls the solution on a copy, compares the answer, and reads the copy back. The read back catches a solution that sorts or moves the input. It gives back `false` on an error, so a loop over thousands of inputs stops at the first one instead of printing thousands of lines.

```go
func check(t *testing.T, nums []int, target, want int) bool {
	t.Helper()

	got := slices.Clone(nums)
	ok := true

	if i := search(got, target); i != want {
		t.Errorf("search(%v, %d) = %d, want %d", nums, target, i, want)
		ok = false
	}
	if !slices.Equal(got, nums) {
		t.Errorf("search(%v, %d) changed the array to %v", nums, target, got)
		ok = false
	}

	return ok
}
```

Drop the read back only when the task asks for a change in place. Then check the whole state of the input instead.

## The tests

Write the ones that fit the task. Most tasks take all of them.

### 1. The table

Named cases in a slice, each one in `t.Run`. Start with every example of `t.md`, under the names `Example1`, `Example2`. Then the edges:

- The smallest input the constraints allow, and the one above it.
- Both ends of the answer: the first index, the last index.
- The answer is absent.
- The extreme values of the constraints, such as `10000` and `-10000`.
- The shapes the task names: sorted, rotated, empty, one element, all equal, all different.

Twenty to forty cases is a normal size. Give each one a name that says the shape, such as `MissBetweenTheParts`, not `Case7`.

### 2. Every short input

Walk every input of up to 6 to 9 elements, taken from a small value set, and compare with `brute`. This is the test that finds the off-by-one. The input space is small, so it covers what a human never thinks of.

Pick the value set so that misses exist: values `0, 2, 4` leave `1` and `3` as misses. Add a second set that holds negative numbers.

Stop at the first error with the `bool` of `check`.

### 3. Random

Random sizes and random values against `brute`, with a fixed seed such as `rand.NewSource(1)`. A fixed seed gives the same run every time, so a failure comes back.

Run several shapes, because the shape decides what breaks: a narrow value range gives many equal values, a wide range gives few.

### 4. Neighbours

Give a part of a larger array, `arr[2:6]`, where the numbers outside break the order or repeat the answer. A solution that reads outside of the input fails here.

### 5. The largest input

Build the largest input the constraints allow. Get the answer from the construction, not from `brute`, which is too slow at that size. Example: an array of 5000 numbers with a gap of 4 has a known index for every value, and every value plus one is a miss.

This test also shows a solution that is correct but too slow to finish.

### 6. The complexity limit

Only when the problem states one, such as `O(log n)` or `O(n)`.

The constraints are usually too small to separate a fast solution from a slow one, so build an input far above them. Run many queries under a time limit and fail on the limit with a message that names the reason.

```go
if time.Since(start) > limit {
	t.Fatalf("search needed more than %v for the searches in an array of %d numbers, so it is not logarithmic", limit, n)
}
```

Start it with `if testing.Short() { t.Skip("the test needs a large array") }`, so `go test -short` stays quick.

### 7. Benchmarks

One per interesting shape: a hit at the start, a hit at the end, a miss, the largest input. Use `for b.Loop()`.

## Style

- A doc comment over every test and every helper. Say what the test catches. "It catches a solution that drops the last number of the range" is useful. "Tests the function" is not.
- An error message prints the input, the answer, and the wanted answer, in that order: `search(%v, %d) = %d, want %d`.
- Fixed seeds, one per test, so two tests do not share a stream.
- `slices.Clone`, `slices.Equal` and `slices.Sort` in place of hand written loops.
- Table names in CamelCase, because they become the name of the subtest.
- The comments follow the repo style: simple words, active voice, short sentences.
