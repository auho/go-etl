package load

import (
	"github.com/auho/go-etl/v3/insight/assistant"
	"github.com/auho/go-etl/v3/insight/assistant/excel/reader"
	"github.com/auho/go-etl/v3/insight/assistant/schema/create"
	simpledb "github.com/auho/go-simple-db/v3"
)

var _ Resource = (*Rows)(nil)

type Rows struct {
	baseResource
	Titles // column title of save to db

	Rows assistant.Entity
}

func (r *Rows) Prepare() error {
	return r.Titles.prepare()
}

func (r *Rows) Name() string {
	return r.Rows.Name()
}

func (r *Rows) Tabler() create.Tabler {
	return create.NewRowsTable(r.Rows)
}

func (r *Rows) SheetData(excel *reader.Excel) (reader.SheetDataReader, error) {
	return r.readSheetData(excel, r.buildSheetConfig())
}

func (r *Rows) DB() *simpledb.SimpleDB {
	return r.Rows.DB()
}
