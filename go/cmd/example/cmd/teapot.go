package cmd

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"image/jpeg"
	"io"
	"math"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/carlosroman/aun-otra-ray-tracer/go/internal/object"
	"github.com/carlosroman/aun-otra-ray-tracer/go/internal/output"
	"github.com/carlosroman/aun-otra-ray-tracer/go/internal/ray"
	"github.com/carlosroman/aun-otra-ray-tracer/go/internal/scene"
)

const (
	ratio = float64(16) / 9
)

var (
	//go:embed models/utah-teapot.obj
	teapotHiRes []byte

	//go:embed models/utah-teapot-low.obj
	teapotLowRes []byte
)

func init() {
	// TODO: Move flags to root/global
	teapotCmd.Flags().BoolVarP(&isJpeg, "jpeg", "j", false, "Switch output to jpeg")
	teapotCmd.Flags().StringVarP(&filename, "filename", "f", "example", "Filename of the output")
	teapotCmd.Flags().Int16VarP(&samplesPerPixel, "samples", "s", 4, "Number of samples per pixel")
	teapotCmd.Flags().Int64VarP(&nx, "width", "w", 640, "Image width in pixels")

	teapotCmd.Flags().BoolVarP(&lowRes, "low-res", "l", false, "Select between low res and high rest")
	rootCmd.AddCommand(teapotCmd)
}

var teapotCmd = &cobra.Command{
	Use:   "teapot",
	Short: "Render a teapot",
	Long:  `This renders the Utah teapot, https://en.wikipedia.org/wiki/Utah_teapot, in front of a simple background.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		start := time.Now()
		var scale = float64(1) / float64(samplesPerPixel)
		var ny = float64(nx) / ratio
		fmt.Println(fmt.Sprintf("Generating image %v,%v with samples: %v and scale: %v", nx, ny, samplesPerPixel, scale))

		var teapot = teapotHiRes
		if lowRes {
			teapot = teapotLowRes
		}

		var world = scene.NewWorld()
		//world.AddLight(object.NewPointLight(ray.NewPoint(50, 100, 20), object.NewColor(.5, .5, .5)))
		world.AddLight(object.NewPointLight(ray.NewPoint(2, 50, 100), object.NewColor(.5, .5, .5)))

		// Floor
		planesMaterial := object.DefaultMaterial()
		checkered := object.NewCheckerPattern(
			object.NewColor(0.35, 0.35, 0.35),
			object.NewColor(0.4, 0.4, 0.4))
		planesMaterial.Pattern = checkered
		planesMaterial.Ambient = 1
		planesMaterial.Diffuse = 0
		planesMaterial.Specular = 0
		planesMaterial.Reflective = 0.1

		p := object.NewPlane(
			object.WithMaterial(planesMaterial),
		)

		p2 := object.NewPlane(
			object.WithMaterial(planesMaterial),
			object.WithTransform(ray.Translation(0, 0, -10).Multiply(
				ray.Rotation(ray.X, math.Pi/2))),
		)

		teapotMaterial := object.DefaultMaterial()
		teapotMaterial.Color = object.NewColor(1, 0.3, 0.2)
		teapotMaterial.Shininess = 5
		teapotMaterial.Specular = 0.4

		var obj object.Object
		obj, err = loadTeapot(teapot,
			object.WithMaterial(teapotMaterial),
			object.WithTransform(
				ray.Translation(0, 0, 0).Multiply(
					ray.Rotation(ray.Y, math.Pi*23/22).Multiply(
						ray.Rotation(ray.X, -math.Pi/2).Multiply(
							ray.Scaling(0.3, 0.3, 0.3))))))
		if err != nil {
			return err
		}

		world.AddObjects(p, p2, obj)

		var c scene.Camera
		c, err = getCamera(int(nx), int(ny))
		if err != nil {
			fmt.Println(err)
			return err
		}
		return renderScene(c, world, start, filename)
	},
}

func renderScene(c scene.Camera, world scene.World, start time.Time, fname string) (err error) {
	//img := scene.Render(camera, world)
	img := scene.MultiThreadedRender(c, world, 24, 1024)
	generateImg := img.GenerateImg()

	var outFile *os.File
	if isJpeg {
		outFile, err = os.Create(fmt.Sprintf("%s.jpg", fname))
		if err != nil {
			fmt.Println(err)
			return err
		}

		err = jpeg.Encode(outFile, generateImg, nil)
		if err != nil {
			fmt.Println(err)
			return err
		}
	} else {

		ppmImage, err := output.NewPPMOutput(generateImg)
		if err != nil {
			fmt.Println(err)
			return err
		}
		outFile, err = os.Create(fmt.Sprintf("%s.ppm", fname))
		if err != nil {
			fmt.Println(err)
			return err
		}
		writer := bufio.NewWriter(outFile)
		_, err = io.Copy(writer, ppmImage)
		if err != nil {
			fmt.Println(err)
			return err
		}
		if err = writer.Flush(); err != nil {
			return err
		}
	}

	err = outFile.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}

	elapsed := time.Since(start)
	fmt.Println(fmt.Sprintf("Wrote: %s in %s", outFile.Name(), elapsed))
	return err
}

func loadTeapot(tp []byte, opts ...object.Option) (o object.Object, err error) {
	reader := bytes.NewReader(tp)
	wavObj, err := object.NewWavefrontObj(reader, opts...)
	if err != nil {
		return nil, err
	}

	o = wavObj.Object()
	return o, err
}

func getCamera(nx, ny int) (camera scene.Camera, err error) {
	camera, err = scene.NewBasicCamera(nx, ny, math.Pi/3)
	if err != nil {
		return nil, err
	}
	err = camera.SetTransform(
		ray.ViewTransform(
			ray.NewPoint(0, 7, 13), ray.NewPoint(0, 1, 0), ray.NewVec(0, 1, 0)))
	return camera, err
}
