package _import

import (
	"fmt"
	"strconv"
	"unicode/utf8"

	"github.com/auho/go-etl/v2/insight/model"
	"github.com/auho/go-etl/v2/insight/model/excel/read"
	"github.com/auho/go-etl/v2/insight/model/table"
	simpleDb "github.com/auho/go-simple-db/v2"
)

var _ sourceor = (*RuleImport)(nil)

func RunRuleImport(db *simpleDb.SimpleDB, ri *RuleImport) error {
	return RunImport(db, ri)
}

type RuleImport struct {
	Source
	KeywordIndex int      // keyword 在 sheet 中的列，从 0 开始
	Titles       []string // save to db 的 columns
	Rule         model.Ruler
}

func (ri *RuleImport) GetTable() table.Tabler {
	return table.NewRuleTable(ri.Rule)
}

func (ri *RuleImport) GetTitles() []string {
	return ri.Titles
}

func (ri *RuleImport) GetSheetData() (read.SheetDataor, error) {
	sheetData, err := read.NewSheetDataNoTitle(ri.XlsxPath, ri.SheetName, ri.StartRow)
	if err != nil {
		return nil, fmt.Errorf("NewSheetDataNoTitle error; %w", sheetData)
	}

	err = sheetData.ReadData()
	if err != nil {
		return nil, fmt.Errorf("ReadData error; %w", err)
	}

	// drop duplicates TODO add if
	for i, title := range ri.Titles {
		if title == ri.Rule.KeywordName() {
			ri.ColumnDropDuplicates = append(ri.ColumnDropDuplicates, i)
			break
		}
	}

	// keyword len of string
	err = sheetData.HandlerRows(func(rows [][]string) ([][]string, error) {
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
