package build

import (
	"github.com/auho/go-etl/v2/insight/app"
	"github.com/spf13/cobra"
)

var _app *app.Application

var buildCmd *cobra.Command

func init() {
	buildCmd = &cobra.Command{
		Use: "build",
	}

	buildCmd.AddGroup(
		&cobra.Group{ID: "data", Title: "data"},
		&cobra.Group{ID: "rule", Title: "rule"},
		&cobra.Group{ID: "all", Title: "all"},
	)

	buildCmd.AddCommand()
}

func Initial(parentCmd *cobra.Command) {
	_app = app.APP

	parentCmd.AddCommand(buildCmd)
}
