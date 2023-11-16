package main

import (
	"fmt"
	"os"

	"github.com/auho/go-etl/v2/example/demand"
	"github.com/auho/go-etl/v2/insight/app"
	"github.com/spf13/cobra"
)

func initial(confName string, rootCmd *cobra.Command) {
	if confName == "" {
		confName = "develop"
		//panic("conf name is empty")
	}

	app.NewApp(confName)

	ss := app.APP.State()
	for _, _s := range ss {
		fmt.Println(_s)
	}

	fmt.Println()

	demand.Initial(rootCmd)
}

func main() {
	var confName string

	var rootCmd = &cobra.Command{
		Use: "example",
	}

	rootCmd.PersistentFlags().StringVarP(&confName, "config", "c", "", "config")

	initial(confName, rootCmd)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
