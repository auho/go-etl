package build

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Build001Cmd = &cobra.Command{
	Use: "001",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(cmd.Use)

		return nil
	},
}
