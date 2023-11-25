package demand

import (
	"github.com/auho/go-etl/v2/example/demand/build"
	"github.com/auho/go-etl/v2/example/demand/layout"
	"github.com/auho/go-etl/v2/example/demand/query"
	"github.com/auho/go-etl/v2/example/demand/tag"
	"github.com/auho/go-etl/v2/insight/app"
	"github.com/spf13/cobra"
)

var _app *app.Application

func Initial(parentCmd *cobra.Command) {
	_app = app.APP

	parentCmd.AddGroup(&cobra.Group{ID: "all", Title: "all"})

	parentCmd.AddCommand(_stateCmd)

	_app.PreFun(layout.Initial)

	build.Initial(parentCmd)
	tag.Initial(parentCmd)
	query.Initial(parentCmd)

	// other initial
}

var _stateCmd = &cobra.Command{
	Use: "state",
	Run: func(cmd *cobra.Command, args []string) {
		_app.State()
	},
}
