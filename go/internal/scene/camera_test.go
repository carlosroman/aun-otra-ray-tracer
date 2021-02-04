package scene_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/object"
	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
	"github.com/carlosroman/aun-otra-ray-trace/go/internal/scene"
)

func TestNewBasicCamera(t *testing.T) {
	camera, err := scene.NewBasicCamera(160, 120, math.Pi/2)
	require.NoError(t, err)
	assert.Equal(t, 160, camera.HSize())
	assert.Equal(t, 120, camera.VSize())
	assert.Equal(t, math.Pi/2, camera.FieldOfView())
}

func TestNewCamera(t *testing.T) {
	c, err := scene.NewCamera(200, 100, defaultFrom, defaultTo, defaultUp)
	require.NoError(t, err)
	assert.Equal(t, 200, c.HSize())
	assert.Equal(t, 100, c.VSize())
	assert.Equal(t, ray.NewPoint(0, 0, 0), c.Origin())
}

func TestCamera_RayForPixel(t *testing.T) {
	testCases := []struct {
		name      string
		nx        float64
		ny        float64
		origin    ray.Vector
		direction ray.Vector
		transform ray.Matrix
	}{
		{
			name:      "Constructing a ray through the center of the canvas",
			nx:        100,
			ny:        50,
			origin:    ray.NewPoint(0, 0, 0),
			direction: ray.NewVec(0, 0, -1),
		},
		{
			name:      "Constructing a ray through a corner of the canvas",
			nx:        0,
			ny:        0,
			origin:    ray.NewPoint(0, 0, 0),
			direction: ray.NewVec(0.66519, 0.33259, -0.66851),
		},
		{
			name:      "Constructing a ray when the camera is transformed",
			nx:        100,
			ny:        50,
			transform: ray.Rotation(ray.Y, math.Pi/4).Multiply(ray.Translation(0, -2, 5)),
			origin:    ray.NewPoint(0, 2, -5),
			direction: ray.NewVec(math.Sqrt(2)/2, 0, -math.Sqrt(2)/2),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			c, err := scene.NewBasicCamera(201, 101, math.Pi/2)
			require.NoError(t, err)
			if tt.transform != nil {
				require.NoError(t, c.SetTransform(tt.transform))
			}
			pixel := c.RayForPixel(tt.nx, tt.ny)
			require.NotNil(t, pixel)
			assert.Equal(t, tt.origin, pixel.Origin())
			assertVectorEqual(t, tt.direction, pixel.Direction())
		})
	}
}

func assertVectorEqual(t *testing.T, expected ray.Vector, actual ray.Vector) {
	assert.InDelta(t, expected.GetX(), actual.GetX(), 0.0001)
	assert.InDelta(t, expected.GetY(), actual.GetY(), 0.0001)
	assert.InDelta(t, expected.GetZ(), actual.GetZ(), 0.0001)
	assert.InDelta(t, expected.GetW(), actual.GetW(), 0.0001)
}

func TestCamera_PixelSize(t *testing.T) {
	testCases := []struct {
		name     string
		hSize    int
		vSize    int
		expected float64
	}{
		{
			name:     "The pixel size for a horizontal canvas",
			hSize:    200,
			vSize:    125,
			expected: 0.01,
		},
		{
			name:     "The pixel size for a vertical canvas",
			hSize:    125,
			vSize:    200,
			expected: 0.01,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			c, err := scene.NewBasicCamera(tt.hSize, tt.vSize, math.Pi/2)
			require.NoError(t, err)
			actual := c.PixelSize()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

var (
	defaultFrom = ray.NewPoint(0, 0, 0)
	defaultTo   = ray.NewPoint(0, 0, -1)
	defaultUp   = ray.NewVec(0, 1, 0)
)

func TestRender(t *testing.T) {
	w, c := setupRenderTest(t)
	img := scene.Render(c, w)
	assertRenderImg(t, img)
}

func TestMultiThreadedRender(t *testing.T) {
	w, c := setupRenderTest(t)
	img := scene.MultiThreadedRender(c, w, 10, 100)
	assertRenderImg(t, img)
}

func assertRenderImg(t *testing.T, img scene.Canvas) {
	assert.NotNil(t, img)
	assertColorEqual(t, object.NewColor(0.38066, 0.47583, 0.2855), img[5][5])
}

func setupRenderTest(t *testing.T) (scene.World, scene.Camera) {
	w, err := scene.DefaultWorld()
	require.NoError(t, err)
	c, err := scene.NewBasicCamera(11, 11, math.Pi/2)
	require.NoError(t, err)
	from := ray.NewPoint(0, 0, -5)
	to := ray.NewPoint(0, 0, 0)
	up := ray.NewVec(0, 1, 0)
	require.NoError(t, c.SetTransform(ray.ViewTransform(from, to, up)))
	return w, c
}

var benchImg scene.Canvas

func BenchmarkRender(b *testing.B) {
	var canvas scene.Canvas
	w, err := scene.DefaultWorld()
	require.NoError(b, err)
	c, err := scene.NewBasicCamera(11, 11, math.Pi/2)
	require.NoError(b, err)
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		canvas = scene.Render(c, w)
	}
	benchImg = canvas
}

func BenchmarkMultiThreadedRender(b *testing.B) {
	var canvas scene.Canvas
	w, err := scene.DefaultWorld()
	require.NoError(b, err)
	c, err := scene.NewBasicCamera(11, 11, math.Pi/2)
	require.NoError(b, err)
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		canvas = scene.MultiThreadedRender(c, w, 10, 100)
	}
	benchImg = canvas
}
