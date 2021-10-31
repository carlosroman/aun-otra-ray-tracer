package cmd

import (
	"fmt"
	"math"
	"time"

	"github.com/carlosroman/aun-otra-ray-tracer/go/internal/object"
	"github.com/carlosroman/aun-otra-ray-tracer/go/internal/ray"
	"github.com/carlosroman/aun-otra-ray-tracer/go/internal/scene"
	"github.com/spf13/cobra"
)

var (
	cubeFilename string
)

func init() {
	// TODO: Move flags to root/global
	cubesCmd.Flags().BoolVarP(&isJpeg, "jpeg", "j", false, "Switch output to jpeg")
	cubesCmd.Flags().StringVarP(&cubeFilename, "filename", "f", "cube", "Filename of the output")
	cubesCmd.Flags().Int16VarP(&samplesPerPixel, "samples", "s", 4, "Number of samples per pixel")
	cubesCmd.Flags().Int64VarP(&nx, "width", "w", 640, "Image width in pixels")
	rootCmd.AddCommand(cubesCmd)
}

var cubesCmd = &cobra.Command{
	Use:   "cubes",
	Short: "Render a basic scene",
	Long:  `This renders the renders a simple 3D scene made up of basic objects.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		start := time.Now()
		var scale = float64(1) / float64(samplesPerPixel)
		var ny = float64(nx) / ratio
		fmt.Println(fmt.Sprintf("Generating image %v,%v with samples: %v and scale: %v", nx, ny, samplesPerPixel, scale))

		var world = scene.NewWorld()
		var camera scene.Camera
		camera, err = getBasicCamera(int(nx), int(ny))
		if err != nil {
			return err
		}
		world, err = getBasicRoom()
		if err != nil {
			return err
		}

		middle := object.DefaultSphere()
		err = middle.SetTransform(ray.Translation(-0.5, 1, 0.5))
		if err != nil {
			return err
		}
		middleM := middle.Material()
		middleM.Color = object.NewColor(0.1, 1, 0.5)
		middleM.Diffuse = 0.7
		middleM.Specular = 0.3
		middle.SetMaterial(middleM)

		right := object.DefaultSphere()
		err = right.SetTransform(ray.Translation(1.5, 0.5, -0.5).Multiply(ray.Scaling(0.5, 0.5, 0.5)))
		if err != nil {
			return err
		}
		rightM := right.Material()
		rightM.Color = object.NewColor(0.5, 1, 0.1)
		rightM.Diffuse = 0.7
		rightM.Specular = 0.3
		rightM.Reflective = 1.0
		right.SetMaterial(rightM)

		left := object.NewCube()
		err = left.SetTransform(ray.Translation(-1.5, 0.33, -0.75).Multiply(ray.Scaling(0.33, 0.33, 0.33)))
		if err != nil {
			return err
		}
		leftM := left.Material()
		leftM.Pattern = object.NewCheckerPattern(object.Red, object.Green)
		err = leftM.Pattern.SetTransform(ray.Translation(1.5, 0.5, -0.5).Multiply(ray.Scaling(0.5, 0.5, 0.5)))
		if err != nil {
			return err
		}
		leftM.Color = object.NewColor(1, 0.8, 0.1)
		leftM.Diffuse = 0.7
		leftM.Specular = 0.3
		left.SetMaterial(leftM)

		world.AddObjects(middle, right, left)
		return renderScene(camera, world, start, cubeFilename)
	},
}

func getBasicCamera(nx, ny int) (camera scene.Camera, err error) {
	camera, err = scene.NewBasicCamera(nx, ny, math.Pi/3)
	if err != nil {
		return nil, err
	}
	err = camera.SetTransform(
		ray.ViewTransform(
			ray.NewPoint(0, 1.5, -5), ray.NewPoint(0, 1, 0), ray.NewVec(0, 1, 0)))
	return camera, err
}

func getBasicRoom() (world scene.World, err error) {
	floor := object.NewPlane()
	err = floor.SetTransform(ray.Scaling(10, 0.01, 10))
	if err != nil {
		return nil, err
	}
	floorM := floor.Material()
	floorM.Color = object.NewColor(1, 0.9, 0.9)
	floorM.Specular = 0
	floorM.Pattern = object.NewCheckerPattern(object.White, object.Black)
	err = floorM.Pattern.SetTransform(ray.Scaling(0.1, 0.01, 0.1))
	if err != nil {
		return nil, err
	}
	floor.SetMaterial(floorM)

	leftWall := object.NewPlane()
	err = leftWall.SetTransform(
		ray.Translation(0, 0, 5).
			Multiply(ray.Rotation(ray.Y, -math.Pi/4)).
			Multiply(ray.Rotation(ray.X, math.Pi/2)).
			Multiply(ray.Scaling(10, 0.01, 10)))
	if err != nil {
		return nil, err
	}
	leftWall.SetMaterial(floorM)

	rightWall := object.NewPlane()
	err = rightWall.SetTransform(
		ray.Translation(0, 0, 5).
			Multiply(ray.Rotation(ray.Y, math.Pi/4)).
			Multiply(ray.Rotation(ray.X, math.Pi/2)).
			Multiply(ray.Scaling(10, 0.01, 10)))
	if err != nil {
		return nil, err
	}
	rightWall.SetMaterial(floorM)

	world = scene.NewWorld()
	world.AddLight(object.NewPointLight(ray.NewPoint(-10, 10, -10), object.White))
	world.AddObjects(floor, leftWall, rightWall)
	return world, err
}
