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
	Titles // column title of save to db
	Rule   assistant.Ruler
}

func (rs *RuleResource) Prepare() error {
	return rs.Titles.prepare()
}

func (rs *RuleResource) GetName() string {
	return rs.Rule.GetName()
}

func (rs *RuleResource) GetTable() buildtable.Tabler {
	return buildtable.NewRuleTable(rs.Rule)
}

func (rs *RuleResource) GetSheetData(excel *read.Excel) (read.SheetDataor, error) {
	sheetData, err := rs.readSheetData(excel, rs.buildSheetConfig())
	if err != nil {
		return nil, fmt.Errorf("readSheetData error; %w", err)
	}

	keywordIndex := -1
	// drop duplicates TODO add if
	for i, title := range rs.titlesKey {
		if title == rs.Rule.KeywordName() {
			keywordIndex = rs.titlesIndex[i]
			rs.ColumnDropDuplicates = append(rs.ColumnDropDuplicates, i)
			break
		}
	}

	if keywordIndex < 0 {
		return nil, fmt.Errorf("keyword index error")
	}

	// keyword len of string
	err = sheetData.HandlerRows(func(rows [][]string) ([][]string, error) {
		rs.titlesKey = append(rs.titlesKey, rs.Rule.KeywordLenName())

		for i, row := range rows {
			row = append(row, strconv.Itoa(utf8.RuneCountInString(row[keywordIndex])))
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
