package task

import (
	"github.com/auho/go-etl/v3/job"
	"github.com/auho/go-etl/v3/job/extract"
	"github.com/auho/go-etl/v3/job/extract/tag"
	"github.com/auho/go-etl/v3/job/transform"
	"github.com/auho/go-etl/v3/job/transform/collect"
	simpledb "github.com/auho/go-simple-db/v3"
)

var _ job.Table = (*_jobSource)(nil)
var _ job.Table = (*_jobTarget)(nil)
var _ job.CleanResource = (*_cleanResource)(nil)

var _ruler extract.Rule

type _jobSource struct{}

func (_ _jobSource) IDName() string            { return "id" }
func (_ _jobSource) TableName() string         { return "source" }
func (_ _jobSource) GetDB() *simpledb.SimpleDB { return nil }

type _jobTarget struct{}

func (_ _jobTarget) IDName() string            { return "id" }
func (_ _jobTarget) TableName() string         { return "target" }
func (_ _jobTarget) GetDB() *simpledb.SimpleDB { return nil }

type _cleanResource struct{}

func (_ _cleanResource) Source() job.Table  { return &_jobSource{} }
func (_ _cleanResource) Data() job.Table    { return &_jobTarget{} }
func (_ _cleanResource) Deleted() job.Table { return &_jobTarget{} }

func ExampleNewClean() {
	_mode := transform.NewUpdate(collect.NewKeysAll([]string{"key1", "key2"}), tag.NewKey(_ruler), nil)
	_ = NewClean(
		&_cleanResource{},
		[]transform.UpdateOperator{_mode},
		WithCleanConfig(CleanConfig{
			NotTruncate:  false,
			AddExtraTags: false,
			BatchSize:    0,
			Concurrency:  0,
			Keys:         []string{"key3", "key4"},
		}),
	)
}

func ExampleNewInsert() {

	_mode := transform.NewInsert(collect.NewKeysAll([]string{"key1", "key2"}), tag.NewKey(_ruler), nil)
	_insMulti1 := transform.NewInsert(collect.NewKeysAll([]string{"key1", "key2"}), tag.NewKey(_ruler), nil)
	_insMulti2 := transform.NewInsert(collect.NewKeysAll([]string{"key1", "key2"}), tag.NewLabel(_ruler), nil)
	_modeMulti := transform.NewInsertStack(_insMulti1, _insMulti2)
	_insCross1 := transform.NewInsert(collect.NewKeysAll([]string{"key1", "key2"}), tag.NewMostKey(_ruler), nil)
	_insCross2 := transform.NewInsert(collect.NewKeysAll([]string{"key1", "key2"}), tag.NewMostText(_ruler), nil)
	_modeCross := transform.NewInsertCross(_insCross1, _insCross2)
	_insSpread1 := transform.NewInsert(collect.NewKeysAll([]string{"key1", "key2"}), tag.NewKey(_ruler), nil)
	_insSpread2 := transform.NewInsert(collect.NewKeysAll([]string{"key1", "key2"}), tag.NewKey(_ruler), nil)
	_modeSpread := transform.NewInsertSpread(_insSpread1, _insSpread2)

	_ = NewInsert(&_jobTarget{}, _mode, WithInsertConfig(InsertConfig{
		NotTruncate: false,
		BatchSize:   0,
		Concurrency: 0,
		ExtraKeys:   nil,
	}))

	_ = NewInsert(&_jobTarget{}, _modeMulti)
	_ = NewInsert(&_jobTarget{}, _modeCross)
	_ = NewInsert(&_jobTarget{}, _modeSpread)
}

func ExampleNewTransfer() {
	_mode := transform.NewTransfer(
		[]string{"key1", "key2"},
		map[string]string{"key1": "alias1"},
		map[string]any{"fixed1": "fixed value"},
	)

	_ = NewTransfer(&_jobTarget{}, _mode)
}

func ExampleNewUpdate() {
	_mode := transform.NewUpdate(collect.NewKeysAll([]string{"key1", "key2"}), tag.NewKey(_ruler), nil)

	_ = NewUpdate(&_jobSource{}, []transform.UpdateOperator{_mode})

	_ = NewUpdateTransfer(&_jobSource{}, &_jobTarget{}, []transform.UpdateOperator{_mode}, WithUpdateTransferConfig(UpdateTransferConfig{
		NotTruncate: false,
		BatchSize:   0,
		Concurrency: 0,
	}))
}
