package task

import (
	"github.com/auho/go-etl/v2/job"
	simpleDb "github.com/auho/go-simple-db/v2"
)

var _ job.Source = (*sourceTest)(nil)
var _ job.Target = (*targetTagATest)(nil)
var _ job.Target = (*targetTagA1Test)(nil)
var _ job.Target = (*targetTagA2Test)(nil)
var _ job.Target = (*targetTransferTest)(nil)
var _ job.Target = (*targetUpdateTransferTest)(nil)
var _ job.Target = (*targetCleanDataTest)(nil)
var _ job.Target = (*targetCleanDeletedTest)(nil)
var _ job.CleanResource = (*targetCleanTest)(nil)

// sourceTest
type sourceTest struct{}

func (s sourceTest) GetIdName() string {
	return _pkName
}

func (s sourceTest) TableName() string {
	return _dataTable
}

func (s sourceTest) GetDB() *simpleDb.SimpleDB {
	return _db
}

// targetTest
type targetTest struct{}

func (t targetTest) GetIdName() string {
	return "id"
}

func (t targetTest) GetDB() *simpleDb.SimpleDB {
	return _db
}

// targetTagATest
type targetTagATest struct {
	targetTest
}

func (t targetTagATest) TableName() string {
	return _tagATable
}

// targetTagA1Test
type targetTagA1Test struct {
	targetTagATest
}

func (t targetTagA1Test) TableName() string {
	return _tagATable + "1"
}

// targetTagA2Test
type targetTagA2Test struct {
	targetTagATest
}

func (t targetTagA2Test) TableName() string {
	return _tagATable + "2"
}

// targetTransferTest
type targetTransferTest struct {
	targetTest
}

func (t targetTransferTest) TableName() string {
	return _transferTable
}

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
}

func (t targetCleanTest) SourceTarget() job.Target {
	return &sourceTest{}
}

func (t targetCleanTest) DataTarget() job.Target {
	return &targetCleanDataTest{}
}

func (t targetCleanTest) DeletedTarget() job.Target {
	return &targetCleanDeletedTest{}
}
