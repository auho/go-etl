package tag

import (
	"fmt"

	"github.com/spf13/cobra"
)

var t001Cmd = &cobra.Command{
	Use:     "001",
	GroupID: "tag",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Use)
	},
}
