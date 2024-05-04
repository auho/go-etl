package read

type SheetDataor interface {
	ReadData() error
	HandlerRows(fn func(rows [][]string) ([][]string, error)) error
	GetRows() [][]string
	GetRowsWithAny() [][]any
}

type Config struct {
	SheetName  string
	SheetIndex int   // sheet index，从 1 开始
	StartRow   int   // 数据开始的行数，从 1 开始
	EndRow     int   // 数据结束的行数，从 1 开始
	ColsIndex  []int // columns index 从 0 开始
}
