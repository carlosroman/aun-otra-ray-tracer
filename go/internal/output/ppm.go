package output

import (
	"bytes"
	"fmt"
	"image"
	"io"
)

func NewPPMOutput(img image.Image) (r io.Reader, err error) {
	var b bytes.Buffer
	bounds := img.Bounds()
	_, err = b.WriteString("P3\n")
	if err != nil {
		return nil, err
	}

	_, err = b.WriteString(fmt.Sprintf("%v %v\n", bounds.Dx(), bounds.Dy()))
	if err != nil {
		return nil, err
	}

	_, err = b.WriteString("255\n")
	if err != nil {
		return nil, err
	}
	for y := bounds.Dy() - 1; y >= 0; y-- {
		for x := 0; x < bounds.Dx(); x++ {
			color := img.At(x, y)
			red, green, blue, _ := color.RGBA()
			_, err := b.WriteString(fmt.Sprintf("%v %v %v\n", uint8(red>>8), uint8(green>>8), uint8(blue>>8)))
			if err != nil {
				return nil, err
			}
		}
	}
	return &b, err
}
