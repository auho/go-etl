package load

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/auho/go-etl/v3/insight/assistant"
	"github.com/auho/go-etl/v3/insight/assistant/excel/reader"
	"github.com/auho/go-etl/v3/insight/assistant/schema/create"
	simpledb "github.com/auho/go-simple-db/v3"
)

var _ Resource = (*Rule)(nil)

type Rule struct {
	baseResource
	Titles // column title of save to db

	Rule assistant.Rule
}

func (r *Rule) Prepare() error {
	return r.Titles.prepare()
}

func (r *Rule) Name() string {
	return r.Rule.Name()
}

func (r *Rule) Tabler() create.Tabler {
	return create.NewRuleTable(r.Rule)
}

func (r *Rule) SheetData(excel *reader.Excel) (reader.SheetDataReader, error) {
	sheetData, err := r.readSheetData(excel, r.buildSheetConfig())
	if err != nil {
		return nil, fmt.Errorf("readSheetData: %w", err)
	}

	keywordIndex := -1
	// drop duplicates TODO add if
	for i, title := range r.titlesKey {
		if title == r.Rule.KeywordName() {
			keywordIndex = r.titlesIndex[i]

			if !r.Rule.Config().AllowKeywordDuplicate() {
				r.columnDropDuplicates = append(r.columnDropDuplicates, i)
			}

			break
		}
	}

	//if keywordIndex < 0 {
	//	return nil, fmt.Errorf("keyword index error")
	//}

	// keyword len of string
	err = sheetData.HandleRows(func(rows [][]string) ([][]string, error) {
		r.titlesKey = append(r.titlesKey, r.Rule.KeywordLenName())

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

func (r *Rule) DB() *simpledb.SimpleDB {
	return r.Rule.DB()
}
