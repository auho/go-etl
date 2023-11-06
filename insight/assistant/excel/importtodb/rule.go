package importtodb

import (
	"fmt"
	"strconv"
	"unicode/utf8"

	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/excel/read"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure/buildtable"
	simpleDb "github.com/auho/go-simple-db/v2"
)

var _ Resourcer = (*RuleResource)(nil)

type RuleResource struct {
	Resource
	Titles
	KeywordIndex int // keyword 在 sheet 中的列，从 1 开始
	Rule         assistant.Ruler
}

func (rs *RuleResource) Prepare() error {
	return rs.Titles.prepare()
}

func (rs *RuleResource) getKeywordIndex() int {
	return rs.KeywordIndex - 1
}

func (rs *RuleResource) GetTable() buildtable.Tabler {
	return buildtable.NewRuleTable(rs.Rule)
}

func (rs *RuleResource) GetSheetData(excel *read.Excel) (read.SheetDataor, error) {
	sheetData, err := rs.readSheetData(excel, rs.buildSheetConfig())
	if err != nil {
		return nil, fmt.Errorf("readSheetData error; %w", err)
	}

	// drop duplicates TODO add if
	for i, title := range rs.titlesKey {
		if title == rs.Rule.KeywordName() {
			rs.ColumnDropDuplicates = append(rs.ColumnDropDuplicates, i)
			break
		}
	}

	// keyword len of string
	err = sheetData.HandlerRows(func(rows [][]string) ([][]string, error) {
		rs.titlesKey = append(rs.titlesKey, rs.Rule.KeywordLenName())

		for i, row := range rows {
			row = append(row, strconv.Itoa(utf8.RuneCountInString(row[rs.getKeywordIndex()])))
			rows[i] = row
		}

		return rows, nil
	})

	if err != nil {
		return nil, fmt.Errorf("keyWord len error; %w", err)
	}

	return sheetData, nil
}

func (rs *RuleResource) GetDB() *simpleDb.SimpleDB {
	return rs.Rule.GetDB()
}
