package dbimporter

import (
	"github.com/auho/go-etl/v3/insight/assistant/excel/reader"
	"github.com/auho/go-etl/v3/insight/assistant/schema"
	"github.com/auho/go-etl/v3/insight/assistant/schema/buildtable"
	simpledb "github.com/auho/go-simple-db/v3"
)

type Resource interface {
	GetIsRecreateTable() bool
	GetIsAppendData() bool
	GetIsShowSql() bool
	GetBatchInsertSize() int
	GetColumnDropDuplicates() []int
	GetDB() *simpledb.SimpleDB

	Prepare() error
	GetName() string
	GetTable() buildtable.Tabler
	GetTitlesName() []string
	GetTitlesIndex() []int
	GetSheetData(*reader.Excel) (reader.SheetDataReader, error)

	ExecCommand(*schema.Command)
	AfterDo(Resource) error
}

type ResourceBase struct {
	SheetName            string
	SheetIndex           int                   // sheet index，从 1 开始
	StartRow             int                   // 数据开始的行数，从 1 开始
	EndRow               int                   // 数据结束的行数，从 1 开始
	BatchInsertSize      int                   // 数据批量插入 size
	IsRecreateTable      bool                  // true: recreate table; false: not recreate table;
	IsAppendData         bool                  // true: append data; false truncate table
	IsShowSql            bool                  // 是否显示 sql
	ColumnDropDuplicates []int                 // [column index] drop duplicates for column
	CommandFun           func(*schema.Command) // recreate table 时执行的 func
	PostFun              func(Resource) error  // 导入后的执行的 func
}

func (s *ResourceBase) buildSheetConfig() reader.Config {
	return reader.Config{
		SheetName:  s.SheetName,
		SheetIndex: s.SheetIndex,
		StartRow:   s.StartRow,
		EndRow:     s.EndRow,
	}
}

func (s *ResourceBase) ExecCommand(command *schema.Command) {
	if s.CommandFun != nil {
		s.CommandFun(command)
	}
}

func (s *ResourceBase) AfterDo(resource Resource) error {
	if s.PostFun != nil {
		return s.PostFun(resource)
	}

	return nil
}

func (s *ResourceBase) GetIsRecreateTable() bool {
	return s.IsRecreateTable
}

func (s *ResourceBase) GetIsAppendData() bool {
	return s.IsAppendData
}

func (s *ResourceBase) GetIsShowSql() bool {
	return s.IsShowSql
}

func (s *ResourceBase) GetBatchInsertSize() int {
	if s.BatchInsertSize <= 0 {
		s.BatchInsertSize = 2000
	}

	return s.BatchInsertSize
}

func (s *ResourceBase) GetColumnDropDuplicates() []int {
	return s.ColumnDropDuplicates
}
