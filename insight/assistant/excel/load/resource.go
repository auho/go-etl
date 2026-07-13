package load

import (
	"github.com/auho/go-etl/v3/insight/assistant/excel/reader"
	"github.com/auho/go-etl/v3/insight/assistant/schema"
	"github.com/auho/go-etl/v3/insight/assistant/schema/create"
	simpledb "github.com/auho/go-simple-db/v3"
)

// Resource defines the contract for loading Excel sheet data into a database table.
// Each implementation wraps an entity (Raw, Rows, or Rule) and provides
// sheet reading, table creation, and data insertion behavior.
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

// baseResource holds the shared configuration and hooks embedded by all
// concrete Resource implementations (Raw, Rows, Rule).
type baseResource struct {
	SheetName            string
	SheetIndex           int                   // sheet index, starting from 1
	StartRow             int                   // data start row, starting from 1
	EndRow               int                   // data end row, starting from 1
	BatchInsertSize      int                   // batch insert size; defaults to 2000 if <= 0
	IsRecreateTable      bool                  // true: recreate table; false: not recreate table
	IsAppendData         bool                  // true: append data; false: truncate table
	IsShowSql            bool                  // whether to print SQL
	columnDropDuplicates []int                 // column indexes to drop duplicates on
	CommandFun           func(*schema.Command) // function executed when recreating table
	AfterFun             func(Resource) error  // function executed after import
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
