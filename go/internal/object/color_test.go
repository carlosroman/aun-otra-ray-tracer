package object_test

import (
	"testing"

	"github.com/carlosroman/aun-otra-ray-tracer/go/internal/object"
)

func TestRGB_Multiply(t *testing.T) {
	c1 := object.NewColor(1, 0.2, 0.4)
	c2 := object.NewColor(0.9, 1, 0.1)
	expected := object.NewColor(0.9, 0.2, 0.04)
	actual := c1.Multiply(c2)
	assertColorEqual(t, expected, actual)
	assertColorEqual(t, object.NewColor(1, 0.2, 0.4), c1)
	assertColorEqual(t, object.NewColor(0.9, 1, 0.1), c2)
}

func TestRGB_Add(t *testing.T) {
	c1 := object.NewColor(0.9, 0.6, 0.75)
	c2 := object.NewColor(0.7, 0.1, 0.25)
	expected := object.NewColor(1.6, 0.7, 1.0)
	actual := c1.Add(c2)
	assertColorEqual(t, expected, actual)
	assertColorEqual(t, object.NewColor(0.9, 0.6, 0.75), c1)
	assertColorEqual(t, object.NewColor(0.7, 0.1, 0.25), c2)
}

func TestRGB_Subtract(t *testing.T) {
	c1 := object.NewColor(0.9, 0.6, 0.75)
	c2 := object.NewColor(0.7, 0.1, 0.25)
	expected := object.NewColor(0.2, 0.5, 0.5)
	actual := c1.Subtract(c2)
	assertColorEqual(t, expected, actual)
	assertColorEqual(t, object.NewColor(0.9, 0.6, 0.75), c1)
	assertColorEqual(t, object.NewColor(0.7, 0.1, 0.25), c2)
}

func TestRGB_MultiplyBy(t *testing.T) {
	c := object.NewColor(0.2, 0.3, 0.4)
	expected := object.NewColor(0.4, 0.6, 0.8)
	actual := c.MultiplyBy(2)
	assertColorEqual(t, expected, actual)
	assertColorEqual(t, object.NewColor(0.2, 0.3, 0.4), c)
}

func TestBlack(t *testing.T) {
	assertColorEqual(t, object.NewColor(0, 0, 0), object.Black)
}

func TestWhite(t *testing.T) {
	assertColorEqual(t, object.NewColor(1, 1, 1), object.White)
}
