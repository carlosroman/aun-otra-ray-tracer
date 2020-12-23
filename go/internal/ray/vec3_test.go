package ray_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

func TestVec3_Translation(t *testing.T) {
	transform := ray.Translation(5, -3, 2)
	p := ray.NewPoint(-3, 4, 5)

	t.Run("Multiplying by a translation matrix", func(t *testing.T) {
		expected := ray.NewPoint(2, 1, 7)
		assert.Equal(t, expected, transform.MultiplyByVector(p))
	})

	t.Run("Multiplying by the inverse of a translation matrix", func(t *testing.T) {
		inv, err := transform.Inverse()
		require.NoError(t, err)
		expected := ray.NewPoint(-8, 7, 3)
		assert.Equal(t, expected, inv.MultiplyByVector(p))

	})

	t.Run("Translation does not affect vectors", func(t *testing.T) {
		v := ray.NewVec(-3, 4, 5)
		assert.Equal(t, v, transform.MultiplyByVector(v))
	})
}

func assertVec(t *testing.T, expected, actual ray.Vector) {
	assert.InDelta(t, expected.GetX(), actual.GetX(), 0.00001)
	assert.InDelta(t, expected.GetY(), actual.GetY(), 0.00001)
	assert.InDelta(t, expected.GetZ(), actual.GetZ(), 0.00001)
}

func TestNewPoint(t *testing.T) {
	point := ray.NewPoint(11, 12, 13)
	assert.Equal(t, float64(11), point.GetX())
	assert.Equal(t, float64(12), point.GetY())
	assert.Equal(t, float64(13), point.GetZ())
	assert.Equal(t, float64(1), point.GetW())
}

func TestNewVec(t *testing.T) {
	vector := ray.NewVec(11, 12, 13)
	assert.Equal(t, float64(11), vector.GetX())
	assert.Equal(t, float64(12), vector.GetY())
	assert.Equal(t, float64(13), vector.GetZ())
	assert.Equal(t, float64(0), vector.GetW())
}
