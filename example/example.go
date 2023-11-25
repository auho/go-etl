package main

import (
	"log"
	"os"
	"runtime/pprof"

	"github.com/auho/go-etl/v2/example/demand"
	"github.com/auho/go-etl/v2/insight/app"
	"github.com/spf13/cobra"
)

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
	var confName string
	rootCmd.PersistentFlags().StringVarP(&confName, "config", "c", "", "config")

	if confName == "" {
		confName = "develop"
		//panic("conf name is empty")
	}

	// init app
	app.NewApp(confName)
	app.APP.PrintlnState()

	rootCmd.Use = app.APP.Name

	// initial demand
	demand.Initial(rootCmd)
}
