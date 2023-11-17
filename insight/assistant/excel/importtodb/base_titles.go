package importtodb

import (
	"fmt"
	"sort"

	"github.com/auho/go-etl/v2/insight/assistant/excel/read"
)

// Titles
// column title of save to db
type Titles struct {
	TitlesWithIndex map[int]string // map[sheet column index ]save to db 的 columns；index 从 1 开始
	Titles          []string       // []save to db 的 columns; 从第一个 column 开始，连续不间断；此选择优先
	titlesKey       []string
	titlesIndex     []int // column index in sheet；从 0 开始
}

func (t *Titles) prepare() error {
	t.buildTitlesKey()

	return t.check()
}

func (t *Titles) GetTitlesName() []string {
	return t.titlesKey
}

func (t *Titles) GetTitlesIndex() []int {
	return t.titlesIndex
}

func (t *Titles) readSheetData(excel *read.Excel, sheetConfig read.Config) (*read.SheetDataNoTitle, error) {
	sheetConfig.ColsIndex = t.titlesIndex
	sheetData, err := read.NewSheetDataNoTitle(excel, sheetConfig)
	if err != nil {
		return nil, fmt.Errorf("NewSheetDataNoTitle error; %w", err)
	}

	err = sheetData.ReadData()
	if err != nil {
		return nil, fmt.Errorf("ReadData error; %w", err)
	}

	return sheetData, nil
}

// buildTitlesKey
// []string titles
// []int columns index of title
func (t *Titles) buildTitlesKey() {
	if len(t.Titles) > 0 {
		t.TitlesWithIndex = make(map[int]string, len(t.Titles))
		for i, title := range t.Titles {
			t.TitlesWithIndex[i+1] = title // 因传入的 index 从 1 开始
		}
	}

	for index := range t.TitlesWithIndex {
		t.titlesIndex = append(t.titlesIndex, index-1) // 因传入的 index 从 1 开始
	}

	sort.Slice(t.titlesIndex, func(i, j int) bool {
		return t.titlesIndex[i] < t.titlesIndex[j]
	})

	for _, index := range t.titlesIndex {
		t.titlesKey = append(t.titlesKey, t.TitlesWithIndex[index+1]) // 因传入的 index 从 1 开始
	}
}

func (t *Titles) check() error {
	if len(t.titlesKey) <= 0 {
		return fmt.Errorf("titles key no exists")
	}

	for i, index := range t.titlesIndex {
		if index < 0 {
			return fmt.Errorf("title[%s] index[%d] is error", t.titlesKey[i], i)
		}
	}

	return nil
}
