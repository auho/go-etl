package query

import (
	"github.com/auho/go-etl/v2/insight/app"
	"github.com/spf13/cobra"
)

var _app *app.Application

var queryCmd *cobra.Command

func init() {
	queryCmd = &cobra.Command{
		Use: "query",
	}

	queryCmd.AddGroup(
		&cobra.Group{ID: "query", Title: "query"},
		&cobra.Group{ID: "all", Title: "all"},
	)

	queryCmd.AddCommand()
}

func Initial(parentCmd *cobra.Command) {
	_app = app.APP

	parentCmd.AddCommand(queryCmd)
}
