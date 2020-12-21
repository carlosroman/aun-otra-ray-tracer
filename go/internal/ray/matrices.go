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

func (m Matrix) Multiply(by Matrix) (result Matrix) {
	result = NewMatrix(len(m), len(m[0]))

	for row := range m {
		for col := range m[row] {
			var val float64
			for byCol := range by {
				val = val + (m[row][byCol] * by[byCol][col])
			}
			result[row][col] = val
		}
	}
	return result
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
