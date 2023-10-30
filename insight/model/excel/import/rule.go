package _import

import (
	"fmt"
	"strconv"
	"unicode/utf8"

	"github.com/auho/go-etl/v2/insight/model"
	"github.com/auho/go-etl/v2/insight/model/excel/read"
	simpleDb "github.com/auho/go-simple-db/v2"
)

type RuleImport struct {
	XlsxPath     string
	SheetName    string
	StartRowNo   int      // 数据开始的行数， 从 0 开始
	KeywordIndex int      // keyword 在 sheet 中的列，从 0 开始
	Titles       []string // save to db 的 columns
	Rule         *model.Rule
	NeedTruncate bool
}

func RunRuleImport(db *simpleDb.SimpleDB, ri *RuleImport) error {
	sheetData, err := read.NewSheetDataNoTitle(ri.XlsxPath, ri.SheetName)
	if err != nil {
		return fmt.Errorf("NewSheetDataNoTitle error; %w", sheetData)
	}

	err = sheetData.ReadData()
	if err != nil {
		return fmt.Errorf("ReadData error; %w", err)
	}

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
		return fmt.Errorf("HandlerRows error; %w", err)
	}

	err = db.Truncate(ri.Rule.TableName())
	if err != nil {
		return fmt.Errorf("truncate error; %w", err)
	}

	err = db.BulkInsertFromSliceSlice(ri.Rule.TableName(), ri.Titles, sheetData.GetRowsWithAny(), 1000)
	if err != nil {
		return err
	}

	return nil
}
