package query

import (
	"fmt"

	"github.com/spf13/cobra"
)

var q001Cmd = &cobra.Command{
	Use:     "001",
	GroupID: "query",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(cmd.Use)
		
		return nil
	},
}
