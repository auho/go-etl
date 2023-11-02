package exporttodb

import (
	"github.com/auho/go-etl/v2/insight/model/excel/read"
	"github.com/auho/go-etl/v2/insight/model/table"
)

type sourceor interface {
	GetIsRecreateTable() bool
	GetIsAppendData() bool
	GetIsShowSql() bool
	GetColumnDropDuplicates() []int

	GetTable() table.Tabler
	GetTitles() []string
	GetSheetData() (read.SheetDataor, error)
}

type Source struct {
	XlsxPath             string
	SheetName            string
	StartRow             int   // 数据开始的行数， 从 0 开始
	IsRecreateTable      bool  // 是否 recreate table
	IsAppendData         bool  // 是否 append data
	IsShowSql            bool  // 是否显示 sql
	ColumnDropDuplicates []int // drop duplicates for column
}

func (s *Source) GetIsRecreateTable() bool {
	return s.IsRecreateTable
}

func (s *Source) GetIsAppendData() bool {
	return s.IsAppendData
}

func (s *Source) GetIsShowSql() bool {
	return s.IsShowSql
}

func (s *Source) GetColumnDropDuplicates() []int {
	return s.ColumnDropDuplicates
}
