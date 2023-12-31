package main

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"github.com/auho/go-etl/v2/example/demand"
	"github.com/auho/go-etl/v2/insight/app"
	"github.com/spf13/cobra"
)

var env = "develop"
var version string
var lastDate string
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
			if confName == "" {
				confName = env
				//panic("conf name is empty")
			}

			app.APP.Build(confName)
			app.APP.PrintlnState()
			err1 := app.APP.RunPreRunE(cmd)
			if err1 != nil {
				return err1
			}

			return nil
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
	rootCmd.PersistentFlags().StringVarP(&confName, "config", "c", "", "config")

	// init app
	app.NewApp()

	rootCmd.Use = app.APP.Name

	fmt.Println("env:", env)
	fmt.Println("version:", version)
	fmt.Println("last date:", lastDate)
	fmt.Println()

	// initial demand
	demand.Initial(rootCmd)
}
