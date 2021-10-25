package ray_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/carlosroman/aun-otra-ray-tracer/go/internal/ray"
)

func Test_Dot(t *testing.T) {
	a := ray.NewVec(1, 2, 3)
	b := ray.NewVec(2, 3, 4)
	actual := ray.Dot(a, b)
	assert.Equal(t, float64(20), actual)
	assert.Equal(t, float64(20), a.Dot(b))
}

func Test_Cross(t *testing.T) {
	a := ray.NewVec(1, 2, 3)
	b := ray.NewVec(2, 3, 4)

	assertVec(t, ray.NewVec(-1, 2, -1), ray.Cross(a, b))
	assertVec(t, ray.NewVec(1, -2, 1), ray.Cross(b, a))
}

func TestTranslation(t *testing.T) {
	expected := ray.NewMatrix(4, 4,
		ray.RowValues{1, 0, 0, 11},
		ray.RowValues{0, 1, 0, 12},
		ray.RowValues{0, 0, 1, 13},
		ray.RowValues{0, 0, 0, 1},
	)

	assert.Equal(t, expected, ray.Translation(11, 12, 13))
}

func TestScaling(t *testing.T) {
	expected := ray.NewMatrix(4, 4,
		ray.RowValues{11, 0, 0, 0},
		ray.RowValues{0, 12, 0, 0},
		ray.RowValues{0, 0, 13, 0},
		ray.RowValues{0, 0, 0, 1},
	)

	assert.Equal(t, expected, ray.Scaling(11, 12, 13))
}
