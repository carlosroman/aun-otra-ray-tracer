package cmd

import (
	"bufio"
	"fmt"
	"image/jpeg"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/carlosroman/aun-otra-ray-tracer/go/internal/output"
	"github.com/carlosroman/aun-otra-ray-tracer/go/internal/scene"
)

var (
	samplesPerPixel int16
	filename        string
	lowRes          bool
	nx              int64
	isJpeg          bool
	rootCmd         = &cobra.Command{
		Use:   "example",
		Short: "Some example 3D renders",
		Long:  "A selection of 3D renders using the engine build here",
	}
)

func initConfig() {
	viper.AutomaticEnv()
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
}

func renderScene(c scene.Camera, world scene.World, start time.Time, fname string) (err error) {
	//img := scene.Render(camera, world)
	workerCount := runtime.GOMAXPROCS(0)
	fmt.Println(fmt.Sprintf("Setting no of workers to: %v", workerCount))
	img := scene.MultiThreadedRender(c, world, 8, 1024)
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

func MaxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}
