package validsudoku

func isValidSudoku(board [][]byte) bool {

	var rows, cols, squares [9][9]bool

	for row, values := range board {
		for col, value := range values {
			if value == '.' {
				continue
			}

			v := value - '1'
			sq := (row/3)*3 + col/3
			if rows[row][v] || cols[col][v] || squares[sq][v] {
				return false
			}

			rows[row][v], cols[col][v], squares[sq][v] = true, true, true
		}
	}
	return true
}
