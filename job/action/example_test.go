package action

import (
	"github.com/auho/go-etl/v2/job"
	"github.com/auho/go-etl/v2/job/means"
	"github.com/auho/go-etl/v2/job/means/tag"
	"github.com/auho/go-etl/v2/job/mode"
	simpleDb "github.com/auho/go-simple-db/v2"
)

var _ job.Source = (*_source)(nil)
var _ job.Target = (*_target)(nil)
var _ job.CleanResource = (*_cleanResource)(nil)

var _ruler means.Ruler

type _source struct{}

func (_ _source) GetIdName() string         { return "id" }
func (_ _source) TableName() string         { return "source" }
func (_ _source) GetDB() *simpleDb.SimpleDB { return nil }

type _target struct{}

func (_ _target) GetIdName() string         { return "id" }
func (_ _target) TableName() string         { return "target" }
func (_ _target) GetDB() *simpleDb.SimpleDB { return nil }

type _cleanResource struct{}

func (_ _cleanResource) SourceTarget() job.Target  { return &_source{} }
func (_ _cleanResource) DataTarget() job.Target    { return &_target{} }
func (_ _cleanResource) DeletedTarget() job.Target { return &_target{} }

func ExampleNewClean() {
	_mode := mode.NewUpdate([]string{"key1", "key2"}, tag.NewKey(_ruler))
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

	_mode := mode.NewInsert([]string{"key1", "key2"}, tag.NewKey(_ruler))
	_modeMulti := mode.NewInsertMulti([]string{"key1", "key2"}, tag.NewKey(_ruler), tag.NewLabel(_ruler))
	_modeCross := mode.NewInsertCross([]string{"key1", "key2"}, tag.NewMostKey(_ruler), tag.NewMostText(_ruler))
	_modeSpread := mode.NewInsertSpread([]string{"key1", "key2"}, tag.NewKey(_ruler), tag.NewKey(_ruler))

	_ = NewInsert(&_target{}, _mode, WithInsertConfig(InsertConfig{
		NotTruncate: false,
		BatchSize:   0,
		Concurrency: 0,
		ExtraKeys:   nil,
	}))

	_ = NewInsert(&_target{}, _modeMulti)
	_ = NewInsert(&_target{}, _modeCross)
	_ = NewInsert(&_target{}, _modeSpread)
}

func ExampleNewTransfer() {
	_mode := mode.NewTransfer(
		[]string{"key1", "key2"},
		map[string]string{"key1": "alias1"},
		map[string]any{"fixed1": "fixed value"},
	)

	_ = NewTransfer(&_target{}, _mode)
}

func ExampleNewUpdate() {
	_mode := mode.NewUpdate([]string{"key1", "key2"}, tag.NewKey(_ruler), tag.NewLabel(_ruler))

	_ = NewUpdate(&_source{}, []mode.UpdateModer{_mode})

	_ = NewUpdateAndTransfer(&_source{}, &_target{}, []mode.UpdateModer{_mode}, WithUpdateConfig(UpdateConfig{
		NotTruncate: false,
		BatchSize:   0,
		Concurrency: 0,
	}))
}
