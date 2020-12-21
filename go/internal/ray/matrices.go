package ray

import "fmt"

type Matrix [][]float64

func NewMatrix(rows, columns int) Matrix {
	m := make(Matrix, rows)
	for r := 0; r < rows; r++ {
		m[r] = make([]float64, columns)
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
