package load

import (
	"github.com/auho/go-etl/v3/insight/assistant/excel/reader"
	"github.com/auho/go-etl/v3/insight/assistant/schema"
	"github.com/auho/go-etl/v3/insight/assistant/schema/create"
	simpledb "github.com/auho/go-simple-db/v3"
)

type Resource interface {
	GetIsRecreateTable() bool
	GetIsAppendData() bool
	GetIsShowSql() bool
	GetBatchInsertSize() int
	GetColumnDropDuplicates() []int
	DB() *simpledb.SimpleDB

	Prepare() error
	Name() string
	Tabler() create.Tabler
	TitlesName() []string
	TitlesIndex() []int
	SheetData(*reader.Excel) (reader.SheetDataReader, error)

	ExecCommand(*schema.Command)
	AfterDo(Resource) error
}

type baseResource struct {
	SheetName            string
	SheetIndex           int                   // sheet index，从 1 开始
	StartRow             int                   // 数据开始的行数，从 1 开始
	EndRow               int                   // 数据结束的行数，从 1 开始
	BatchInsertSize      int                   // 数据批量插入 size
	IsRecreateTable      bool                  // true: recreate table; false: not recreate table;
	IsAppendData         bool                  // true: append data; false truncate table
	IsShowSql            bool                  // 是否显示 sql
	columnDropDuplicates []int                 // [column index] drop duplicates for column
	CommandFun           func(*schema.Command) // recreate table 时执行的 func
	AfterFun             func(Resource) error  // 导入后的执行的 func
}

func (b *baseResource) buildSheetConfig() reader.Config {
	return reader.Config{
		SheetName:  b.SheetName,
		SheetIndex: b.SheetIndex,
		StartRow:   b.StartRow,
		EndRow:     b.EndRow,
	}
}

func (b *baseResource) ExecCommand(command *schema.Command) {
	if b.CommandFun != nil {
		b.CommandFun(command)
	}
}

func (b *baseResource) AfterDo(resource Resource) error {
	if b.AfterFun != nil {
		return b.AfterFun(resource)
	}

	return nil
}

func (b *baseResource) GetIsRecreateTable() bool {
	return b.IsRecreateTable
}

func (b *baseResource) GetIsAppendData() bool {
	return b.IsAppendData
}

func (b *baseResource) GetIsShowSql() bool {
	return b.IsShowSql
}

func (b *baseResource) GetBatchInsertSize() int {
	if b.BatchInsertSize <= 0 {
		b.BatchInsertSize = 2000
	}

	return b.BatchInsertSize
}

func (b *baseResource) GetColumnDropDuplicates() []int {
	return b.columnDropDuplicates
}
