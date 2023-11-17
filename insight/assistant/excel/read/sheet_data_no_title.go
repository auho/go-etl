package read

var _ SheetDataor = (*SheetDataNoTitle)(nil)

type SheetDataNoTitle struct {
	sheetData
}

func NewSheetDataNoTitleWithPath(xlsxPath string, config Config) (*SheetDataNoTitle, error) {
	excel, err := NewExcel(xlsxPath)
	if err != nil {
		return nil, err
	}

	return NewSheetDataNoTitle(excel, config)
}

func NewSheetDataNoTitle(excel *Excel, config Config) (*SheetDataNoTitle, error) {
	sd := &SheetDataNoTitle{}
	sd.excel = excel
	sd.config = config

	return sd, nil
}
func (sd *SheetDataNoTitle) ReadData() error {
	return sd.readSheet()
}
