package read

var _ SheetDataor = (*SheetDataWithTitle)(nil)

type SheetDataWithTitle struct {
	SheetData
	titles []string          // []name in sheet
	alias  map[string]string // map[name in sheet ] name of logic
}

func NewSheetDataWithTitle(xlsxPath, sheetName string, startRow int, alias map[string]string) (*SheetDataWithTitle, error) {
	sd := &SheetDataWithTitle{}
	sd.xlsxPath = xlsxPath
	sd.sheetName = sheetName
	sd.startRow = startRow
	sd.alias = alias

	return sd, nil
}

func (sd *SheetDataWithTitle) GetTitles() []string {
	return sd.titles
}

func (sd *SheetDataWithTitle) GetAlias() map[string]string {
	return sd.alias
}

func (sd *SheetDataWithTitle) GetTitlesWithAlias() []string {
	var titles []string
	for _, title := range sd.titles {
		if alias, ok := sd.alias[title]; ok {
			title = alias
		}

		titles = append(titles, title)
	}

	return titles
}

func (sd *SheetDataWithTitle) ReadData() error {
	err := sd.readFromSheet()
	if err != nil {
		return err
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

	return nil
}
