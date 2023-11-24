package demand

import (
	"github.com/auho/go-etl/v2/example/demand/query"
	"github.com/auho/go-etl/v2/insight/app"
	"github.com/spf13/cobra"
)

func init() {
	queryCmd.AddCommand(
		_queryAllCmd,
	)
}

var queryCmd = &cobra.Command{
	Use: "query",
}

var _queryAllCmd = &cobra.Command{
	Use: "all",
	RunE: func(cmd *cobra.Command, args []string) error {
		app.APP.AddCommand(query.Query001Cmd)

		return app.APP.RunCommandE(args)
	},
}
