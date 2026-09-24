package bfs

const (
	dino   = 'D'
	intern = 'S'
	wall   = '*'
)

type office struct {
	rowCount int
	colCount int
	cells    []rune
}

func newOffice(grid []string) *office {
	rows := len(grid)
	cols := len(grid[0])
	cells := make([]rune, 0, rows*cols)

	for _, line := range grid {
		cells = append(cells, []rune(line)...)
	}

	return &office{rowCount: rows, colCount: cols, cells: cells}
}

func (o *office) at(i int) rune {
	return o.cells[i]
}

func (o *office) neighbours(i int) []int {
	buf := make([]int, 0, 4)

	// left
	if i%o.colCount > 0 && o.at(i-1) != wall {
		buf = append(buf, i-1)
	}

	// right
	if i%o.colCount != o.colCount-1 && o.at(i+1) != wall {
		buf = append(buf, i+1)
	}

	//top
	if i-o.colCount >= 0 && o.at(i-o.colCount) != wall {
		buf = append(buf, i-o.colCount)
	}

	// bottom
	if i+o.colCount < len(o.cells) && o.at(i+o.colCount) != wall {
		buf = append(buf, i+o.colCount)
	}
	return buf
}

func minTime(grid []string) int {
	office := newOffice(grid)

	waiting := 0
	invited := make([]bool, len(office.cells))
	currentWave := make([]int, 0, len(office.cells))

	for i, role := range office.cells {
		if role == wall {
			continue
		}
		if role == dino {
			invited[i] = true
			currentWave = append(currentWave, i)
			continue
		}
		// role = intern
		waiting++
	}

	time := 0
	nextWave := make([]int, len(office.cells))
	for waiting > 0 && len(currentWave) > 0 {
		nextWave = nextWave[:0]

		for _, index := range currentWave {
			for _, position := range office.neighbours(index) {
				if invited[position] || office.at(position) != intern {
					continue
				}
				invited[position] = true
				waiting--
				nextWave = append(nextWave, position)
			}
		}
		time++
		if waiting == 0 {
			break
		}
		currentWave, nextWave = nextWave, currentWave
	}

	if waiting > 0 {
		return -1
	}
	return time
}
