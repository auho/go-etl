package importtodb

import (
	"fmt"
	"log"

	"github.com/auho/go-etl/v2/insight/assistant/model"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
)

func ExampleRunImportToDb() {
	// RunImportToDb
	err := RunImportToDb("one.xlsx",
		// rule
		&RuleResource{
			Resource: Resource{
				SheetName:       "Sheet1",
				StartRow:        2,
				IsRecreateTable: true,
			},
			Titles: Titles{
				Titles: []string{"column1_name"}, // save to db 的 columns; 从第一个 column 开始，连续不间断；此选择优
			},
			Rule: model.NewRuleSimple("rule_name", nil, nil),
		},
		// rows
		&RowsResource{
			Resource: Resource{
				SheetName:       "Sheet2",
				StartRow:        2,
				IsRecreateTable: true,
				CommandFun: func(command *tablestructure.Command) {
					command.AddString("two_1")
					command.AddString("two_2")
				},
				PostFun: func(resource Resourcer) error {
					err1 := resource.GetDB().GormDB().
						Table(resource.GetTable().GetTableName()).
						Where(fmt.Sprintf("`%s` = ?", "two_1"), "value").
						UpdateColumn("two_2", "").Error
					if err1 != nil {
						return fmt.Errorf("update error; %w", err1)
					}
					return nil
				},
			},
			Titles: Titles{
				Titles: []string{"two_1", "two_2"}, // columns name to db
			},
			Rows: model.NewRows("two", "id", nil),
		},
	)

	if err != nil {
		log.Fatalln(err)
	}
}
