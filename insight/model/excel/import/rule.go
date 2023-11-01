package _import

import (
	"fmt"
	"strconv"
	"unicode/utf8"

	"github.com/auho/go-etl/v2/insight/model"
	"github.com/auho/go-etl/v2/insight/model/excel/read"
	"github.com/auho/go-etl/v2/insight/model/table"
	"github.com/auho/go-etl/v2/insight/model/tool/slices"
	simpleDb "github.com/auho/go-simple-db/v2"
)

func RunRuleImport(db *simpleDb.SimpleDB, ri *RuleImport) error {
	return ri.Import(db)
}

type RuleImport struct {
	XlsxPath             string
	SheetName            string
	StartRowNo           int      // 数据开始的行数， 从 0 开始
	KeywordIndex         int      // keyword 在 sheet 中的列，从 0 开始
	IsRecreateTable      bool     // 是否 recreate table
	IsAppendData         bool     // 是否 append data
	IsShowSql            bool     // 是否显示 sql
	ColumnDropDuplicates []int    // drop duplicates for column
	Titles               []string // save to db 的 columns
	Rule                 model.Ruler
}

func (ri *RuleImport) Import(db *simpleDb.SimpleDB) error {
	err := ri.buildTable(db)
	if err != nil {
		return fmt.Errorf("buildTable error; %w", err)
	}
	sheetData, err := ri.readSheetData()
	if err != nil {
		return fmt.Errorf("readSheetData error; %w", err)
	}

	err = ri.importToTable(db, sheetData)
	if err != nil {
		return fmt.Errorf("importToTable error; %w", err)
	}

	return nil
}

func (ri *RuleImport) buildTable(db *simpleDb.SimpleDB) error {
	ruleTable := table.NewRuleTable(ri.Rule)

	if ri.IsShowSql {
		fmt.Println(ruleTable.GetTable().SqlForCreate())
	}

	_, err := db.GetTableColumns(ri.Rule.TableName())
	if err != nil {
		ri.IsRecreateTable = true
	} else {
		if ri.IsRecreateTable {
			err = db.Drop(ri.Rule.TableName())
			if err != nil {
				return fmt.Errorf("drop; %w", err)
			}
		}
	}

	if ri.IsRecreateTable {
		err = ruleTable.Build(db)
		if err != nil {
			return fmt.Errorf("build error; %w", err)
		}
	}

	return nil
}

func (ri *RuleImport) readSheetData() (*read.SheetDataNoTitle, error) {
	sheetData, err := read.NewSheetDataNoTitle(ri.XlsxPath, ri.SheetName)
	if err != nil {
		return nil, fmt.Errorf("NewSheetDataNoTitle error; %w", sheetData)
	}

	err = sheetData.ReadData()
	if err != nil {
		return nil, fmt.Errorf("ReadData error; %w", err)
	}

	// drop duplicates
	if len(ri.ColumnDropDuplicates) <= 0 {
		for i, title := range ri.Titles {
			if title == ri.Rule.KeywordName() {
				ri.ColumnDropDuplicates = append(ri.ColumnDropDuplicates, i)
				break
			}
		}
	}

	if len(ri.ColumnDropDuplicates) > 0 {
		err = sheetData.HandlerRows(func(rows [][]string) ([][]string, error) {
			rows = slices.SliceSliceDropDuplicates(rows, ri.ColumnDropDuplicates)

			return rows, nil
		})
		if err != nil {
			return nil, fmt.Errorf("drop duplicates error; %w", err)
		}
	}

	// keyword len of string
	err = sheetData.HandlerRows(func(rows [][]string) ([][]string, error) {
		rows = rows[ri.StartRowNo:]

		ri.Titles = append(ri.Titles, ri.Rule.KeywordLenName())

		for i, row := range rows {
			row = append(row, strconv.Itoa(utf8.RuneCountInString(row[ri.KeywordIndex])))
			rows[i] = row
		}

		return rows, nil
	})

	if err != nil {
		return nil, fmt.Errorf("keyWord len error; %w", err)
	}

	return sheetData, nil
}

func (ri *RuleImport) importToTable(db *simpleDb.SimpleDB, sheetData *read.SheetDataNoTitle) error {
	var err error

	if !ri.IsAppendData {
		err = db.Truncate(ri.Rule.TableName())
		if err != nil {
			return fmt.Errorf("truncate error; %w", err)
		}
	}

	err = db.BulkInsertFromSliceSlice(ri.Rule.TableName(), ri.Titles, sheetData.GetRowsWithAny(), 1000)
	if err != nil {
		return err
	}

	return nil
}
