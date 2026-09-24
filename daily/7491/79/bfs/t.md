# The Dino dinosaurs invite the interns to a party

The Dino dinosaurs want to invite every intern in every open space to a big party, but they have no means of communication.

The office is a rectangular grid of size n×m. Every cell holds one of three characters:

`*` — empty space, which cannot pass an invitation;
`S` — an intern who does not know about the party yet;
`D` — a Dino dinosaur;

At time 0 every `D` cell already spreads the invitation. In each time unit every Dino and every already invited intern pass the invitation at the same time to all neighbour interns. Two cells are neighbours when they share a side: above, below, left or right.

Find the smallest number of time units that every intern needs to get an invitation.

If it is impossible to invite every intern, print -1.

How many time units do all the interns need to get an invitation to the party?

## Input format

The first line holds two natural numbers:

n — the number of rows of the matrix (`1 <= n <= 10^4`);

m — the number of columns of the matrix (`1 <= m <= 10^4`);

After that comes a two dimensional array n×m, with `n*m <= 10^4` (the open space), built from the characters:

empty cells (the character `*`);
interns (the character `S`);
Dino dinosaurs (the character `D`);

## Output format

Print a single whole number:

the smallest time in which every intern gets an invitation;
-1, if at least one intern cannot get an invitation;
if the grid holds no intern, print 0;

## Examples

### Example 1

```
Input:
2 2
D*
*S

Output:
-1
```

### Example 2

```
Input:
3 5
D****
SSSSS
****D

Output:
3
```

### Example 3

```
Input:
5 5
D****
SSDSD
*****
SSSSS
****D

Output:
5
```

## Note

cells with the character `*` are closed: an invitation cannot pass through them;
if the matrix holds several Dino dinosaurs, all of them start to invite at the same time;
if the matrix holds no intern, the answer is 0;
if at least one intern sits in a separate open space with no Dino, it is impossible to invite that intern, and the answer is -1;
An open space is a largest connected area of `S` and `D` cells, where a move is allowed only between cells that share a side. It is possible to invite every intern if and only if every open space that holds interns also holds at least one Dino.

## Constraints

- `1 <= n <= 10^4`
- `1 <= m <= 10^4`
- `n * m <= 10^4`
- every cell is one of `*`, `S`, `D`
