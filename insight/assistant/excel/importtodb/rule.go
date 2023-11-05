package importtodb

import (
	"fmt"
	"strconv"
	"unicode/utf8"

	"github.com/auho/go-etl/v2/insight/assistant/excel/read"
	"github.com/auho/go-etl/v2/insight/assistant/model"
	"github.com/auho/go-etl/v2/insight/assistant/table"
	simpleDb "github.com/auho/go-simple-db/v2"
)

var _ resourcer = (*RuleResource)(nil)

type RuleResource struct {
	Resource
	KeywordIndex int      // keyword 在 sheet 中的列，从 0 开始
	Titles       []string // save to db 的 columns
	Rule         model.Ruler
}

func (re *RuleResource) GetTable() table.Tabler {
	return table.NewRuleTable(re.Rule)
}

func (re *RuleResource) GetTitles() []string {
	return re.Titles
}

func (re *RuleResource) GetSheetData() (read.SheetDataor, error) {
	sheetData, err := read.NewSheetDataNoTitle(re.XlsxPath, re.SheetName, re.StartRow)
	if err != nil {
		return nil, fmt.Errorf("NewSheetDataNoTitle error; %w", err)
	}

	err = sheetData.ReadData()
	if err != nil {
		return nil, fmt.Errorf("ReadData error; %w", err)
	}

	// drop duplicates TODO add if
	for i, title := range re.Titles {
		if title == re.Rule.KeywordName() {
			re.ColumnDropDuplicates = append(re.ColumnDropDuplicates, i)
			break
		}
	}

	// keyword len of string
	err = sheetData.HandlerRows(func(rows [][]string) ([][]string, error) {
		re.Titles = append(re.Titles, re.Rule.KeywordLenName())

		for i, row := range rows {
			row = append(row, strconv.Itoa(utf8.RuneCountInString(row[re.KeywordIndex])))
			rows[i] = row
		}

		return rows, nil
	})

	if err != nil {
		return nil, fmt.Errorf("keyWord len error; %w", err)
	}

	return sheetData, nil
}

func (re *RuleResource) Import(db *simpleDb.SimpleDB) error {
	return RunImportToDb(db, re)
}
