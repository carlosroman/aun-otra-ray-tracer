package cmd

import (
	"fmt"
	"math"
	"time"

	"github.com/spf13/cobra"

	"github.com/carlosroman/aun-otra-ray-tracer/go/internal/object"
	"github.com/carlosroman/aun-otra-ray-tracer/go/internal/ray"
	"github.com/carlosroman/aun-otra-ray-tracer/go/internal/scene"
)

var (
	hexagonFilename string
)

func init() {
	// TODO: Move flags to root/global
	hexagonCmd.Flags().BoolVarP(&isJpeg, "jpeg", "j", false, "Switch output to jpeg")
	hexagonCmd.Flags().StringVarP(&hexagonFilename, "filename", "f", "cube", "Filename of the output")
	hexagonCmd.Flags().Int16VarP(&samplesPerPixel, "samples", "s", 4, "Number of samples per pixel")
	hexagonCmd.Flags().Int64VarP(&nx, "width", "w", 640, "Image width in pixels")
	rootCmd.AddCommand(hexagonCmd)
}

var hexagonCmd = &cobra.Command{
	Use:   "hexagon",
	Short: "Render a hexagon",
	Long:  `This renders a simple hexagon made up of spheres and cylinders.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		start := time.Now()
		var scale = float64(1) / float64(samplesPerPixel)
		var ny = float64(nx) / ratio
		fmt.Println(fmt.Sprintf("Generating image %v,%v with samples: %v and scale: %v", nx, ny, samplesPerPixel, scale))

		var camera scene.Camera
		camera, err = getBasicCamera(int(nx), int(ny))
		if err != nil {
			return err
		}
		world, err := getBasicRoom()
		if err != nil {
			return err
		}
		hexMaterial := object.DefaultMaterial()
		hexMaterial.Pattern = object.NewCheckerPattern(object.Red, object.Green)
		hex := object.NewHexagon(
			object.WithMaterial(hexMaterial),
			object.WithTransform(
				ray.Translation(-0.5, 1, 0.5).Multiply(
					ray.Rotation(ray.X, -math.Pi/12))),
		)
		world.AddObjects(hex)
		return renderScene(camera, world, start, hexagonFilename)
	},
}
