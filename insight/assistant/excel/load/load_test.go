package load

import (
	"testing"

	"github.com/auho/go-etl/v3/insight/assistant/entity"
	"github.com/auho/go-etl/v3/insight/assistant/schema"
)

func TestRunLoad_Rows(t *testing.T) {
	skipNoDB(t)

	tableName := "test_load_rows"
	t.Cleanup(func() { dropTable(t, tableName) })

	xlsxPath := createTestExcel(t, "Sheet1", [][]any{
		{"name", "value"},
		{"alice", "100"},
		{"bob", "200"},
	})

	err := RunLoad(xlsxPath,
		&Rows{
			baseResource: baseResource{
				SheetName:       "Sheet1",
				StartRow:        2,
				IsRecreateTable: true,
				CommandFun: func(command *schema.Command) {
					command.AddString("name")
					command.AddString("value")
				},
			},
			Titles: Titles{
				Titles: []string{"name", "value"},
			},
			Rows: entity.NewRows(tableName, "id", _simpleDB),
		},
	)
	if err != nil {
		t.Fatalf("RunLoad: %v", err)
	}

	var count int64
	err = _simpleDB.GormDB().Table(tableName).Count(&count).Error
	if err != nil {
		t.Fatalf("Count: %v", err)
	}
	if count != 2 {
		t.Errorf("row count = %d, want 2", count)
	}
}

func TestRunLoad_Rule(t *testing.T) {
	skipNoDB(t)

	ruleName := "test_load_rule"
	tableName := "rule_" + ruleName
	t.Cleanup(func() { dropTable(t, tableName) })

	rule := entity.NewRuleSimple(ruleName, []string{"label_a"}, _simpleDB)

	nameCol := rule.Name()           // test_load_rule
	labelCol := rule.LabelsName()[0] // label_a
	keywordCol := rule.KeywordName() // test_load_rule_keyword

	xlsxPath := createTestExcel(t, "Sheet1", [][]any{
		{nameCol, labelCol, keywordCol},
		{"rule1", "labelA", "apple"},
		{"rule2", "labelB", "car"},
	})

	err := RunLoad(xlsxPath,
		&Rule{
			baseResource: baseResource{
				SheetName:       "Sheet1",
				StartRow:        2,
				IsRecreateTable: true,
			},
			Titles: Titles{
				Titles: []string{nameCol, labelCol, keywordCol},
			},
			Rule: rule,
		},
	)
	if err != nil {
		t.Fatalf("RunLoad: %v", err)
	}

	var count int64
	err = _simpleDB.GormDB().Table(tableName).Count(&count).Error
	if err != nil {
		t.Fatalf("Count: %v", err)
	}
	if count != 2 {
		t.Errorf("row count = %d, want 2", count)
	}
}

func TestRunLoad_Raw(t *testing.T) {
	skipNoDB(t)

	tableName := "test_load_raw"
	t.Cleanup(func() { dropTable(t, tableName) })

	xlsxPath := createTestExcel(t, "Sheet1", [][]any{
		{"col1", "col2"},
		{"value1", "value2"},
		{"value3", "value4"},
	})

	err := RunLoad(xlsxPath,
		&Raw{
			baseResource: baseResource{
				SheetName:       "Sheet1",
				StartRow:        1,
				IsRecreateTable: true,
				CommandFun: func(command *schema.Command) {
					command.AddString("col1")
					command.AddString("col2")
				},
			},
			Rows: entity.NewRows(tableName, "id", _simpleDB),
		},
	)
	if err != nil {
		t.Fatalf("RunLoad: %v", err)
	}

	var count int64
	err = _simpleDB.GormDB().Table(tableName).Count(&count).Error
	if err != nil {
		t.Fatalf("Count: %v", err)
	}
	if count != 2 {
		t.Errorf("row count = %d, want 2", count)
	}
}

func TestRunLoad_AppendData(t *testing.T) {
	skipNoDB(t)

	tableName := "test_load_append"
	t.Cleanup(func() { dropTable(t, tableName) })

	xlsxPath := createTestExcel(t, "Sheet1", [][]any{
		{"name", "value"},
		{"alice", "100"},
		{"bob", "200"},
	})

	rowsEntity := entity.NewRows(tableName, "id", _simpleDB)

	// first import: create table and insert 2 rows
	err := RunLoad(xlsxPath,
		&Rows{
			baseResource: baseResource{
				SheetName:       "Sheet1",
				StartRow:        2,
				IsRecreateTable: true,
				CommandFun: func(command *schema.Command) {
					command.AddString("name")
					command.AddString("value")
				},
			},
			Titles: Titles{
				Titles: []string{"name", "value"},
			},
			Rows: rowsEntity,
		},
	)
	if err != nil {
		t.Fatalf("first RunLoad: %v", err)
	}

	// second import: append data (not recreate, append)
	err = RunLoad(xlsxPath,
		&Rows{
			baseResource: baseResource{
				SheetName:       "Sheet1",
				StartRow:        2,
				IsRecreateTable: false,
				IsAppendData:    true,
				CommandFun: func(command *schema.Command) {
					command.AddString("name")
					command.AddString("value")
				},
			},
			Titles: Titles{
				Titles: []string{"name", "value"},
			},
			Rows: entity.NewRows(tableName, "id", _simpleDB),
		},
	)
	if err != nil {
		t.Fatalf("second RunLoad: %v", err)
	}

	var count int64
	err = _simpleDB.GormDB().Table(tableName).Count(&count).Error
	if err != nil {
		t.Fatalf("Count: %v", err)
	}
	if count != 4 {
		t.Errorf("row count = %d, want 4 (2 + 2 appended)", count)
	}
}
