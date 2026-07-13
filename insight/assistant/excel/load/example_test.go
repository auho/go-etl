package load

import (
	"fmt"
	"log"

	"github.com/auho/go-etl/v3/insight/assistant/entity"
	"github.com/auho/go-etl/v3/insight/assistant/schema"
)

func ExampleRunLoad() {
	// RunLoad
	err := RunLoad("one.xlsx",
		// rule
		&Rule{
			baseResource: baseResource{
				SheetName:       "Sheet1",
				StartRow:        2,
				IsRecreateTable: true,
			},
			Titles: Titles{
				Titles: []string{"column1_name"}, // save to db 的 columns; 从第一个 column 开始，连续不间断；此选择优
			},
			Rule: entity.NewRuleSimple("rule_name", nil, nil),
		},
		// rows
		&Rows{
			baseResource: baseResource{
				SheetName:       "Sheet2",
				StartRow:        2,
				IsRecreateTable: true,
				CommandFun: func(command *schema.Command) {
					command.AddString("two_1")
					command.AddString("two_2")
				},
				AfterFun: func(resource Resource) error {
					err1 := resource.DB().GormDB().
						Table(resource.Tabler().GetTableName()).
						Where(fmt.Sprintf("`%s` = ?", "two_1"), "value").
						UpdateColumn("two_2", "").Error
					if err1 != nil {
						return fmt.Errorf("UpdateColumn: %w", err1)
					}
					return nil
				},
			},
			Titles: Titles{
				Titles: []string{"two_1", "two_2"}, // columns name to db
			},
			Rows: entity.NewRows("two", "id", nil),
		},
	)

	if err != nil {
		log.Fatalln(err)
	}
}
