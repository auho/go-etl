package tag

import (
	"fmt"

	"github.com/spf13/cobra"
)

var _tag001Cmd = &cobra.Command{
	Use: "tag.001",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Use)
	},
}
