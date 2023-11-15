package cmd

import (
	"fmt"
	"os"

	"github.com/auho/go-etl/v2/example/app/application"
	"github.com/auho/go-etl/v2/example/demand"
	"github.com/spf13/cobra"
)

func initial() {
	if confName == "" {
		confName = "develop"
		//panic("conf name is empty")
	}

	application.NewApp(confName)

	ss := application.App.State()
	for _, _s := range ss {
		fmt.Println(_s)
	}

	fmt.Println()

	demand.Initial(rootCmd)
}

var confName string

var rootCmd = &cobra.Command{
	Use: "example",
}

func Execute() {
	rootCmd.PersistentFlags().StringVarP(&confName, "config", "c", "", "config")

	initial()

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
