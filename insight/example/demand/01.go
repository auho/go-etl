package demand

import "github.com/spf13/cobra"

var _01Cmd = &cobra.Command{
	Use: "01",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println(01)
	},
}
