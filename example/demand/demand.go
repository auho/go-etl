package demand

import (
	"github.com/auho/go-etl/v2/example/demand/layout"
	"github.com/auho/go-etl/v2/insight/app"
	"github.com/spf13/cobra"
)

func Initial(parentCmd *cobra.Command) {
	layout.Initial()

	parentCmd.AddCommand(_stateCmd)
	parentCmd.AddCommand(queryCmd)
	parentCmd.AddCommand(tagCmd)
}

var _stateCmd = &cobra.Command{
	Use: "state",
	Run: func(cmd *cobra.Command, args []string) {
		app.APP.State()
	},
}
