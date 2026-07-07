package main

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"github.com/auho/go-etl/v3/example/demand"
	"github.com/auho/go-etl/v3/insight/app"
	"github.com/spf13/cobra"
)

var env = "develop"
var version string
var buildInfo string
var confName string

func main() {
	cpuFile, err := os.Create("cpu.pprof")
	if err != nil {
		log.Fatalln(err)
	}

	memFile, err := os.Create("mem.pprof")
	if err != nil {
		log.Fatalln(err)
	}

	err = pprof.StartCPUProfile(cpuFile)
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		pprof.StopCPUProfile()
		_ = pprof.WriteHeapProfile(memFile)
	}()

	// root cmd
	var rootCmd = &cobra.Command{
		Use: "root",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return app.App.RunPreRunE(cmd)
		},
	}

	// initial root cmd
	initial(rootCmd)

	// execute root cmd
	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func initial(rootCmd *cobra.Command) {
	// init app
	app.NewApp()

	rootCmd.Use = app.App.Name
	rootCmd.PersistentFlags().StringVarP(&confName, "config", "c", "", "config")

	fmt.Println("env:", env)
	fmt.Println("version:", version)
	fmt.Println("build info:", buildInfo)
	fmt.Println()

	// app start
	app.App.AddPreRunE(func() error {
		_confName := ""
		if confName != "" {
			_confName = confName
		} else {
			if app.App.ConfName != "" {
				_confName = app.App.ConfName
			} else {
				_confName = env
			}
		}

		app.App.Build(_confName)
		app.App.PrintlnState()

		return nil
	})

	// initial demand
	demand.Initial(rootCmd)
}
