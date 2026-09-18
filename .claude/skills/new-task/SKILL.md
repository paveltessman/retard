---
name: new-task
description: Set up a Go algorithm practice task in this prep repo. It makes the folder, writes t.md with the problem text formatted, writes an empty implementation stub, and writes a test file that is proven against a reference solution and against common wrong solutions. Use it whenever the user pastes a LeetCode or interview problem and wants a folder, a t.md, tests, or a place to start solving. Examples that must trigger it - "/new-task 7491/82/t/two_sum <problem text>", "set up the next task", "add tests for this problem", "format t.md and add tests", "I am going to solve X, prepare it". Also use it when the user names a task folder that does not exist yet.
---

# new-task

Set up one practice task so the user can solve it and check the answer alone.

```
/new-task <path/to/problem> <problem text>
```

## What you produce

Three files in the task folder:

| File | Content |
| --- | --- |
| `t.md` | The problem text, formatted. The wording stays as it came. |
| `<name>.go` | The package clause and an empty stub with the signature. |
| `<name>_test.go` | The tests. |

The user solves the task with no help. Read `references/no-spoilers.md` before you write the stub comment and the test comments. It says what you can say and what gives the answer away.

## Read the arguments

The first argument is the folder. The rest is the problem text.

- A path with a slash, such as `7491/82/t/search_rotated`, is relative to the repo root.
- A bare name, such as `two_sum`, goes into the newest `*/t/` folder of the repo. Find it with `ls -dt 7491/*/t | head -1`. Name the folder you picked in your report.
- No problem text, and `t.md` already exists: format that file and keep going.
- No problem text and no `t.md`: ask for the text. You cannot write tests without it.

The problem text often comes from a paste of a web page. It arrives as one flat block.

## Names

The folder name drives every other name. For a folder `search_rotated`:

| Thing | Value | Rule |
| --- | --- | --- |
| Package | `searchrotated` | The folder name with the underscores removed. |
| Files | `search_rotated.go`, `search_rotated_test.go` | The folder name. |
| Function | `search` | The name from the problem, as the site gives it. |

The function name comes from the problem, not from the folder. `search_rotated` holds `search`, and `subarray_sum` holds `subarraySum`.

## Step 1. Write t.md

Format the text. Do not rewrite it. The user reads this file in place of the web page, so a changed wording can change the task.

Use this shape:

````markdown
# <number>. <Title>

<the statement, one line per paragraph, backticks around nums, target, k>

## Examples

### Example 1

```
Input:  nums = [4,5,6,7,0,1,2], target = 0
Output: 4
```

## Constraints

- `1 <= nums.length <= 5000`
- `-10^4 <= nums[i] <= 10^4`
````

Two fixes belong here, because a paste loses them:

- Exponents. `-104` means `-10^4`. `2 * 104` means `2 * 10^4`. `-231` means `-2^31`. Read the number as the problem, not as the paste.
- Code marks. Put backticks around the names of the arguments, the values, and the code fragments.

Keep every example and every constraint. The tests come from them.

## Step 2. Write the stub

```go
package searchrotated

// search gives the index of target in nums, or -1 when nums does not hold
// target. nums is an ascending array of distinct numbers, possibly left
// rotated at an unknown index. The search must run in O(log n) time.
func search(nums []int, target int) int {
	// TODO: implement
	return -1
}
```

The comment names the contract: the input, the answer, and the limit the problem sets. It does not name a technique.

If the problem needs a type, such as `ListNode` or `TreeNode`, define it in this file.

If the file already holds an implementation, leave it alone. Take the signature from it and write only the tests.

## Step 3. Write the tests

Read `references/test-recipe.md` first. It holds the list of tests, the helpers, and the shape of each one.

Read the nearest finished task as well, such as `7491/82/t/search_rotated/search_rotated_test.go`. It shows the house style: the fixed seeds, the doc comment over every test, and the error messages that print the input.

## Step 4. Prove the tests

A test file that nobody proved is worse than no test file. It gives the user a green run on a wrong solution.

Work in the scratchpad folder. The repo never sees a solution.

```bash
S=<scratchpad>
skill=.claude/skills/new-task

# A correct solution must pass.
$skill/scripts/verify_tests.sh 7491/82/t/search_rotated $S/good.go
# Every wrong solution must fail.
$skill/scripts/verify_tests.sh 7491/82/t/search_rotated $S/bad1.go
```

Write one correct solution and 3 to 5 wrong ones. Pick the wrong ones from the mistakes this task invites: an off-by-one at a bound, a `<` in place of a `<=`, the answer of a neighbour index, a solution that changes the input, and a solution that is correct but too slow. The slow one proves the complexity test.

The script builds a throwaway module from the test files plus the solution you give it, so it needs a solution file that also defines any helper type.

A test that no wrong solution fails is dead weight. Cut it or make it sharper.

Run `gofmt -l` and `go vet` on the task folder as well.

Delete the scratchpad solutions at the end. Never copy one into the repo.

## Step 5. Report

Name the three files. Say how to run the tests:

```
go test ./7491/82/t/search_rotated/
go test -short ./7491/82/t/search_rotated/   # skips the slow timing test
```

Say what the tests cover, in a short list. Say which wrong solutions you tried and how many tests caught each one. That number is what tells the user the tests are worth trusting.

Do not describe the algorithm. Do not say which approach the tests expect.
