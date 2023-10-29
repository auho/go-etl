package cmd

import (
	"os"

	"github.com/auho/go-etl/v2/insight/app"
	"github.com/auho/go-etl/v2/insight/example/demand"
	"github.com/spf13/cobra"
)

func initial() {
	demand.Initial(rootCmd)
}

var confName string

var rootCmd = &cobra.Command{
	Use: "example",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if confName == "" {
			panic("conf name is empty")
		}

		demand.App = app.NewApp(confName)
	},
}

func Execute() {
	rootCmd.PersistentFlags().StringVarP(&confName, "config", "c", "", "config")

	initial()

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
