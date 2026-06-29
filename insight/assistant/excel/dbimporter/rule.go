package dbimporter

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/auho/go-etl/v3/insight/assistant"
	"github.com/auho/go-etl/v3/insight/assistant/excel/reader"
	"github.com/auho/go-etl/v3/insight/assistant/schema/buildtable"
	simpledb "github.com/auho/go-simple-db/v3"
)

var _ Resource = (*RuleResource)(nil)

type RuleResource struct {
	ResourceBase
	Titles // column title of save to db
	Rule   assistant.Rule
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

func (rs *RuleResource) GetSheetData(excel *reader.Excel) (reader.SheetDataReader, error) {
	sheetData, err := rs.readSheetData(excel, rs.buildSheetConfig())
	if err != nil {
		return nil, fmt.Errorf("readSheetData: %w", err)
	}

	keywordIndex := -1
	// drop duplicates TODO add if
	for i, title := range rs.titlesKey {
		if title == rs.Rule.KeywordName() {
			keywordIndex = rs.titlesIndex[i]

			if !rs.Rule.Config().AllowKeywordDuplicate() {
				rs.ColumnDropDuplicates = append(rs.ColumnDropDuplicates, i)
			}

			break
		}
	}

	//if keywordIndex < 0 {
	//	return nil, fmt.Errorf("keyword index error")
	//}

	// keyword len of string
	err = sheetData.HandleRows(func(rows [][]string) ([][]string, error) {
		rs.titlesKey = append(rs.titlesKey, rs.Rule.KeywordLenName())

		var _newRows [][]string
		for _, row := range rows {
			keywordLen := "0"

			if keywordIndex > -1 {
				_ky := row[keywordIndex]
				_ky = strings.TrimSpace(_ky)
				if _ky == "" {
					continue
				}

				row[keywordIndex] = _ky
				keywordLen = strconv.Itoa(utf8.RuneCountInString(_ky))
			}

			row = append(row, keywordLen)
			_newRows = append(_newRows, row)
		}

		return _newRows, nil
	})

	if err != nil {
		return nil, fmt.Errorf("HandleRows: %w", err)
	}

	return sheetData, nil
}

func (rs *RuleResource) GetDB() *simpledb.SimpleDB {
	return rs.Rule.GetDB()
}
