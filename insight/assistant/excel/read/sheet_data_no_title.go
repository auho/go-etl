package read

var _ SheetDataor = (*SheetDataNoTitle)(nil)

type SheetDataNoTitle struct {
	sheetData
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
