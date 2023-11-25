package tag

import (
	"fmt"

	"github.com/spf13/cobra"
)

var t001Cmd = &cobra.Command{
	Use:     "001",
	GroupID: "tag",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(cmd.Use)

		return nil
	},
}
