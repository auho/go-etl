package load

import (
	"fmt"
	"sort"

	"github.com/auho/go-etl/v3/insight/assistant/excel/reader"
)

// Titles defines the mapping between sheet columns and database column names.
// TitlesWithIndex takes precedence over Titles when both are set.
type Titles struct {
	TitlesWithIndex map[int]string // map[sheet column index]db column name; index starts from 1; takes precedence over Titles
	Titles          []string       // db column names, mapped to consecutive sheet columns starting from the first
	titlesKey       []string       // resolved db column names in sheet column order
	titlesIndex     []int          // column indexes in sheet, starting from 0
}

func (t *Titles) prepare() error {
	t.buildTitlesKey()

	return t.check()
}

func (t *Titles) TitlesName() []string {
	return t.titlesKey
}

func (t *Titles) TitlesIndex() []int {
	return t.titlesIndex
}

func (t *Titles) readSheetData(excel *reader.Excel, sheetConfig reader.Config) (*reader.SheetDataNoTitle, error) {
	sheetConfig.ColsIndex = t.titlesIndex
	sheetData, err := reader.NewSheetDataNoTitle(excel, sheetConfig)
	if err != nil {
		return nil, fmt.Errorf("NewSheetDataNoTitle: %w", err)
	}

	err = sheetData.ReadData()
	if err != nil {
		return nil, fmt.Errorf("ReadData: %w", err)
	}

	return sheetData, nil
}

// buildTitlesKey resolves titlesKey and titlesIndex from Titles and TitlesWithIndex.
// TitlesWithIndex overrides Titles for the same column index.
func (t *Titles) buildTitlesKey() {
	_titlesWithIndex := make(map[int]string)

	// titles: consecutive columns starting from index 1
	for i, title := range t.Titles {
		_titlesWithIndex[i+1] = title // index starts from 1
	}

	// titles with index: override any existing entries
	for i, title := range t.TitlesWithIndex {
		_titlesWithIndex[i] = title
	}

	// collect column indexes (convert to 0-based)
	for index := range _titlesWithIndex {
		t.titlesIndex = append(t.titlesIndex, index-1) // convert to 0-based
	}

	// sort indexes ascending
	sort.Slice(t.titlesIndex, func(i, j int) bool {
		return t.titlesIndex[i] < t.titlesIndex[j]
	})

	// build keys in sorted index order
	for _, index := range t.titlesIndex {
		t.titlesKey = append(t.titlesKey, _titlesWithIndex[index+1]) // convert back to 1-based
	}
}

func (t *Titles) check() error {
	if len(t.titlesKey) <= 0 {
		return fmt.Errorf("titles key does not exist")
	}

	for i, index := range t.titlesIndex {
		if index < 0 {
			return fmt.Errorf("title[%s] index[%d] is invalid", t.titlesKey[i], i)
		}
	}

	return nil
}
