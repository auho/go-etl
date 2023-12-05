package query

import (
	"log"

	"github.com/auho/go-etl/v2/insight/assistant/query/source"
)

func ExampleNewQuery() {
	_q, err := NewQuery("xlsxName", "xlsxDir")
	if err != nil {
		log.Fatalln(err)
	}

	// append mode
	_q.AddAppend(source.NewRows(source.Source{
		Name:  "one",
		Table: nil,
		DB:    nil,
	}))

	// rows stack
	_q.AddAppend(source.NewRowsStack("name",
		source.Source{
			Table: nil,
			DB:    nil,
		},
		source.Source{
			Table: nil,
			DB:    nil,
		},
	))

	// spread mode
	_q.AddSpread(source.NewPlaceholder(source.Source{
		HasNamePrefix: false,
		Name:          "",
		Table:         nil,
		DB:            nil,
	}).WithItems(nil))

	// save
	err = _q.Save()
	if err != nil {
		log.Fatalln(err)
	}
}

func ExampleNewQueryWithPath() {
	_q, err := NewQueryWithPath("xlsxPath")
	if err != nil {
		log.Fatalln(err)
	}

	// append mode
	_q.AddAppend(source.NewRows(source.Source{
		Name:  "one",
		Table: nil,
		DB:    nil,
	}))

	// rows stack
	_q.AddAppend(source.NewRowsStack("name",
		source.Source{
			Table: nil,
			DB:    nil,
		},
		source.Source{
			Table: nil,
			DB:    nil,
		},
	))

	// spread mode
	_q.AddSpread(source.NewPlaceholder(source.Source{
		HasNamePrefix: false,
		Name:          "",
		Table:         nil,
		DB:            nil,
	}).WithItems(nil))

	// save
	err = _q.Save()
	if err != nil {
		log.Fatalln(err)
	}
}
