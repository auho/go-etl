package runner

import (
	"github.com/auho/go-etl/v3/internal/testutil/mysql"
	"github.com/auho/go-etl/v3/task"
	simpledb "github.com/auho/go-simple-db/v3"
)

var _ task.Table = (*sourceTest)(nil)
var _ task.Table = (*targetTagATest)(nil)
var _ task.Table = (*targetTransferTest)(nil)
var _ task.Table = (*targetUpdateTransferTest)(nil)
var _ task.Table = (*targetCleanDataTest)(nil)
var _ task.Table = (*targetCleanDeletedTest)(nil)
var _ task.CleanResource = (*targetCleanTest)(nil)

// sourceTest
type sourceTest struct {
	tableName string
}

func (s *sourceTest) IDName() string {
	return _pkName
}

func (s *sourceTest) TableName() string {
	return s.tableName
}

func (s *sourceTest) GetDB() *simpledb.SimpleDB {
	simpleDB, _ := mysql.NewDB()
	return simpleDB
}

// targetTest
type targetTest struct{}

func (t targetTest) IDName() string {
	return "id"
}

func (t targetTest) GetDB() *simpledb.SimpleDB {
	simpleDB, _ := mysql.NewDB()
	return simpleDB
}

// targetTagATest
type targetTagATest struct {
	targetTest
}

func (t targetTagATest) TableName() string {
	return _tagATable
}

// targetTransferTest
type targetTransferTest struct {
	targetTest
}

func (t targetTransferTest) TableName() string {
	return _transferTable
}

// targetUpdateTransferTest
type targetUpdateTransferTest struct {
	targetTest
}

func (t targetUpdateTransferTest) TableName() string {
	return _updateAndTransferTable
}

// targetCleanDataTest
type targetCleanDataTest struct {
	targetTest
}

func (t targetCleanDataTest) TableName() string {
	return _cleanDataTable
}

// targetCleanDeletedTest
type targetCleanDeletedTest struct {
	targetTest
}

func (t targetCleanDeletedTest) TableName() string {
	return _deletedDataTable
}

// targetCleanTest
type targetCleanTest struct {
	targetTest
	source task.Table
}

func (t *targetCleanTest) Source() task.Table {
	return t.source
}

func (t *targetCleanTest) Data() task.Table {
	return &targetCleanDataTest{}
}

func (t *targetCleanTest) Deleted() task.Table {
	return &targetCleanDeletedTest{}
}
