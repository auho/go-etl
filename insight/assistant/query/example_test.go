package query

import (
	"log"

	"github.com/auho/go-etl/v2/insight/assistant/query/source"
)

func ExampleNewQuery() {
	_q, err := NewQuery("xlsxName", "xlsxPath")
	if err != nil {
		log.Fatalln(err)
	}

	_q.AddAppend(source.NewRows(source.Source{
		HasNamePrefix: false,
		Name:          "one",
		Table:         nil,
		DB:            nil,
	}))

	_q.AddSpread(source.NewPlaceholder(source.Source{
		HasNamePrefix: false,
		Name:          "",
		Table:         nil,
		DB:            nil,
	}).WithItems(nil))

	err = _q.Save()
	if err != nil {
		log.Fatalln(err)
	}
}
