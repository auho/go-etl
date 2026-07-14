package reader

// SheetDataReader defines the contract for reading and transforming sheet data from an Excel file.
// Implementations include SheetDataNoTitle (no header row) and SheetDataWithTitle (with header row).
type SheetDataReader interface {
	// ReadData reads rows from the Excel sheet into memory.
	ReadData() error
	// TransformRows applies fn to the current rows, replacing them with the result.
	TransformRows(fn func(rows [][]string) ([][]string, error)) error
	// GetRows returns the current rows as [][]string.
	GetRows() [][]string
	// GetRowsAsAny returns the current rows as [][]any for use with generic insert APIs.
	GetRowsAsAny() [][]any
}

// Config specifies which sheet and row range to read from an Excel file.
type Config struct {
	SheetName     string
	SheetIndex    int   // sheet index, starting from 1
	StartRow      int   // data start row, starting from 1
	EndRow        int   // data end row, starting from 1; 0 means no limit
	ColumnIndexes []int // column indexes to select, starting from 0; empty selects all columns
}
