package query

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	Cmd.AddCommand(
		_allCmd,
	)
}

var Cmd = &cobra.Command{
	Use: "query",
}

var _allCmd = &cobra.Command{
	Use: "all",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(cmd.Use)

		return nil
	},
}
