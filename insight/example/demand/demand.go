package demand

import (
	"github.com/auho/go-etl/insight/app"
	"github.com/spf13/cobra"
)

var App *app.App

func Initial(parentCmd *cobra.Command) {
	parentCmd.AddCommand(_01Cmd)
}
