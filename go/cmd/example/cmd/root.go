package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
