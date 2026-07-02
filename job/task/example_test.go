package task

import (
	"github.com/auho/go-etl/v3/job"
	"github.com/auho/go-etl/v3/job/extract"
	"github.com/auho/go-etl/v3/job/extract/tag"
	"github.com/auho/go-etl/v3/job/transform"
	"github.com/auho/go-etl/v3/job/transform/collector"
	simpledb "github.com/auho/go-simple-db/v3"
)

var _ job.Table = (*jobSource)(nil)
var _ job.Table = (*jobTarget)(nil)
var _ job.CleanResource = (*cleanResource)(nil)

var ruler extract.Rule

type jobSource struct{}

func (_ jobSource) IDName() string            { return "id" }
func (_ jobSource) TableName() string         { return "source" }
func (_ jobSource) GetDB() *simpledb.SimpleDB { return nil }

type jobTarget struct{}

func (_ jobTarget) IDName() string            { return "id" }
func (_ jobTarget) TableName() string         { return "target" }
func (_ jobTarget) GetDB() *simpledb.SimpleDB { return nil }

type cleanResource struct{}

func (_ cleanResource) Source() job.Table  { return &jobSource{} }
func (_ cleanResource) Data() job.Table    { return &jobTarget{} }
func (_ cleanResource) Deleted() job.Table { return &jobTarget{} }

func ExampleNewClean() {
	mode := transform.NewUpdate(collector.NewKeysAll([]string{"key1", "key2"}, tag.NewKey(ruler)), nil)
	_ = NewClean(
		&cleanResource{},
		[]transform.UpdateOperator{mode},
		WithCleanConfig(CleanConfig{
			SkipTruncate: false,
			AddExtraTags: false,
			Keys:         []string{"key3", "key4"},
		}),
	)
}

func ExampleNewInsert() {
	mode := transform.NewInsert(collector.NewKeysAll([]string{"key1", "key2"}, tag.NewKey(ruler)), nil)
	insMulti1 := transform.NewInsert(collector.NewKeysAll([]string{"key1", "key2"}, tag.NewKey(ruler)), nil)
	insMulti2 := transform.NewInsert(collector.NewKeysAll([]string{"key1", "key2"}, tag.NewKey(ruler)), nil)
	modeMulti := transform.NewInsertStack(insMulti1, insMulti2)
	insCross1 := transform.NewInsert(collector.NewKeysAll([]string{"key1", "key2"}, tag.NewMostKey(ruler)), nil)
	insCross2 := transform.NewInsert(collector.NewKeysAll([]string{"key1", "key2"}, tag.NewMostText(ruler)), nil)
	modeCross := transform.NewInsertCross(insCross1, insCross2)
	insSpread1 := transform.NewInsert(collector.NewKeysAll([]string{"key1", "key2"}, tag.NewKey(ruler)), nil)
	insSpread2 := transform.NewInsert(collector.NewKeysAll([]string{"key1", "key2"}, tag.NewKey(ruler)), nil)
	modeSpread := transform.NewInsertSpread(insSpread1, insSpread2)

	_ = NewInsert(&jobTarget{}, mode, WithInsertConfig(InsertConfig{
		SkipTruncate: false,
		ExtraKeys:    nil,
	}))

	_ = NewInsert(&jobTarget{}, modeMulti)
	_ = NewInsert(&jobTarget{}, modeCross)
	_ = NewInsert(&jobTarget{}, modeSpread)
}

func ExampleNewTransfer() {
	mode := transform.NewTransfer(
		[]string{"key1", "key2"},
		map[string]string{"key1": "alias1"},
		map[string]any{"fixed1": "fixed value"},
	)

	_ = NewTransfer(&jobTarget{}, mode)
}

func ExampleNewUpdate() {
	mode := transform.NewUpdate(collector.NewKeysAll([]string{"key1", "key2"}, tag.NewKey(ruler)), nil)

	_ = NewUpdate(&jobSource{}, []transform.UpdateOperator{mode})

	_ = NewUpdateTransfer(&jobSource{}, &jobTarget{}, []transform.UpdateOperator{mode}, WithUpdateTransferConfig(UpdateTransferConfig{
		SkipTruncate: false,
	}))
}
