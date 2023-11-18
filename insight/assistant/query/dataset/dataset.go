package dataset

// Dataset
// data set
// 处理前的
type Dataset struct {
	Name   string   // dataset name
	Keys   []string // item key name
	Titles []string // dataset item data title
	Sets   []Set
}

// Data
// 处理后的
type Data struct {
	Names      []string           // data name 保存 name 的顺序
	Rows       map[string][][]any // map[name]rows
	RowsAmount map[string]int     // rows num (不包含 title)
	Amount     int                // rows total num (不包含 title)
}

func (d *Data) addRows() {

}

func (d *Data) addRowsWithTitles(name string, titles []any, rows [][]any) {
	if len(d.Rows) <= 0 {
		d.Rows = make(map[string][][]any)
	}

	if len(d.RowsAmount) <= 0 {
		d.RowsAmount = make(map[string]int)
	}

	d.Names = append(d.Names, name)
	d.RowsAmount[name] = len(rows)
	d.Amount += d.RowsAmount[name]
	d.Rows[name] = append([][]any{titles}, rows...)
}
