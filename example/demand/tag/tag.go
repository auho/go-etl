package tag

import (
	"github.com/auho/go-etl/v2/insight/app"
	"github.com/spf13/cobra"
)

var _app *app.Application

var tagCmd *cobra.Command

func init() {
	tagCmd = &cobra.Command{
		Use: "tag",
	}

	tagCmd.AddGroup(
		&cobra.Group{ID: "tag", Title: "tag"},
		&cobra.Group{ID: "all", Title: "all"},
	)

	tagCmd.AddCommand()
}

func Initial(parentCmd *cobra.Command) {
	_app = app.APP

	parentCmd.AddCommand(tagCmd)
}
