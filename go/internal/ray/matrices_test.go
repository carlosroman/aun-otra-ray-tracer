package ray_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

func TestMatrix_SetAndGet(t *testing.T) {
	// Given 4 x 4 Matrix
	m := ray.NewMatrix(4, 4)

	for i := 0; i < 4; i++ {
		m.Set(0, i, float64(i)+1)
	}

	for i := 0; i < 4; i++ {
		m.Set(1, i, float64(i)+5.5)
	}

	for i := 0; i < 4; i++ {
		m.Set(2, i, float64(i)+9)
	}

	for i := 0; i < 4; i++ {
		m.Set(3, i, float64(i)+13.5)
	}

	assert.Equal(t, float64(1), m.Get(0, 0))
	assert.Equal(t, float64(4), m.Get(0, 3))
	assert.Equal(t, 5.5, m.Get(1, 0))
	assert.Equal(t, 7.5, m.Get(1, 2))
	assert.Equal(t, float64(11), m.Get(2, 2))
	assert.Equal(t, 13.5, m.Get(3, 0))
	assert.Equal(t, 15.5, m.Get(3, 2))
}

func TestNewMatrix_2x2(t *testing.T) {
	m := ray.NewMatrix(2, 2)
	m.Set(0, 0, -3)
	m.Set(0, 1, 5)
	m.Set(1, 0, 1)
	m.Set(1, 1, -2)

	assert.Equal(t, float64(-3), m[0][0])
	assert.Equal(t, float64(5), m[0][1])
	assert.Equal(t, float64(1), m[1][0])
	assert.Equal(t, float64(-2), m[1][1])
}

func TestNewMatrix_3x3(t *testing.T) {
	m := ray.NewMatrix(3, 3)
	m[0][0] = -3
	m[0][1] = 5
	m[0][2] = 0
	m[1][0] = 1
	m[1][1] = -2
	m[1][2] = -7
	m[2][0] = 0
	m[2][1] = 1
	m[2][2] = 1

	assert.Equal(t, float64(-3), m.Get(0, 0))
	assert.Equal(t, float64(-2), m.Get(1, 1))
	assert.Equal(t, float64(1), m.Get(2, 2))
}

func TestNewMatrix_direct(t *testing.T) {
	// Given 4 x 4 Matrix
	m := ray.NewMatrix(4, 4)

	for i := 0; i < 4; i++ {
		m[0][i] = float64(i) + 1
	}

	for i := 0; i < 4; i++ {
		m[1][i] = float64(i) + 5.5
	}

	for i := 0; i < 4; i++ {
		m[2][i] = float64(i) + 9
	}

	for i := 0; i < 4; i++ {
		m[3][i] = float64(i) + 13.5
	}

	assert.Equal(t, float64(1), m[0][0])
	assert.Equal(t, float64(4), m[0][3])
	assert.Equal(t, 5.5, m[1][0])
	assert.Equal(t, 7.5, m[1][2])
	assert.Equal(t, float64(11), m[2][2])
	assert.Equal(t, 13.5, m[3][0])
	assert.Equal(t, 15.5, m[3][2])
}

func TestMatrix_Equal(t *testing.T) {
	a := Get4x4()
	b := Get4x4()
	assert.Equal(t, a, b)
	b.Set(1, 1, 10)
	assert.NotEqual(t, t, a, b)
}

func TestMatrix_Multiply_4x4_by_4x4(t *testing.T) {
	a := Get4x4()

	b := ray.NewMatrix(4, 4)
	b.SetRow(0, -2, 1, 2, 3)
	b.SetRow(1, 3, 2, 1, -1)
	b.SetRow(2, 4, 3, 6, 5)
	b.SetRow(3, 1, 2, 7, 8)

	expected := ray.NewMatrix(4, 4,
		ray.RowValues{20, 22, 50, 48},
		ray.RowValues{44, 54, 114, 108},
		ray.RowValues{40, 58, 110, 102},
		ray.RowValues{16, 26, 46, 42},
	)

	assert.Equal(t, expected, a.Multiply(b))
}

func TestMatrix_Multiply_3x2_by_2x2(t *testing.T) {
	a := ray.NewMatrix(3, 2)
	a.SetRow(0, 4, 5)
	a.SetRow(1, 3, 4)
	a.SetRow(2, 1, 2)

	b := ray.NewMatrix(2, 2)
	b.SetRow(0, 3, 1)
	b.SetRow(1, 0, 3)

	expected := ray.NewMatrix(3, 2)
	expected.SetRow(0, 12, 19)
	expected.SetRow(1, 9, 15)
	expected.SetRow(2, 3, 7)
	assert.Equal(t, expected, a.Multiply(b))
}

func TestIdentityMatrix(t *testing.T) {
	a := ray.IdentityMatrix(4, 4)
	expected := ray.NewMatrix(4, 4)
	expected.SetRow(0, 1, 0, 0, 0)
	expected.SetRow(1, 0, 1, 0, 0)
	expected.SetRow(2, 0, 0, 1, 0)
	expected.SetRow(3, 0, 0, 0, 1)
	assert.Equal(t, expected, a)

	// Multiplying a matrix by the identity matrix returns the same matrix
	assert.Equal(t, Get4x4(), Get4x4().Multiply(a))
}

func TestMatrix_Transpose(t *testing.T) {
	a := ray.NewMatrix(4, 4)
	a.SetRow(0, 0, 9, 3, 0)
	a.SetRow(1, 9, 8, 0, 8)
	a.SetRow(2, 1, 8, 5, 3)
	a.SetRow(3, 0, 0, 5, 8)

	expected := ray.NewMatrix(4, 4)
	expected.SetRow(0, 0, 9, 1, 0)
	expected.SetRow(1, 9, 8, 8, 0)
	expected.SetRow(2, 3, 0, 5, 5)
	expected.SetRow(3, 0, 8, 3, 8)

	assert.Equal(t, expected, a.Transpose())
	assert.Equal(t, ray.IdentityMatrix(4, 4), ray.IdentityMatrix(4, 4).Transpose())
}

func TestMatrix_Multiply_4x4_by_4x1(t *testing.T) {
	a := ray.NewMatrix(4, 4,
		ray.RowValues{1, 2, 3, 4},
		ray.RowValues{2, 4, 4, 2},
		ray.RowValues{8, 6, 4, 1},
		ray.RowValues{0, 0, 0, 1},
	)

	expected := ray.NewMatrix(4, 1,
		ray.RowValues{18},
		ray.RowValues{24},
		ray.RowValues{33},
		ray.RowValues{1},
	)
	assert.Equal(t, expected, a.MultiplyByTuple(1, 2, 3, 1))
}

func TestMatrix_Determinant(t *testing.T) {

	testCases := []struct {
		name     string
		matrix   ray.Matrix
		expected float64
	}{
		{
			name: "2x2",
			matrix: ray.NewMatrix(2, 2,
				ray.RowValues{1, 5},
				ray.RowValues{-3, 2},
			),
			expected: 17,
		},
		{
			name: "3x3",
			matrix: ray.NewMatrix(3, 3,
				ray.RowValues{1, 2, 6},
				ray.RowValues{-5, 8, -4},
				ray.RowValues{2, 6, 4},
			),
			expected: -196,
		},
		{
			name: "4x4",
			matrix: ray.NewMatrix(4, 4,
				ray.RowValues{-2, -8, 3, 5},
				ray.RowValues{-3, 1, 7, 3},
				ray.RowValues{1, 2, -9, 6},
				ray.RowValues{-6, 7, 7, -9},
			),
			expected: -4071,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expected, testCase.matrix.Determinant())
		})
	}
}

func TestMatrix_Minor(t *testing.T) {
	type args struct {
		row int
		col int
	}
	testCases := []struct {
		name     string
		args     args
		matrix   ray.Matrix
		expected float64
	}{
		{
			name: "Calculating a minor of a 3x3 matrix",
			args: args{
				row: 1,
				col: 0,
			},
			matrix: ray.NewMatrix(3, 3,
				ray.RowValues{3, 5, 0},
				ray.RowValues{2, -1, -7},
				ray.RowValues{6, -1, 5},
			),
			expected: 25,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			b := testCase.matrix.SubMatrix(testCase.args.row, testCase.args.col)
			assert.Equal(t,
				testCase.expected,
				testCase.matrix.Minor(testCase.args.row, testCase.args.col), "Did not match expected")
			assert.Equal(t,
				b.Determinant(),
				testCase.matrix.Minor(testCase.args.row, testCase.args.col), "Did not equal determinant of the sub-matrix")
		})
	}
}

func TestMatrix_Cofactor(t *testing.T) {
	a := ray.NewMatrix(3, 3,
		ray.RowValues{3, 5, 0},
		ray.RowValues{2, -1, -7},
		ray.RowValues{6, -1, 5},
	)

	assert.Equal(t, float64(-12), a.Minor(0, 0))
	assert.Equal(t, float64(-12), a.Cofactor(0, 0))

	assert.Equal(t, float64(25), a.Minor(1, 0))
	assert.Equal(t, float64(-25), a.Cofactor(1, 0))
}

func TestMatrix_SubMatrix(t *testing.T) {
	type args struct {
		row int
		col int
	}
	testCases := []struct {
		name     string
		args     args
		matrix   func() ray.Matrix
		expected ray.Matrix
	}{
		{
			name: "A sub-matrix of a 3x3 matrix is a 2x2 matrix",
			args: args{
				row: 0,
				col: 2,
			},
			matrix: func() ray.Matrix {
				return ray.NewMatrix(3, 3,
					ray.RowValues{1, 5, 0},
					ray.RowValues{-3, 2, 7},
					ray.RowValues{0, 6, -3},
				)
			},
			expected: ray.NewMatrix(2, 2,
				ray.RowValues{-3, 2},
				ray.RowValues{0, 6},
			),
		},
		{
			name: "A sub-matrix of a 3x3 removing first col",
			args: args{
				row: 1,
				col: 0,
			},
			matrix: func() ray.Matrix {
				return ray.NewMatrix(3, 3,
					ray.RowValues{3, 5, 0},
					ray.RowValues{2, -1, -7},
					ray.RowValues{6, -1, 5},
				)
			},
			expected: ray.NewMatrix(2, 2,
				ray.RowValues{5, 0},
				ray.RowValues{-1, 5},
			),
		},
		{
			name: "A sub-matrix of a 4x4 matrix is a 3x3 matrix",
			args: args{
				row: 2,
				col: 1,
			},
			matrix: func() ray.Matrix {
				return ray.NewMatrix(4, 4,
					ray.RowValues{-6, 1, 1, 6},
					ray.RowValues{-8, 5, 8, 6},
					ray.RowValues{-1, 0, 8, 2},
					ray.RowValues{-7, 1, -1, 1},
				)
			},
			expected: ray.NewMatrix(3, 3,
				ray.RowValues{-6, 1, 6},
				ray.RowValues{-8, 8, 6},
				ray.RowValues{-7, -1, 1},
			),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			matrix := testCase.matrix()
			assert.Equal(t, testCase.expected, matrix.SubMatrix(testCase.args.row, testCase.args.col), "SubMatrix matches")
			assert.Equal(t, testCase.matrix(), matrix, "Original matrix not altered")
		})
	}
}

func Get4x4() (m ray.Matrix) {
	m = ray.NewMatrix(4, 4)
	val := float64(1)
	gotNine := false
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			m.Set(row, col, val)
			if !gotNine {
				if val == 9 {
					gotNine = true
				} else {
					val++
				}
			}
			if gotNine {
				val--
			}
		}
	}
	return
}

func TestNewMatrix(t *testing.T) {
	a := ray.NewMatrix(2, 2)
	a.SetRow(0, 1, 2)
	a.SetRow(1, 3, 4)

	b := ray.NewMatrix(2, 2,
		ray.RowValues{1, 2},
		ray.RowValues{3, 4},
	)

	c := ray.NewMatrix(2, 2)
	c[0][0] = 1
	c[0][1] = 2
	c[1][0] = 3
	c[1][1] = 4

	assert.Equal(t, a, b)
	assert.Equal(t, b, c)
	assert.Equal(t, a, c)
}
