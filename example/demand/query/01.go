package query

import (
	"fmt"

	"github.com/spf13/cobra"
)

var _query001Cmd = &cobra.Command{
	Use: "query.001",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Use)
	},
}
