package importtodb

import (
	"github.com/auho/go-etl/v2/insight/assistant/excel/read"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure/buildtable"
	simpleDb "github.com/auho/go-simple-db/v2"
)

type Resourcer interface {
	GetIsRecreateTable() bool
	GetIsAppendData() bool
	GetIsShowSql() bool
	GetColumnDropDuplicates() []int
	GetDB() *simpleDb.SimpleDB

	Prepare() error
	GetName() string
	GetTable() buildtable.Tabler
	GetTitlesName() []string
	GetTitlesIndex() []int
	GetSheetData(*read.Excel) (read.SheetDataor, error)

	CommandExec(*tablestructure.Command)
	PostDo(Resourcer) error
}

type Resource struct {
	SheetName            string
	SheetIndex           int                           // sheet index，从 1 开始
	StartRow             int                           // 数据开始的行数，从 1 开始
	EndRow               int                           // 数据结束的行数，从 1 开始
	IsRecreateTable      bool                          // 是否 recreate table
	IsAppendData         bool                          // 是否 append data
	IsShowSql            bool                          // 是否显示 sql
	ColumnDropDuplicates []int                         // [column index] drop duplicates for column
	CommandFun           func(*tablestructure.Command) // recreate table 时执行的 func
	PostFun              func(Resourcer) error         // 导入后的执行的 func
}

func (s *Resource) buildSheetConfig() read.Config {
	return read.Config{
		SheetName:  s.SheetName,
		SheetIndex: s.SheetIndex,
		StartRow:   s.StartRow,
		EndRow:     s.EndRow,
	}
}

func (s *Resource) CommandExec(command *tablestructure.Command) {
	if s.CommandFun != nil {
		s.CommandFun(command)
	}
}

func (s *Resource) PostDo(resource Resourcer) error {
	if s.PostFun != nil {
		return s.PostFun(resource)
	}

	return nil
}

func (s *Resource) GetIsRecreateTable() bool {
	return s.IsRecreateTable
}

func (s *Resource) GetIsAppendData() bool {
	return s.IsAppendData
}

func (s *Resource) GetIsShowSql() bool {
	return s.IsShowSql
}

func (s *Resource) GetColumnDropDuplicates() []int {
	return s.ColumnDropDuplicates
}
