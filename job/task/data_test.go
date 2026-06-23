package task

import (
	"github.com/auho/go-etl/v2/job"
	simpledb "github.com/auho/go-simple-db/v2"
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

func (s sourceTest) GetIDName() string {
	return _pkName
}

func (s sourceTest) TableName() string {
	return _dataTable
}

func (s sourceTest) GetDB() *simpledb.SimpleDB {
	return _simpleDB
}

// targetTest
type targetTest struct{}

func (t targetTest) GetIDName() string {
	return "id"
}

func (t targetTest) GetDB() *simpledb.SimpleDB {
	return _simpleDB
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

func (t targetCleanTest) Source() job.Target {
	return &sourceTest{}
}

func (t targetCleanTest) Data() job.Target {
	return &targetCleanDataTest{}
}

func (t targetCleanTest) Deleted() job.Target {
	return &targetCleanDeletedTest{}
}
