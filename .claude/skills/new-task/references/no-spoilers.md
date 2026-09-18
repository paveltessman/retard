# Do not give the answer away

`strat.md` sets the rule for these days: no LLM while the user solves. The user opens an LLM only after the task is solved, or after 40 minutes with no progress. The skill runs before the work starts, so everything it leaves behind is on the screen while the user thinks.

A file that names the technique takes the task away. The user then writes down an answer that they read, and the practice gives nothing.

So the three files carry the contract, and nothing about the method.

## What each file can say

| File | Can say | Must not say |
| --- | --- | --- |
| `t.md` | The problem text as it came, and the limit the problem states. | A hint, a note, a link, an approach. |
| `<name>.go` | The input, the answer, the complexity limit of the problem. | The technique, the data structure, the name of the algorithm. |
| `<name>_test.go` | The shape of the inputs and what a wrong answer looks like. | The steps of a correct solution. |

## The stub comment

Good, because every word comes from the problem:

```go
// search gives the index of target in nums, or -1 when nums does not hold
// target. nums is an ascending array of distinct numbers, possibly left
// rotated at an unknown index. The search must run in O(log n) time.
```

Bad, because it hands over the plan:

```go
// search finds the rotation point first, then runs a binary search over the
// half that holds target.
```

## The test comments

A test comment says what a wrong answer looks like, not how to get a right one.

Good: "It catches a solution that drops the number at the end of the range."

Bad: "It catches a solution that compares nums[mid] with nums[hi] in place of nums[lo]."

The second one names the comparison the user has to find. Keep the description on the level of the answer, not the code.

The same rule holds for the names in the table. `MissBetweenTheParts` describes the input. `ForgotThePivotCheck` describes the solution.

## The report in the chat

Say what the files hold and how to run the tests. Say which wrong solutions the tests caught, in the words of the answer: "a solution that returns the neighbour index", "a solution that is correct but linear".

Do not print the reference solution. Do not name the algorithm. If the user asks for the solution later, that is a new request, and it is theirs to make.
