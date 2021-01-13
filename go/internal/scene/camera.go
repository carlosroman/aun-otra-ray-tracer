package scene

import (
	"math"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

func NewCamera(hSize, vSize int, from, to, vup ray.Vector) Camera {
	return newCamera(hSize, vSize, math.Pi/2, ray.ViewTransform(from, to, vup))
}

func NewBasicCamera(hSize, vSize int, fieldOfView float64) Camera {
	return newCamera(hSize, vSize, fieldOfView, ray.IdentityMatrix(4, 4))
}

func newCamera(hSize, vSize int, fieldOfView float64, transform ray.Matrix) Camera {
	pixelSize, halfWidth, halfHeight := calculatePixelSize(hSize, vSize, fieldOfView)
	return &camera{
		hSize:       hSize,
		vSize:       vSize,
		fieldOfView: fieldOfView,
		origin:      ray.NewPoint(0, 0, 0),
		focalLength: 1.0,
		transform:   transform,
		pixelSize:   pixelSize,
		halfWidth:   halfWidth,
		halfHeight:  halfHeight,
	}
}

type Camera interface {
	HSize() int
	VSize() int
	Origin() ray.Vector
	PixelSize() float64
	FocalLength() float64
	RayForPixel(nx, ny float64) ray.Ray
	FieldOfView() float64
	SetTransform(by ray.Matrix)
}

type camera struct {
	hSize       int
	vSize       int
	fieldOfView float64
	focalLength float64
	origin      ray.Vector
	transform   ray.Matrix
	pixelSize   float64
	halfWidth   float64
	halfHeight  float64
}

func (c camera) HSize() int {
	return c.hSize
}

func (c camera) VSize() int {
	return c.vSize
}

func (c camera) FocalLength() float64 {
	return c.focalLength
}

func (c camera) FieldOfView() float64 {
	return c.fieldOfView
}

func calculatePixelSize(hSize, vSize int, fieldOfView float64) (pixelSize, halfWidth, halfHeight float64) {
	// half_view ← tan(camera.field_of_view / 2)
	halfView := math.Tan(fieldOfView / 2)
	aspect := float64(hSize) / float64(vSize)
	if aspect >= 1 {
		halfWidth = halfView
		halfHeight = halfView / aspect
	} else {
		halfWidth = halfView * aspect
		halfHeight = halfView
	}
	pixelSize = (halfWidth * 2) / float64(hSize)
	return
}

func (c camera) PixelSize() float64 {
	return c.pixelSize
}

func (c *camera) SetTransform(by ray.Matrix) {
	c.transform = by
}

func (c camera) RayForPixel(nx, ny float64) ray.Ray {
	// the offset from the edge of the canvas to the pixel's center

	xOffset := (nx + 0.5) * c.pixelSize
	yOffset := (ny + 0.5) * c.pixelSize

	// the untransformed coordinates of the pixel in world space.
	// (remember that the camera looks toward -z, so +x is to the *left*.)
	worldX := c.halfWidth - xOffset
	worldY := c.halfHeight - yOffset

	// using the camera matrix, transform the canvas point and the origin
	// and then compute the ray's direction vector.
	// (remember that the canvas is at z=-1)

	// pixel ← inverse(camera.transform) * point(world_x, world_y, -1)
	inv, _ := c.transform.Inverse()
	pixel := inv.MultiplyByVector(ray.NewPoint(worldX, worldY, -c.focalLength))
	origin := inv.MultiplyByVector(c.origin)
	// lower_left_corner = origin - horizontal/2 - vertical/2 - vec3(0, 0, focal_length)
	// lower_left_corner + u*horizontal + v*vertical - origin

	// direction ← normalize(pixel - origin)
	direction := pixel.Subtract(origin).Normalize()
	return ray.NewRayAt(origin, direction)
}

func (c camera) Origin() ray.Vector {
	return c.origin
}

func Render(c Camera, w World) Canvas {
	canvas := NewCanvas(c.HSize(), c.VSize())
	const cl = 255.99
	for y := 0; y < c.VSize()-1; y++ {
		for x := 0; x < c.HSize()-1; x++ {
			r := c.RayForPixel(float64(x), float64(y))
			canvas[x][y] = w.ColorAt(r)
		}
	}
	return canvas
}
