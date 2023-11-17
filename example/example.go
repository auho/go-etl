package main

import (
	"os"

	"github.com/auho/go-etl/v2/example/demand"
	"github.com/auho/go-etl/v2/insight/app"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use: "root",
	}

	initial(rootCmd)

	err := rootCmd.Execute()
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

	app.NewApp(confName)
	app.APP.PrintlnState()

	rootCmd.Use = app.APP.Name

	demand.Initial(rootCmd)
}
