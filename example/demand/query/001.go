package query

import (
	"fmt"

	"github.com/spf13/cobra"
)

var q001Cmd = &cobra.Command{
	Use:     "001",
	GroupID: "query",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Use)
	},
}
