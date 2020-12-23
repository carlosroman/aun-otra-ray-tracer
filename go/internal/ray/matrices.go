package ray

import (
	"fmt"
)

type Matrix [][]float64

type RowValues []float64

func NewMatrix(rows, columns int, rowsValues ...RowValues) Matrix {
	m := make(Matrix, rows)
	for r := 0; r < rows; r++ {
		m[r] = make([]float64, columns)
	}
	for r := range rowsValues {
		for i := range rowsValues[r] {
			m[r][i] = rowsValues[r][i]
		}
	}
	return m
}

func (m Matrix) Get(row, column int) float64 {
	return m[row][column]
}

func (m Matrix) Set(row, column int, val float64) {
	m[row][column] = val
}

func (m Matrix) SetRow(row int, vals ...float64) {
	for i := range m[row] {
		m[row][i] = vals[i]
	}
}

func (m Matrix) Transpose() (result Matrix) {
	result = NewMatrix(len(m[0]), len(m))
	for row := range m {
		for col := range m[row] {
			result[col][row] = m[row][col]
		}
	}
	return result
}

func (m Matrix) Multiply(by Matrix) (result Matrix) {
	result = NewMatrix(len(m), len(by[0]))

	for row := range result {
		for col := range result[row] {
			var val float64
			for byCol := range by {
				val = val + (m[row][byCol] * by[byCol][col])
			}
			result[row][col] = val
		}
	}
	return result
}

func (m Matrix) Determinant() float64 {
	if len(m) == 2 && len(m[0]) == 2 {
		return (m[0][0] * m[1][1]) - (m[0][1] * m[1][0])
	}

	var res float64
	for col := range m[0] {
		c := m.Cofactor(0, col)
		res = res + (m[0][col] * c)
	}

	return res
}

func (m Matrix) SubMatrix(row, col int) (result Matrix) {
	result = NewMatrix(len(m)-1, len(m[0])-1)
	for r := range m {
		if r == row {
			continue
		}
		currentRow := r
		if currentRow > row {
			currentRow--
		}
		for c := range m[r] {
			if c == col {
				continue
			}
			currentCol := c
			if currentCol > col {
				currentCol--
			}
			result[currentRow][currentCol] = m[r][c]
		}
	}
	return result
}

func (m Matrix) Minor(row, col int) float64 {
	return m.
		SubMatrix(row, col).
		Determinant()
}

func (m Matrix) Cofactor(row, col int) float64 {
	if (row+col)%2 == 0 {
		return m.Minor(row, col)
	}
	return -m.Minor(row, col)
}

func (m Matrix) MultiplyByTuple(vals ...float64) (result Matrix) {
	tuple := NewMatrix(len(vals), 1)
	for i := range vals {
		tuple[i] = []float64{vals[i]}
	}

	return m.Multiply(tuple)
}

func (m Matrix) String() string {
	s := "[\n"
	for row := range m {
		for col := range m[row] {
			s = s + fmt.Sprintf("%v ", m[row][col])
		}
		s = s + "\n"
	}
	s = s + "]"
	return s
}

func IdentityMatrix(rows, cols int) (result Matrix) {
	result = NewMatrix(rows, cols)

	for row := range result {
		for col := range result[row] {
			if col == row {
				result[row][col] = 1
			} else {
				result[row][col] = 0
			}
		}
	}
	return result
}
