package ray

import (
	"fmt"
	"math"
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
	const colSizeForTuple = 1
	tuple := NewMatrix(len(vals), colSizeForTuple)
	for i := range vals {
		tuple[i] = []float64{vals[i]}
	}

	return m.Multiply(tuple)
}

func (m Matrix) MultiplyByVector(vec Vector) (result Vector) {
	mt := m.MultiplyByTuple(vec.GetX(), vec.GetY(), vec.GetZ(), vec.GetW())
	return newTuple(mt[0][0], mt[1][0], mt[2][0], mt[3][0])
}

func (m Matrix) Inverse() (result Matrix, err error) {
	d := m.Determinant()
	if d == 0 {
		return nil, NonInvertibleErr
	}
	result = NewMatrix(len(m), len(m[0]))
	for row := range m {
		for col := range m[row] {
			result[col][row] = m.Cofactor(row, col) / d
		}
	}
	return result, err
}

func Rotation(axis Axis, by float64) (rotation Matrix) {
	rotation = DefaultIdentityMatrix()
	cosBy := math.Cos(by)
	sinBy := math.Sin(by)
	switch axis {
	case X:
		rotation[1][1] = cosBy
		rotation[1][2] = -sinBy
		rotation[2][1] = sinBy
		rotation[2][2] = cosBy
	case Y:
		rotation[0][0] = cosBy
		rotation[0][2] = sinBy
		rotation[2][0] = -sinBy
		rotation[2][2] = cosBy
	case Z:
		rotation[0][0] = cosBy
		rotation[0][1] = -sinBy
		rotation[1][0] = sinBy
		rotation[1][1] = cosBy
	}
	return rotation
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

func DefaultIdentityMatrix() (identityMatrix Matrix) {
	return IdentityMatrix(defaultIdentityMatrixSize, defaultIdentityMatrixSize)
}

func DefaultIdentityMatrixInverse() (identityMatrixInverse Matrix) {
	identityMatrixInverse, _ = DefaultIdentityMatrix().Inverse()
	return identityMatrixInverse
}

func IdentityMatrix(rows, cols int) (identityMatrix Matrix) {
	identityMatrix = NewMatrix(rows, cols)

	for row := range identityMatrix {
		for col := range identityMatrix[row] {
			if col == row {
				identityMatrix[row][col] = 1
			} else {
				identityMatrix[row][col] = 0
			}
		}
	}
	return identityMatrix
}

func ViewTransform(from, to, up Vector) (result Matrix) {
	fwd := to.Subtract(from).Normalize()
	upN := up.Normalize()
	left := Cross(fwd, upN)
	trueUp := Cross(left, fwd)
	orientation := NewMatrix(4, 4,
		RowValues{left.GetX(), left.GetY(), left.GetZ(), 0},
		RowValues{trueUp.GetX(), trueUp.GetY(), trueUp.GetZ(), 0},
		RowValues{-fwd.GetX(), -fwd.GetY(), -fwd.GetZ(), 0},
		RowValues{0, 0, 0, 1},
	)
	return orientation.Multiply(Translation(-from.GetX(), -from.GetY(), -from.GetZ()))
}
