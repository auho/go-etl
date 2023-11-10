package dataset

import (
	"fmt"
	"time"
)

type Set struct {
	ItemName string
	Sql      string
	Amount   int
	Rows     [][]any
	Duration time.Duration
}

type Dataset struct {
	Name   string   // dataset name
	Titles []string // dataset item data title
	Sets   []Set
}

type Mode string

const ModeAppend Mode = "append"
const ModeSpread Mode = "spread"

type Moder interface {
	Data() (*Data, error)
	Name() string
	Sets() []Set
}

type Data struct {
	Names []string           // 保存 name 的顺序
	Rows  map[string][][]any // map[name][][]any
}

func (d *Data) Add(name string, rows [][]any) {
	if len(d.Rows) <= 0 {
		d.Rows = make(map[string][][]any)
	}

	d.Names = append(d.Names, name)
	d.Rows[name] = rows
}

func NewMode(mode Mode, ds *Dataset) (Moder, error) {
	var dsMode Moder
	switch mode {
	case ModeAppend:
		dsMode = NewAppendMode(ds)
	case ModeSpread:
		dsMode = NewSpreadMode(ds)
	default:
		return nil, fmt.Errorf("dataset mode[%s] error", mode)
	}

	return dsMode, nil
}
