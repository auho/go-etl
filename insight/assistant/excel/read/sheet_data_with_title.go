package read

import (
	"fmt"
)

var _ SheetDataor = (*SheetDataWithTitle)(nil)

type SheetDataWithTitle struct {
	sheetData
	titles []string          // []name in sheet
	alias  map[string]string // map[name in sheet ] name of logic
}

func NewSheetDataWithTitleWithPath(xlsxPath string, config Config, alias map[string]string) (*SheetDataWithTitle, error) {
	excel, err := NewExcel(xlsxPath)
	if err != nil {
		return nil, err
	}

	return NewSheetDataWithTitle(excel, config, alias)
}

func NewSheetDataWithTitle(excel *Excel, config Config, alias map[string]string) (*SheetDataWithTitle, error) {
	sd := &SheetDataWithTitle{}
	sd.excel = excel
	sd.config = config
	sd.alias = alias

	return sd, nil
}

func (sd *SheetDataWithTitle) GetTitles() []string {
	return sd.titles
}

func (sd *SheetDataWithTitle) GetAlias() map[string]string {
	return sd.alias
}

func (sd *SheetDataWithTitle) genTitlesWithAlias() {
	var titles []string
	for _, title := range sd.titles {
		if alias, ok := sd.alias[title]; ok {
			title = alias
		}

		titles = append(titles, title)
	}

	sd.titles = titles
}

func (sd *SheetDataWithTitle) ReadData() error {
	err := sd.readSheet()
	if err != nil {
		return fmt.Errorf("readSheet error; %w", err)
	}

	sd.titles = sd.rows[0]
	titleLen := len(sd.titles)

	sd.rows = sd.rows[1:]

	for i, row := range sd.rows {
		if len(row) >= titleLen {
			row = row[0:titleLen]
		}

		sd.rows[i] = row
	}

	sd.genTitlesWithAlias()

	return nil
}
