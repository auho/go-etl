package query

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Query001Cmd = &cobra.Command{
	Use: "001",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Use)
	},
}
