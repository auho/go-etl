package demand

import (
	"github.com/auho/go-etl/v2/flow"
	"github.com/auho/go-etl/v2/means/segword"
	"github.com/auho/go-etl/v2/mode"
	"github.com/spf13/cobra"
)

var _01Cmd = &cobra.Command{
	Use: "01",
	Run: func(cmd *cobra.Command, args []string) {
		flow.InsertFlow(
			App.DB,
			"goods",
			"id",
			"goods_product_name_words",
			mode.NewInsert([]string{"product_name"}, segword.NewSegWordsMeans()),
			[]string{""},
		)
	},
}
