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
var _ job.Target = (*targetCleanTest)(nil)

// sourceTest
type sourceTest struct {
}

func (s sourceTest) GetIdName() string {
	return _pkName
}

func (s sourceTest) TableName() string {
	return _dataTable
}

func (s sourceTest) GetDB() *simpleDb.SimpleDB {
	return _db
}

// targetTagATest
type targetTagATest struct {
}

func (t targetTagATest) GetIdName() string {
	return "id"
}

func (t targetTagATest) TableName() string {
	return _tagATable
}

func (t targetTagATest) GetDB() *simpleDb.SimpleDB {
	return _db
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
}

func (t targetTransferTest) GetIdName() string {
	return "id"
}

func (t targetTransferTest) TableName() string {
	return _transferTable
}

func (t targetTransferTest) GetDB() *simpleDb.SimpleDB {
	return _db
}

type targetUpdateTransferTest struct {
}

func (t targetUpdateTransferTest) GetIdName() string {
	return "id"
}

func (t targetUpdateTransferTest) TableName() string {
	return _updateAndTransferTable
}

func (t targetUpdateTransferTest) GetDB() *simpleDb.SimpleDB {
	return _db
}

// targetCleanTest
type targetCleanTest struct {
}

func (t targetCleanTest) GetIdName() string {
	return "id"
}

func (t targetCleanTest) TableName() string {
	return _cleanTable
}

func (t targetCleanTest) GetDB() *simpleDb.SimpleDB {
	return _db
}
