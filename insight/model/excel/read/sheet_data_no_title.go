package read

type SheetDataNoTitle struct {
	SheetData
}

func NewSheetDataNoTitle(xlsxPath, sheetName string) (*SheetDataNoTitle, error) {
	sd := &SheetDataNoTitle{}
	sd.xlsxPath = xlsxPath
	sd.sheetName = sheetName

	return sd, nil
}

func (sd *SheetDataNoTitle) ReadData() error {
	return sd.readFromSheet()
}
