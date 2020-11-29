package ray

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_vec3_Add(t *testing.T) {
	v1 := NewVec(3, -2, 5)
	v2 := NewVec(-2, 3, 1)

	actual := v1.Add(v2)
	assert.NotNil(t, actual)

	expected := NewVec(1, 1, 6)
	assertVec(t, expected, actual)
}

func Test_vec3_Subtract(t *testing.T) {
	v1 := NewVec(3, 2, 1)
	v2 := NewVec(5, 6, 7)

	actual := v1.Subtract(v2)
	assert.NotNil(t, actual)

	expected := NewVec(-2, -4, -6)
	assertVec(t, expected, actual)
}

func Test_vec3_Negate(t *testing.T) {
	v1 := NewVec(1, -2, 3)

	actual := v1.Negate()
	assert.NotNil(t, actual)

	expected := NewVec(-1, 2, -3)
	assertVec(t, expected, actual)
}

func Test_vec3_Magnitude(t *testing.T) {
	testCases := []struct {
		name              string
		v                 vec3
		expectedMagnitude float64
	}{
		{
			name:              "x",
			v:                 NewVec(1, 0, 0),
			expectedMagnitude: 1,
		},
		{
			name:              "y",
			v:                 NewVec(0, 1, 0),
			expectedMagnitude: 1,
		},
		{
			name:              "z",
			v:                 NewVec(0, 0, 1),
			expectedMagnitude: 1,
		},
		{
			name:              "positive",
			v:                 NewVec(1, 2, 3),
			expectedMagnitude: 3.74165738677,
		},
		{
			name:              "negative",
			v:                 NewVec(-1, -2, -3),
			expectedMagnitude: 3.74165738677,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := testCase.v.Magnitude()
			assert.NotNil(t, actual)
			assert.InDelta(t, testCase.expectedMagnitude, actual, 0.00001)
		})
	}
}
func Test_vec3_Multiply(t *testing.T) {
	testCases := []struct {
		name     string
		by       float64
		v        vec3
		expected vec3
	}{
		{
			name:     "scalar",
			v:        NewVec(1, -2, 3),
			by:       3.5,
			expected: NewVec(3.5, -7, 10.5),
		},
		{
			name:     "fraction",
			v:        NewVec(1, -2, 3),
			by:       0.5,
			expected: NewVec(0.5, -1, 1.5),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := testCase.v.Multiply(testCase.by)
			assert.NotNil(t, actual)
			assertVec(t, testCase.expected, actual)
		})
	}
}

func Test_vec3_Divide(t *testing.T) {
	v1 := NewVec(1, -2, 3)

	actual := v1.Divide(2)
	assert.NotNil(t, actual)

	expected := NewVec(0.5, -1, 1.5)
	assertVec(t, expected, actual)
}

func Test_vec3_Normalize(t *testing.T) {
	testCases := []struct {
		name     string
		v        vec3
		expected vec3
	}{
		{
			name:     "simple",
			v:        NewVec(4, 0, 0),
			expected: NewVec(1, 0, 0),
		},
		{
			name:     "complex",
			v:        NewVec(1, 2, 3),
			expected: NewVec(0.26726, 0.53452, 0.80178),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := testCase.v.Normalize()
			assertVec(t, testCase.expected, actual)
			assert.Equal(t, float64(1), actual.Magnitude())
		})
	}

}

func assertVec(t *testing.T, expected, actual Vector) {
	assert.InDelta(t, expected.GetX(), actual.GetX(), 0.00001)
	assert.InDelta(t, expected.GetY(), actual.GetY(), 0.00001)
	assert.InDelta(t, expected.GetZ(), actual.GetZ(), 0.00001)
}
