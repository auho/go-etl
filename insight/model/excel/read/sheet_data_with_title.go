package read

type SheetDataWithTitle struct {
	SheetData
	titles []string          // []name in sheet
	alias  map[string]string // map[name in sheet ] name of logic
}

func NewSheetDataWithTitle(xlsxPath, sheetName string, alias map[string]string) (*SheetDataWithTitle, error) {
	sd := &SheetDataWithTitle{}
	sd.xlsxPath = xlsxPath
	sd.sheetName = sheetName
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

func (sd *SheetDataWithTitle) ReadData(rows [][]string) error {
	err := sd.readFromSheet()
	if err != nil {
		return err
	}

	sd.titles = rows[0]
	titleLen := len(sd.titles)

	rows = rows[1:]

	for _, row := range rows {
		if len(row) >= titleLen {
			row = row[0:titleLen]
		}

		sd.rows = append(sd.rows, row)
	}

	return nil
}
