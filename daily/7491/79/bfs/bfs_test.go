package bfs

import (
	"fmt"
	"math/rand"
	"slices"
	"strings"
	"testing"
	"time"
)

// brute gives the answer of the task in the slow and plain way: it repeats the
// step that the task describes. In every round it marks the 'S' cells that
// share a side with an invited cell, and it counts the rounds. A round that
// marks nothing while an 'S' cell is still waiting gives -1. The tests compare
// the solution with it.
func brute(grid []string) int {
	n, m := len(grid), len(grid[0])

	invited := make([][]bool, n)
	left := 0

	for i := range n {
		invited[i] = make([]bool, m)

		for j := range m {
			switch grid[i][j] {
			case 'D':
				invited[i][j] = true
			case 'S':
				left++
			}
		}
	}

	steps := 0

	for left > 0 {
		var next [][2]int

		for i := range n {
			for j := range m {
				if grid[i][j] != 'S' || invited[i][j] {
					continue
				}

				if (i > 0 && invited[i-1][j]) || (i+1 < n && invited[i+1][j]) ||
					(j > 0 && invited[i][j-1]) || (j+1 < m && invited[i][j+1]) {
					next = append(next, [2]int{i, j})
				}
			}
		}

		if len(next) == 0 {
			return -1
		}

		for _, c := range next {
			invited[c[0]][c[1]] = true
		}

		left -= len(next)
		steps++
	}

	return steps
}

// show turns a grid into one line, so an error message fits on the screen. A
// grid that is too large for that gives its shape instead.
func show(grid []string) string {
	total := 0
	for _, row := range grid {
		total += len(row)
	}

	if total > 400 {
		return fmt.Sprintf("<a grid of %d rows and %d columns>", len(grid), len(grid[0]))
	}

	return strings.Join(grid, "/")
}

// check calls minTime on a copy of grid and compares the answer with want. It
// also reads the copy back, because minTime must leave the grid as it was. It
// gives back false on the first error, so a loop can stop.
func check(t *testing.T, grid []string, want int) bool {
	t.Helper()

	got := slices.Clone(grid)
	ok := true

	if v := minTime(got); v != want {
		t.Errorf("minTime(%s) = %d, want %d", show(grid), v, want)
		ok = false
	}
	if !slices.Equal(got, grid) {
		t.Errorf("minTime(%s) changed the grid to %s", show(grid), show(got))
		ok = false
	}

	return ok
}

// gridFromCode builds the grid of n rows and m columns that belongs to the
// number code, read as a number of base 3 over the characters '*', 'S' and
// 'D'. The codes from 0 to 3^(n*m)-1 give every grid of that shape.
func gridFromCode(code, n, m int) []string {
	const alphabet = "*SD"

	rows := make([]string, n)
	buf := make([]byte, m)

	for i := range n {
		for j := range m {
			buf[j] = alphabet[code%3]
			code /= 3
		}

		rows[i] = string(buf)
	}

	return rows
}

// randomGrid builds a grid of n rows and m columns. weights holds the share of
// '*', 'S' and 'D', so a caller can ask for many walls or for many dinosaurs.
func randomGrid(rnd *rand.Rand, n, m int, weights [3]int) []string {
	const alphabet = "*SD"

	total := weights[0] + weights[1] + weights[2]
	rows := make([]string, n)
	buf := make([]byte, m)

	for i := range n {
		for j := range m {
			r := rnd.Intn(total)

			k := 0
			for r >= weights[k] {
				r -= weights[k]
				k++
			}

			buf[j] = alphabet[k]
		}

		rows[i] = string(buf)
	}

	return rows
}

// fill gives n rows of m copies of c.
func fill(n, m int, c byte) []string {
	row := strings.Repeat(string(c), m)

	rows := make([]string, n)
	for i := range rows {
		rows[i] = row
	}

	return rows
}

// put writes c at row i and column j of grid.
func put(grid []string, i, j int, c byte) {
	buf := []byte(grid[i])
	buf[j] = c
	grid[i] = string(buf)
}

func TestMinTime(t *testing.T) {
	cases := []struct {
		name string
		grid []string
		want int
	}{
		{"Example1", []string{"D*", "*S"}, -1},
		{"Example2", []string{"D****", "SSSSS", "****D"}, 3},
		{"Example3", []string{"D****", "SSDSD", "*****", "SSSSS", "****D"}, 5},

		{"OneCellEmpty", []string{"*"}, 0},
		{"OneCellDino", []string{"D"}, 0},
		{"OneCellIntern", []string{"S"}, -1},

		{"NoInternsAtAll", []string{"D*D", "***", "*DD"}, 0},
		{"OnlyWalls", []string{"***", "***"}, 0},
		{"OnlyDinos", []string{"DD", "DD"}, 0},
		{"OnlyInterns", []string{"SSS", "SSS"}, -1},

		{"InternNextToDino", []string{"DS"}, 1},
		{"InternUnderDino", []string{"D", "S"}, 1},
		{"DinoAfterIntern", []string{"SD"}, 1},
		{"DinoUnderIntern", []string{"S", "D"}, 1},

		{"OneRowFromTheLeft", []string{"DSSSS"}, 4},
		{"OneRowFromTheRight", []string{"SSSSD"}, 4},
		{"OneRowFromBothEnds", []string{"DSSSSD"}, 2},
		{"OneRowDinoInTheMiddle", []string{"SSDSS"}, 2},
		{"OneColumn", []string{"D", "S", "S"}, 2},
		{"OneColumnFromBothEnds", []string{"D", "S", "S", "S", "D"}, 2},

		{"InternBehindAWall", []string{"DS*S"}, -1},
		{"WallBetweenTwoOpenSpacesWithDinos", []string{"DS*SD"}, 1},
		{"WallBetweenTheRows", []string{"SSS", "***", "SSS"}, -1},
		{"DinoOnlyInTheTopHalf", []string{"DSS", "***", "SSS"}, -1},
		{"DinoInEveryHalf", []string{"DSS", "***", "SSD"}, 2},

		{"WallForcesTheLongWay", []string{"D*S", "SSS"}, 4},
		{"LongCorridor", []string{"DS*", "*SS", "SS*"}, 4},
		{"AroundAWall", []string{"DS", "*S"}, 2},

		{"SquareFromTheCorner", []string{"DSS", "SSS", "SSS"}, 4},
		{"SquareFromTheCentre", []string{"SSS", "SDS", "SSS"}, 2},
		{"SquareFromTwoCorners", []string{"DSS", "SSS", "SSD"}, 2},
		{"SquareWithAWallInTheCentre", []string{"DSS", "S*S", "SSS"}, 4},

		{"TwoDinosNextToTheSameIntern", []string{"DSD"}, 1},
		{"DinosAtEveryCorner", []string{"DSD", "SSS", "DSD"}, 2},

		{"TallRectangle", []string{"DS", "SS", "SS", "SS", "SS"}, 5},
		{"WideRectangle", []string{"DSSSS", "SSSSS"}, 5},

		{"InternInAPocket", []string{"*S*", "*D*", "***"}, 1},
		{"PocketWithNoDino", []string{"*S*", "***", "*D*"}, -1},

		{"WallsSplitOneRow", []string{"S*D*S"}, -1},
		{"EveryOtherCellIsAWall", []string{"D*S*D"}, -1},
		{"DiagonalNeighbourIsNotANeighbour", []string{"D*", "*S"}, -1},
		{"DiagonalNeighbourWithAPath", []string{"DS", "SS"}, 2},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			check(t, c.grid, c.want)
		})
	}
}

// TestMinTimeEveryShortGrid walks every grid of up to 9 cells over the three
// characters and compares the solution with brute. The space is small enough
// to cover every shape a person never thinks of: a grid with no intern, a grid
// with no dinosaur, a grid that falls apart into several open spaces, and a
// grid of one row or one column. It catches an answer that is one too large or
// one too small, and an answer that misses a cell at the border.
func TestMinTimeEveryShortGrid(t *testing.T) {
	const maxCells = 9

	for n := 1; n <= maxCells; n++ {
		for m := 1; n*m <= maxCells; m++ {
			codes := 1
			for range n * m {
				codes *= 3
			}

			for code := range codes {
				grid := gridFromCode(code, n, m)

				if !check(t, grid, brute(grid)) {
					return
				}
			}
		}
	}
}

// TestMinTimeRandom compares the solution with brute on random grids of random
// shapes. The shares of the three characters change from round to round: many
// walls give many small open spaces, few walls give one large one, and few
// dinosaurs give long times.
func TestMinTimeRandom(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))

	shapes := [][3]int{
		{1, 8, 1},  // few walls, few dinosaurs
		{4, 5, 1},  // many walls
		{8, 3, 1},  // mostly walls
		{1, 20, 1}, // one large open space, one dinosaur in 22 cells
		{2, 5, 5},  // many dinosaurs
		{1, 1, 0},  // no dinosaur at all
		{1, 0, 1},  // no intern at all
	}

	for _, weights := range shapes {
		for range 200 {
			n := 1 + rnd.Intn(12)
			m := 1 + rnd.Intn(12)

			grid := randomGrid(rnd, n, m, weights)

			if !check(t, grid, brute(grid)) {
				return
			}
		}
	}
}

// TestMinTimeRandomTallAndWide compares the solution with brute on grids of one
// row and on grids of one column. It catches an answer that mixes up the row
// index with the column index.
func TestMinTimeRandomTallAndWide(t *testing.T) {
	rnd := rand.New(rand.NewSource(2))

	for range 500 {
		m := 1 + rnd.Intn(40)

		wide := randomGrid(rnd, 1, m, [3]int{2, 9, 1})
		if !check(t, wide, brute(wide)) {
			return
		}

		tall := randomGrid(rnd, m, 1, [3]int{2, 9, 1})
		if !check(t, tall, brute(tall)) {
			return
		}
	}
}

// TestMinTimeKeepsNeighbours gives a window of a taller grid to minTime. The
// rows above and below the window hold interns and a dinosaur, so a solution
// that reads outside of the rows it got gives a wrong answer here. The test
// also reads the whole grid back.
func TestMinTimeKeepsNeighbours(t *testing.T) {
	full := []string{"SSSS", "SSSS", "DS*S", "SSSS", "SSSS"}
	before := slices.Clone(full)

	grid := full[2:4]

	if got, want := minTime(grid), 5; got != want {
		t.Errorf("minTime(%s) = %d, want %d", show(grid), got, want)
	}
	if !slices.Equal(full, before) {
		t.Errorf("minTime(%s) changed the rows around it to %s", show(before[2:4]), show(full))
	}
}

// TestMinTimeLarge uses grids of 10^4 cells, the largest the constraints allow.
// The answer comes from the way the test builds the grid, because brute is too
// slow at this size. A solution that is correct but too slow runs out of time
// here.
func TestMinTimeLarge(t *testing.T) {
	t.Run("OneLongRow", func(t *testing.T) {
		grid := []string{"D" + strings.Repeat("S", 9999)}
		check(t, grid, 9999)
	})

	t.Run("OneLongColumn", func(t *testing.T) {
		grid := fill(10000, 1, 'S')
		put(grid, 9999, 0, 'D')
		check(t, grid, 9999)
	})

	t.Run("SquareFromTheCorner", func(t *testing.T) {
		grid := fill(100, 100, 'S')
		put(grid, 0, 0, 'D')
		check(t, grid, 198)
	})

	t.Run("SquareFromTheCentre", func(t *testing.T) {
		grid := fill(100, 100, 'S')
		put(grid, 50, 50, 'D')
		check(t, grid, 100)
	})

	t.Run("SquareWithAWallColumn", func(t *testing.T) {
		grid := fill(100, 100, 'S')
		put(grid, 0, 0, 'D')

		for i := range 100 {
			put(grid, i, 50, '*')
		}

		check(t, grid, -1)
	})

	t.Run("SquareWithADinoOnEachSideOfTheWall", func(t *testing.T) {
		grid := fill(100, 100, 'S')
		put(grid, 0, 0, 'D')
		put(grid, 0, 99, 'D')

		for i := range 100 {
			put(grid, i, 50, '*')
		}

		check(t, grid, 148)
	})

	// The walls leave one path of 5050 cells that folds row by row. The
	// dinosaur stands at one end of it, so the far end needs 5049 units.
	t.Run("FoldedPath", func(t *testing.T) {
		const n, m = 100, 100

		grid := make([]string, n)
		for i := range n {
			if i%2 == 0 {
				grid[i] = strings.Repeat("S", m)
				continue
			}

			buf := []byte(strings.Repeat("*", m))
			if i%4 == 1 {
				buf[m-1] = 'S'
			} else {
				buf[0] = 'S'
			}

			grid[i] = string(buf)
		}

		put(grid, 0, 0, 'D')

		check(t, grid, 5049)
	})

	t.Run("NoInternInTenThousandCells", func(t *testing.T) {
		grid := fill(100, 100, '*')
		put(grid, 42, 17, 'D')

		check(t, grid, 0)
	})
}

// TestMinTimeGrowsWithTheCells runs minTime on grids far above the constraints
// and stops it at a time limit. A solution whose work grows with the number of
// cells needs some milliseconds here. A solution that walks the whole grid once
// per time unit needs hours, so it runs into the limit. The constraints alone
// cannot show the difference, because 10^4 cells are too few.
func TestMinTimeGrowsWithTheCells(t *testing.T) {
	if testing.Short() {
		t.Skip("the test needs very large grids")
	}

	const limit = 10 * time.Second

	cases := []struct {
		name string
		grid []string
		want int
	}{
		{"OneRowOfTwoMillion", []string{"D" + strings.Repeat("S", 2000000-1)}, 1999999},
		{"TwoThousandRowsOfOneThousand", nil, 2998},
	}

	tall := fill(2000, 1000, 'S')
	put(tall, 0, 0, 'D')
	cases[1].grid = tall

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			done := make(chan int, 1)
			start := time.Now()

			go func() {
				done <- minTime(c.grid)
			}()

			select {
			case got := <-done:
				if got != c.want {
					t.Fatalf("minTime(%s) = %d, want %d", show(c.grid), got, c.want)
				}
				if spent := time.Since(start); spent > limit {
					t.Fatalf("minTime needed %v for a grid of %d rows and %d columns, so its work grows faster than the number of cells",
						spent, len(c.grid), len(c.grid[0]))
				}
			case <-time.After(limit):
				t.Fatalf("minTime needed more than %v for a grid of %d rows and %d columns, so its work grows faster than the number of cells",
					limit, len(c.grid), len(c.grid[0]))
			}
		})
	}
}

func BenchmarkMinTime(b *testing.B) {
	oneRow := []string{"D" + strings.Repeat("S", 9999)}

	corner := fill(100, 100, 'S')
	put(corner, 0, 0, 'D')

	many := fill(100, 100, 'S')
	for i := 0; i < 100; i += 10 {
		for j := 0; j < 100; j += 10 {
			put(many, i, j, 'D')
		}
	}

	blocked := fill(100, 100, 'S')
	put(blocked, 0, 0, 'D')
	for i := range 100 {
		put(blocked, i, 50, '*')
	}

	empty := fill(100, 100, '*')

	shapes := []struct {
		name string
		grid []string
	}{
		{"OneLongRow", oneRow},
		{"SquareFromTheCorner", corner},
		{"SquareWithManyDinos", many},
		{"SquareWithAWallColumn", blocked},
		{"NoIntern", empty},
	}

	for _, s := range shapes {
		b.Run(s.name, func(b *testing.B) {
			for b.Loop() {
				minTime(s.grid)
			}
		})
	}
}
