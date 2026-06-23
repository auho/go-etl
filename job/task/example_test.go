package task

import (
	"github.com/auho/go-etl/v2/job"
	"github.com/auho/go-etl/v2/job/means"
	"github.com/auho/go-etl/v2/job/means/tag"
	"github.com/auho/go-etl/v2/job/mode"
	simpledb "github.com/auho/go-simple-db/v2"
)

var _ job.Source = (*_jobSource)(nil)
var _ job.Target = (*_jobTarget)(nil)
var _ job.CleanResource = (*_cleanResource)(nil)

var _ruler means.Ruler

type _jobSource struct{}

func (_ _jobSource) GetIDName() string         { return "id" }
func (_ _jobSource) TableName() string         { return "source" }
func (_ _jobSource) GetDB() *simpledb.SimpleDB { return nil }

type _jobTarget struct{}

func (_ _jobTarget) GetIDName() string         { return "id" }
func (_ _jobTarget) TableName() string         { return "target" }
func (_ _jobTarget) GetDB() *simpledb.SimpleDB { return nil }

type _cleanResource struct{}

func (_ _cleanResource) Source() job.Target  { return &_jobSource{} }
func (_ _cleanResource) Data() job.Target    { return &_jobTarget{} }
func (_ _cleanResource) Deleted() job.Target { return &_jobTarget{} }

func ExampleNewClean() {
	_mode := mode.NewUpdate([]string{"key1", "key2"}, tag.NewKey(_ruler).ToMeans())
	_ = NewClean(
		&_cleanResource{},
		[]mode.UpdateModer{_mode},
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

	_mode := mode.NewInsert([]string{"key1", "key2"}, tag.NewKey(_ruler).ToMeans())
	_modeMulti := mode.NewInsertStack([]string{"key1", "key2"}, tag.NewKey(_ruler).ToMeans(), tag.NewLabel(_ruler).ToMeans())
	_modeCross := mode.NewInsertCross([]string{"key1", "key2"}, tag.NewMostKey(_ruler).ToMeans(), tag.NewMostText(_ruler).ToMeans())
	_modeSpread := mode.NewInsertSpread([]string{"key1", "key2"}, tag.NewKey(_ruler).ToMeans(), tag.NewKey(_ruler).ToMeans())

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
	_mode := mode.NewTransfer(
		[]string{"key1", "key2"},
		map[string]string{"key1": "alias1"},
		map[string]any{"fixed1": "fixed value"},
	)

	_ = NewTransfer(&_jobTarget{}, _mode)
}

func ExampleNewUpdate() {
	_mode := mode.NewUpdate([]string{"key1", "key2"}, tag.NewKey(_ruler).ToMeans(), tag.NewLabel(_ruler).ToMeans())

	_ = NewUpdate(&_jobSource{}, []mode.UpdateModer{_mode})

	_ = NewUpdateAndTransfer(&_jobSource{}, &_jobTarget{}, []mode.UpdateModer{_mode}, WithUpdateTransferConfig(UpdateTransferConfig{
		NotTruncate: false,
		BatchSize:   0,
		Concurrency: 0,
	}))
}
