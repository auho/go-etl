package entity

import (
	"github.com/auho/go-etl/v3/task/pipeline/database/etl"
)

var _ etl.CleanResource = (*CleanData)(nil)

type CleanData struct {
	base

	rows    *Rows
	data    *Data
	deleted *Rows
}

func NewCleanData(rows *Rows) *CleanData {
	cd := &CleanData{}
	cd.rows = rows
	cd.data = rows.ToData()
	cd.deleted = rows.ToDeletedRows()

	return cd
}

func (cd *CleanData) Source() etl.Table {
	return cd.rows
}

func (cd *CleanData) Data() etl.Table {
	return cd.data
}

func (cd *CleanData) Deleted() etl.Table {
	return cd.deleted
}
