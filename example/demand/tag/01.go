package tag

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Tag001Cmd = &cobra.Command{
	Use: "001",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Use)
	},
}
