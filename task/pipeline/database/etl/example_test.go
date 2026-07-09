package etl

import (
	"github.com/auho/go-etl/v3/task/extract"
	"github.com/auho/go-etl/v3/task/extract/tag"
	"github.com/auho/go-etl/v3/task/transform"
	"github.com/auho/go-etl/v3/task/transform/collector"
	simpledb "github.com/auho/go-simple-db/v3"
)

var _ Table = (*jobSource)(nil)
var _ Table = (*jobTarget)(nil)
var _ CleanResource = (*cleanResource)(nil)

var ruler extract.Rule

type jobSource struct{}

func (_ jobSource) IDName() string         { return "id" }
func (_ jobSource) TableName() string      { return "source" }
func (_ jobSource) DB() *simpledb.SimpleDB { return nil }

type jobTarget struct{}

func (_ jobTarget) IDName() string         { return "id" }
func (_ jobTarget) TableName() string      { return "target" }
func (_ jobTarget) DB() *simpledb.SimpleDB { return nil }

type cleanResource struct{}

func (_ cleanResource) Source() Table  { return &jobSource{} }
func (_ cleanResource) Data() Table    { return &jobTarget{} }
func (_ cleanResource) Deleted() Table { return &jobTarget{} }

func ExampleNewClean() {
	operator := transform.NewUpdate(collector.NewKeysAll([]string{"key1", "key2"}, tag.NewKey(ruler)), nil)
	_ = NewClean(
		&cleanResource{},
		[]transform.UpdateOperator{operator},
		WithCleanConfig(CleanConfig{
			SkipTruncate: false,
			ExtraTags:    false,
			Keys:         []string{"key3", "key4"},
		}),
	)
}

func ExampleNewInsert() {
	operator := transform.NewInsert(collector.NewKeysAll([]string{"key1", "key2"}, tag.NewKey(ruler)), nil)
	insMulti1 := transform.NewInsert(collector.NewKeysAll([]string{"key1", "key2"}, tag.NewKey(ruler)), nil)
	insMulti2 := transform.NewInsert(collector.NewKeysAll([]string{"key1", "key2"}, tag.NewKey(ruler)), nil)
	operatorMulti := transform.NewInsertStack(insMulti1, insMulti2)
	insCross1 := transform.NewInsert(collector.NewKeysAll([]string{"key1", "key2"}, tag.NewMostKey(ruler)), nil)
	insCross2 := transform.NewInsert(collector.NewKeysAll([]string{"key1", "key2"}, tag.NewMostText(ruler)), nil)
	operatorCross := transform.NewInsertCross(insCross1, insCross2)
	insSpread1 := transform.NewInsert(collector.NewKeysAll([]string{"key1", "key2"}, tag.NewKey(ruler)), nil)
	insSpread2 := transform.NewInsert(collector.NewKeysAll([]string{"key1", "key2"}, tag.NewKey(ruler)), nil)
	operatorSpread := transform.NewInsertSpread(insSpread1, insSpread2)

	_ = NewInsert(&jobTarget{}, operator, WithInsertConfig(InsertConfig{
		SkipTruncate: false,
		ExtraKeys:    nil,
	}))
	_ = NewInsert(&jobTarget{}, operatorMulti)
	_ = NewInsert(&jobTarget{}, operatorCross)
	_ = NewInsert(&jobTarget{}, operatorSpread)
}

func ExampleNewTransfer() {
	operator := transform.NewTransfer(
		[]string{"key1", "key2"},
		map[string]string{"key1": "alias1"},
		map[string]any{"fixed1": "fixed value"},
	)

	_ = NewTransfer(&jobTarget{}, operator)
}

func ExampleNewUpdate() {
	operator := transform.NewUpdate(collector.NewKeysAll([]string{"key1", "key2"}, tag.NewKey(ruler)), nil)

	_ = NewUpdate(&jobSource{}, []transform.UpdateOperator{operator})

	_ = NewUpdateTransfer(&jobSource{}, &jobTarget{}, []transform.UpdateOperator{operator}, WithUpdateTransferConfig(UpdateTransferConfig{
		SkipTruncate: false,
	}))
}
