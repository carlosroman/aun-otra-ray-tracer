package output_test

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"io"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/carlosroman/aun-otra-ray-tracer/go/internal/object"
	"github.com/carlosroman/aun-otra-ray-tracer/go/internal/output"
	"github.com/carlosroman/aun-otra-ray-tracer/go/internal/scene"
)

func TestPPMWriter(t *testing.T) {
	dir, err := os.MkdirTemp("", "TestPPMWriter")
	require.NoError(t, err)
	defer func() {
		_ = os.RemoveAll(dir)
	}()

	testImages := []struct {
		name     string
		image    func() image.Image
		testFile string
	}{
		{
			name:     "black",
			testFile: path.Join("..", "..", "..", "test", "ppm", "black.ppm"),
			image: func() image.Image {
				return image.NewPaletted(image.Rect(0, 0, 10, 10), color.Palette{
					image.Black,
				})
			}},
		{
			name:     "white",
			testFile: path.Join("..", "..", "..", "test", "ppm", "white.ppm"),
			image: func() image.Image {
				return image.NewPaletted(image.Rect(0, 0, 10, 10), color.Palette{
					image.White,
				})
			}},
		{
			name:     "rgb",
			testFile: path.Join("..", "..", "..", "test", "ppm", "rgb.ppm"),
			image: func() image.Image {
				nx := 3
				ny := 2
				img := image.NewRGBA(image.Rect(0, 0, nx, ny))

				// red
				img.Set(0, 0, color.RGBA{
					R: 255, G: 0, B: 0,
					A: 0xff,
				})
				// green
				img.Set(1, 0, color.RGBA{
					R: 0, G: 255, B: 0,
					A: 0xff,
				})
				// blue
				img.Set(2, 0, color.RGBA{
					R: 0, G: 0, B: 255,
					A: 0xff,
				})
				// yellow
				img.Set(0, 1, color.RGBA{
					R: 255, G: 255, B: 0,
					A: 0xff,
				})
				// white
				img.Set(1, 1, color.RGBA{
					R: 255, G: 255, B: 255,
					A: 0xff,
				})
				// black
				img.Set(2, 1, color.RGBA{
					R: 0, G: 0, B: 0,
					A: 0xff,
				})
				return img
			}},
		{
			name:     "color",
			testFile: path.Join("..", "..", "..", "test", "ppm", "color.ppm"),
			image: func() image.Image {
				nx := 200
				ny := 100
				img := image.NewRGBA(image.Rect(0, 0, nx, ny))
				for j := float64(ny - 1); j >= 0; j-- {
					for i := float64(0); i < float64(nx); i++ {
						r := i / float64(nx)
						g := j / float64(ny)
						b := 0.2
						c := color.RGBA{
							R: uint8(r * 255.99),
							G: uint8(g * 255.99),
							B: uint8(b * 255.99),
							A: 0xff,
						}
						img.Set(int(i), ny-int(j)-1, c)
					}
				}
				return img
			}},
	}

	for _, tt := range testImages {
		t.Run(tt.name, func(t *testing.T) {
			img := tt.image()
			ppmImage, err := output.NewPPMOutput(img)
			require.NoError(t, err)

			t.Logf("Tmp Dir: %v", dir)
			file, err := os.CreateTemp(dir, fmt.Sprintf("test_%s.ppm", tt.name))
			require.NoError(t, err)
			t.Logf("File: %v", file.Name())
			writer := bufio.NewWriter(file)
			x, err := io.Copy(writer, ppmImage)
			assert.Greater(t, x, int64(0))
			t.Logf("Copied %v bytes", x)
			assert.NoError(t, err)
			err = writer.Flush()
			assert.NoError(t, err)

			readFile, err := os.ReadFile(file.Name())
			assert.NoError(t, err)
			assert.NotEmpty(t, readFile)

			bytes, err := os.ReadFile(tt.testFile)
			require.NoError(t, err)
			expectedFile := string(bytes)
			assert.Equal(t, expectedFile, string(readFile))
		})
	}
}

func TestNewPPMOutput(t *testing.T) {
	dir, err := os.MkdirTemp("", "TestNewPPMOutput")
	require.NoError(t, err)
	defer func() {
		_ = os.RemoveAll(dir)
	}()

	// Given
	c := scene.NewCanvas(5, 3)

	// And
	c1 := object.NewColor(1.5, 0, 0)
	c2 := object.NewColor(0, 0.5, 0)
	c3 := object.NewColor(-0.5, 0, 1)

	// When
	c.SetColor(0, 0, c1)
	c.SetColor(2, 1, c2)
	c.Set(4, 2, c3.R, c3.G, c3.B)

	// Then
	img := c.GenerateImg()
	ppmImage, err := output.NewPPMOutput(img)
	require.NoError(t, err)

	all, err := io.ReadAll(ppmImage)
	require.NoError(t, err)

	output := string(all)
	t.Log(output)
}
