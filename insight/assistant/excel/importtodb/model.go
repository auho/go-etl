package importtodb

import (
	"fmt"
	"sort"

	"github.com/auho/go-etl/v2/insight/assistant/excel/read"
)

type Titles struct {
	Titles      map[int]string // map[sheet column index ]save to db 的 columns；index 从 1 开始
	titlesKey   []string
	titlesIndex []int // column index in sheet；从 0 开始
}

func (t *Titles) prepare() error {
	t.buildTitlesKey()

	return nil
}

func (t *Titles) GetTitlesKey() []string {
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
	for index := range t.Titles {
		t.titlesIndex = append(t.titlesIndex, index-1) // index 从 1 开始
	}

	sort.Slice(t.titlesIndex, func(i, j int) bool {
		return t.titlesIndex[i] < t.titlesIndex[j]
	})

	for _, index := range t.titlesIndex {
		t.titlesKey = append(t.titlesKey, t.Titles[index+1]) // index 从 1 开始
	}
}
