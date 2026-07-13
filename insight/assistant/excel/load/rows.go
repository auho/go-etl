package load

import (
	"github.com/auho/go-etl/v3/insight/assistant"
	"github.com/auho/go-etl/v3/insight/assistant/excel/reader"
	"github.com/auho/go-etl/v3/insight/assistant/schema/create"
	simpledb "github.com/auho/go-simple-db/v3"
)

var _ Resource = (*RowsResource)(nil)

type RowsResource struct {
	baseResource
	Titles // column title of save to db

	Rows assistant.Entity
}

func (r *RowsResource) Prepare() error {
	return r.Titles.prepare()
}

func (r *RowsResource) Name() string {
	return r.Rows.Name()
}

func (r *RowsResource) Tabler() create.Tabler {
	return create.NewRowsTable(r.Rows)
}

func (r *RowsResource) SheetData(excel *reader.Excel) (reader.SheetDataReader, error) {
	return r.readSheetData(excel, r.buildSheetConfig())
}

func (r *RowsResource) DB() *simpledb.SimpleDB {
	return r.Rows.DB()
}
