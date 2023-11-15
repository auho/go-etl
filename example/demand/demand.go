package demand

import (
	"github.com/auho/go-etl/v2/example/demand/query"
	"github.com/auho/go-etl/v2/example/demand/tag"
	"github.com/spf13/cobra"
)

func Initial(parentCmd *cobra.Command) {
	parentCmd.AddCommand(_stateCmd)
	parentCmd.AddCommand(query.Cmd)
	parentCmd.AddCommand(tag.Cmd)
}

var _stateCmd = &cobra.Command{
	Use: "state",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
