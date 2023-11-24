package demand

import (
	"github.com/auho/go-etl/v2/example/demand/tag"
	"github.com/auho/go-etl/v2/insight/app"
	"github.com/spf13/cobra"
)

func init() {
	tagCmd.AddCommand(
		_tagAllCmd,
	)
}

var tagCmd = &cobra.Command{
	Use: "tag",
}

var _tagAllCmd = &cobra.Command{
	Use: "all",
	RunE: func(cmd *cobra.Command, args []string) error {
		app.APP.AddCommand(tag.Tag001Cmd)

		return app.APP.RunCommandE(args)
	},
}
