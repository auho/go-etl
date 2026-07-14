package reader

var _ SheetDataReader = (*SheetDataNoTitle)(nil)

// SheetDataNoTitle reads sheet data without a title row from Excel.
// All rows are treated as data rows.
type SheetDataNoTitle struct {
	sheetData
}

// NewSheetDataNoTitleWithPath creates a SheetDataNoTitle by opening the xlsx file at xlsxPath.
func NewSheetDataNoTitleWithPath(xlsxPath string, config Config) (*SheetDataNoTitle, error) {
	excel, err := NewExcel(xlsxPath)
	if err != nil {
		return nil, err
	}

	return NewSheetDataNoTitle(excel, config)
}

// NewSheetDataNoTitle creates a SheetDataNoTitle from an existing Excel.
func NewSheetDataNoTitle(excel *Excel, config Config) (*SheetDataNoTitle, error) {
	sd := &SheetDataNoTitle{}
	sd.excel = excel
	sd.config = config

	return sd, nil
}

// ReadData reads all rows from the configured sheet into memory.
func (sd *SheetDataNoTitle) ReadData() error {
	return sd.readSheet()
}
