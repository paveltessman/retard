# Technical interview prep plan (4 days)

Interview: about 60 minutes, 2-3 coding tasks, plus theory questions during the solution. Language: Go.

## The main rule for these 4 days

1. No LLM while you solve. Open an LLM only after the task is solved, or after 40 minutes with no progress.
2. A timer on each task: 25 minutes. If you go over, read the solution, then write the code again from scratch from memory.
3. Write in a plain editor, without autocomplete and without a compiler. Compile only at the end.
4. Speak out loud in Russian, like in the interview. This is a separate skill, and it drops first.
5. Write every algorithm from the "write from memory" list at least twice, on different days.

## Task protocol for the interview

1. Restate the task in your own words. Ask about constraints and edge cases: empty input, duplicates, overflow, value range.
2. Work through an example by hand on paper.
3. Name the naive solution and its complexity. Say out loud that it is naive.
4. Name the improvement and its complexity. Agree with the interviewer before you write code.
5. Write the code. Say what you do while you write.
6. Run the code on the example by hand.
7. Name the time complexity and the memory complexity. Name the edge cases you handled.

---

## Linear structures

Arrays, hash tables, two pointers, sliding window, prefix sums, stack, queue, deque, monotonic stack, linked list.

Write from memory:

- [x] stack and queue, on a slice and on a linked list
- [x] deque on a ring buffer
- [x] reverse a linked list iteratively
- [x] cycle detection in a list (Floyd algorithm, "tortoise and hare")

Tasks:

- [x] Two Sum, Valid Anagram, Group Anagrams
- [x] Contains Duplicate, Top K with a map plus sorting
- [x] Longest Substring Without Repeating Characters
- [x] Minimum Size Subarray Sum
- [x] Move Zeroes, Remove Duplicates from Sorted Array
- [x] Container With Most Water
- [x] Subarray Sum Equals K (prefix sums plus a map)
- [x] Product of Array Except Self
- [x] Valid Parentheses, Min Stack
- [ ] Daily Temperatures, Largest Rectangle in Histogram (monotonic stack)
- [ ] Reverse Linked List, Merge Two Sorted Lists, Linked List Cycle, Remove Nth Node From End, Middle of the Linked List

Go theory in the evening (60-90 minutes):

- slice: len, cap, growth on append, why append can damage the original array, copy, passing a slice into a function
- map: how a hash table works, collisions, amortized O(1) and the worst case, random iteration order, a map is not goroutine safe, you cannot take the address of an element
- string, []byte, []rune, UTF-8, why len(s) is not the number of characters

## Sorting, binary search, heap

Write from memory (this is the main goal of the day):

- [x] merge sort (and the merge of two sorted arrays as a separate function)
- [x] quicksort with Hoare partition, plus quickselect for the k-th order statistic
- [x] binary heap on a slice: siftUp, siftDown, push, pop, build in O(n)
- [x] heapsort
- [x] insertion sort (useful as the answer to "what is faster on small arrays")

Binary search, three templates, write each one from memory:

- [x] classic value search
- [x] lower bound (first element >= x)
- [x] upper bound (first element > x)

Think about the boundary invariants separately. This is where people make the most mistakes.

Tasks:

- [x] Search in Rotated Sorted Array
- [ ] binary search on the answer: Koko Eating Bananas, Capacity To Ship Packages, Split Array Largest Sum
- [ ] Kth Largest Element in an Array (heap and quickselect, both ways)
- [ ] Top K Frequent Elements
- [ ] Merge k Sorted Lists
- [ ] Merge Intervals, Insert Interval, Non-overlapping Intervals
- [ ] task 3 from the contest, solve it yourself and measure the time honestly

Theory in the evening:

- complexity and stability of bubble, insertion, merge, quick, heap, counting, radix
- why quicksort is O(n^2) in the worst case and how people avoid it
- what sort.Slice uses in Go (pdqsort) and why it is not stable, how sort.SliceStable differs
- why stability matters in practice (sorting by several keys)

Practical Go advice: do not write container/heap in the interview. The five method interface takes time, and it is easy to forget Push/Pop with a pointer. Write your own heap on a slice. It is 30 lines and full control.

## Recursion and trees

Write from memory:

- [x] binary tree traversals: preorder, inorder, postorder, recursive and iterative with a stack
- [x] level order traversal (BFS with a queue, split into levels)
- [ ] BST: search, insert, delete a node (three cases), BST validation with min/max bounds
- [ ] tree height, diameter

Tasks:

- Maximum Depth of Binary Tree, Invert Binary Tree, Symmetric Tree
- Binary Tree Level Order Traversal, Right Side View
- Validate Binary Search Tree, Kth Smallest Element in a BST
- Lowest Common Ancestor: in a BST separately, in a general tree separately
- Diameter of Binary Tree, Path Sum, Binary Tree Maximum Path Sum
- Serialize and Deserialize Binary Tree (if time is left)
- backtracking: Subsets, Permutations, Combination Sum

Theory in the evening:

- BST: why it is needed, operation complexity, degeneration into a list, what a balanced tree is, the idea of an AVL tree and a red-black tree at the level of "why and what complexity"
- B-tree and why an index in PostgreSQL is a B-tree, not a BST. Backend interviewers really like this question
- heap against BST: what is faster where, and why

## Graphs

This is the largest gap, so it takes the whole day.

Representations: adjacency list (a map or a slice of slices), adjacency matrix, edge list. Be able to name the memory and the speed of each one.

Write from memory:

- BFS and DFS on an adjacency list (DFS recursive and on a stack)
- BFS on a grid (this is exactly task 1 of the contest)
- connected components search
- shortest path in an unweighted graph with BFS, with path reconstruction
- topological sort: Kahn algorithm on a queue, and the DFS version
- cycle detection: in a directed graph separately (three colors), in an undirected graph separately
- disjoint set union (DSU) with path compression and union by rank
- Dijkstra algorithm on your own heap

Tasks:

- Number of Islands, Max Area of Island, Rotting Oranges (this is task 1 of the contest), 01 Matrix, Surrounded Regions
- Clone Graph
- Course Schedule I and II (topological sort)
- Is Graph Bipartite
- Word Ladder (BFS on an implicit graph)
- Number of Provinces, Redundant Connection (DSU)
- Network Delay Time, Path With Minimum Effort (Dijkstra)

Theory:

- BFS and DFS complexity: O(V + E), and why
- where BFS gives the shortest path, and where it does not
- Dijkstra against BFS against the Bellman-Ford algorithm, what to do with negative weights
- minimum spanning tree: the idea of the Prim and Kruskal algorithms, without an implementation

## Minimal DP block (only if time is left)

In a one hour interview with 2-3 tasks, heavy DP is rare. But one dimensional DP is possible. If you are ahead of the schedule, spend 2 hours:

- Climbing Stairs, House Robber, Coin Change, Longest Increasing Subsequence, 0/1 knapsack
- be able to explain the move from recursion with memoization to a table

## Go theory for backend (30-40 minutes every day, not in one block)

Interviewers ask these questions while you solve the task, so the answer must be short and confident.

Concurrency:

- a goroutine and an OS thread, the G-M-P scheduler in general terms
- channel: buffered and unbuffered, closing a channel, reading from a closed channel, select, default
- sync.Mutex, sync.RWMutex, sync.WaitGroup, sync.Once
- context: cancellation, deadline, passing through layers
- data race, the race detector, a classic deadlock
- patterns: worker pool, fan-in and fan-out, parallelism limit with a semaphore on a channel

Language:

- interface, a nil interface against an interface with a nil pointer inside, type assertion
- value against pointer receivers
- defer: order, the moment arguments are evaluated, defer in a loop
- errors: wrapping with %w, errors.Is, errors.As, panic and recover
- garbage collector, escape analysis, stack against heap

Backend minimum (they can ask one or two questions):

- a database index, when it is not used, EXPLAIN
- transactions and isolation levels
- HTTP: methods, codes, idempotency, keep-alive
- TCP against UDP, what a timeout and a retry are

## The evening before the interview

1. Write BFS on a grid, binary search (lower bound) and a heap from scratch. This is a warmup, not a test.
2. Say the task protocol out loud.
3. Check the equipment: camera, microphone, editor, no second monitor needed, Telegram off.

## About the language

Snippets to remember:

- reading input with bufio.Scanner and a larger buffer (you already have it in task 1)
- sort.Slice and sort.SliceStable with a comparator
- your own heap on a slice
- a queue on a slice with a head index, without shifts
