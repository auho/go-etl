package demand

import (
	"github.com/auho/go-etl/v2/example/demand/build"
	"github.com/auho/go-etl/v2/insight/app"
	"github.com/spf13/cobra"
)

func init() {
	buildCmd.AddCommand(
		_buildAllCmd,
	)
}

var buildCmd = &cobra.Command{
	Use: "build",
}

var _buildAllCmd = &cobra.Command{
	Use: "all",
	RunE: func(cmd *cobra.Command, args []string) error {
		app.APP.AddCommand(build.Build001Cmd)

		return app.APP.RunCommandE(args)
	},
}
