package read

var _ SheetDataor = (*SheetDataNoTitle)(nil)

type SheetDataNoTitle struct {
	SheetData
}

func NewSheetDataNoTitle(xlsxPath, sheetName string, startRow int) (*SheetDataNoTitle, error) {
	sd := &SheetDataNoTitle{}
	sd.xlsxPath = xlsxPath
	sd.sheetName = sheetName
	sd.startRow = startRow

	return sd, nil
}

func (sd *SheetDataNoTitle) ReadData() error {
	return sd.readFromSheet()
}
