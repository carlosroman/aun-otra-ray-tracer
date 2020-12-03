package ray_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

func Test_vec3_Add(t *testing.T) {
	v1 := ray.NewVec(3, -2, 5)
	v2 := ray.NewVec(-2, 3, 1)

	actual := v1.Add(v2)
	assert.NotNil(t, actual)

	expected := ray.NewVec(1, 1, 6)
	assertVec(t, expected, actual)
}

func Test_vec3_Subtract(t *testing.T) {
	v1 := ray.NewVec(3, 2, 1)
	v2 := ray.NewVec(5, 6, 7)

	actual := v1.Subtract(v2)
	assert.NotNil(t, actual)

	expected := ray.NewVec(-2, -4, -6)
	assertVec(t, expected, actual)
}

func Test_vec3_Negate(t *testing.T) {
	v1 := ray.NewVec(1, -2, 3)

	actual := v1.Negate()
	assert.NotNil(t, actual)

	expected := ray.NewVec(-1, 2, -3)
	assertVec(t, expected, actual)
}

func Test_vec3_Magnitude(t *testing.T) {
	testCases := []struct {
		name              string
		v                 ray.Vector
		expectedMagnitude float64
	}{
		{
			name:              "x",
			v:                 ray.NewVec(1, 0, 0),
			expectedMagnitude: 1,
		},
		{
			name:              "y",
			v:                 ray.NewVec(0, 1, 0),
			expectedMagnitude: 1,
		},
		{
			name:              "z",
			v:                 ray.NewVec(0, 0, 1),
			expectedMagnitude: 1,
		},
		{
			name:              "positive",
			v:                 ray.NewVec(1, 2, 3),
			expectedMagnitude: math.Sqrt(14),
		},
		{
			name:              "negative",
			v:                 ray.NewVec(-1, -2, -3),
			expectedMagnitude: math.Sqrt(14),
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
		v        ray.Vector
		expected ray.Vector
	}{
		{
			name:     "scalar",
			v:        ray.NewVec(1, -2, 3),
			by:       3.5,
			expected: ray.NewVec(3.5, -7, 10.5),
		},
		{
			name:     "fraction",
			v:        ray.NewVec(1, -2, 3),
			by:       0.5,
			expected: ray.NewVec(0.5, -1, 1.5),
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
	v1 := ray.NewVec(1, -2, 3)

	actual := v1.Divide(2)
	assert.NotNil(t, actual)

	expected := ray.NewVec(0.5, -1, 1.5)
	assertVec(t, expected, actual)
}

func Test_vec3_Normalize(t *testing.T) {
	testCases := []struct {
		name     string
		v        ray.Vector
		expected ray.Vector
	}{
		{
			name:     "simple",
			v:        ray.NewVec(4, 0, 0),
			expected: ray.NewVec(1, 0, 0),
		},
		{
			name:     "complex",
			v:        ray.NewVec(1, 2, 3),
			expected: ray.NewVec(0.26726, 0.53452, 0.80178),
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

func assertVec(t *testing.T, expected, actual ray.Vector) {
	assert.InDelta(t, expected.GetX(), actual.GetX(), 0.00001)
	assert.InDelta(t, expected.GetY(), actual.GetY(), 0.00001)
	assert.InDelta(t, expected.GetZ(), actual.GetZ(), 0.00001)
}
