package dataset

// Dataset
// data set
// 处理前
type Dataset struct {
	Name   string   // dataset name
	Titles []string // dataset item data title
	Sets   []Set
}

// Data
// 最后处理后的
type Data struct {
	Names []string           // 保存 name 的顺序
	Rows  map[string][][]any // map[name][][]any
}

func (d *Data) add(name string, rows [][]any) {
	if len(d.Rows) <= 0 {
		d.Rows = make(map[string][][]any)
	}

	d.Names = append(d.Names, name)
	d.Rows[name] = rows
}
