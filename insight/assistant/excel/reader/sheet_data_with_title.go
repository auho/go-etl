package reader

import (
	"fmt"
)

var _ SheetDataReader = (*SheetDataWithTitle)(nil)

// SheetDataWithTitle reads sheet data with a title row from Excel.
// The first row is treated as column titles; remaining rows are data.
// An optional alias map translates sheet column names to logical column names.
type SheetDataWithTitle struct {
	sheetData
	titles []string          // column names from the sheet's first row
	alias  map[string]string // map[sheet column name]logical column name
}

// NewSheetDataWithTitleWithPath creates a SheetDataWithTitle by opening the xlsx file at xlsxPath.
func NewSheetDataWithTitleWithPath(xlsxPath string, config Config, alias map[string]string) (*SheetDataWithTitle, error) {
	excel, err := NewExcel(xlsxPath)
	if err != nil {
		return nil, err
	}

	return NewSheetDataWithTitle(excel, config, alias)
}

// NewSheetDataWithTitle creates a SheetDataWithTitle from an existing Excel.
func NewSheetDataWithTitle(excel *Excel, config Config, alias map[string]string) (*SheetDataWithTitle, error) {
	sd := &SheetDataWithTitle{}
	sd.excel = excel
	sd.config = config
	sd.alias = alias

	return sd, nil
}

// GetTitles returns the resolved column names (after alias translation).
func (sd *SheetDataWithTitle) GetTitles() []string {
	return sd.titles
}

// GetAlias returns the alias map from sheet column names to logical column names.
func (sd *SheetDataWithTitle) GetAlias() map[string]string {
	return sd.alias
}

// applyAliasToTitles replaces each title with its alias if one is defined.
func (sd *SheetDataWithTitle) applyAliasToTitles() {
	var titles []string
	for _, title := range sd.titles {
		if alias, ok := sd.alias[title]; ok {
			title = alias
		}

		titles = append(titles, title)
	}

	sd.titles = titles
}

// ReadData reads rows from the configured sheet, extracts the title row,
// trims data rows to the title length, and applies alias translation.
func (sd *SheetDataWithTitle) ReadData() error {
	err := sd.readSheet()
	if err != nil {
		return fmt.Errorf("readSheet: %w", err)
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

	sd.applyAliasToTitles()

	return nil
}
