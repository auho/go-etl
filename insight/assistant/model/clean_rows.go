package model

import (
	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/job"
)

var _ job.CleanResource = (*CleanRows)(nil)

type CleanRows struct {
	model

	raw     assistant.Entity
	rows    *Rows
	deleted *Rows
}

func NewCleanRows(newName string, raw assistant.Entity) *CleanRows {
	cd := &CleanRows{}
	cd.raw = raw
	cd.rows = NewRows(newName, raw.IDName(), raw.GetDB())
	cd.deleted = cd.rows.ToDeletedRows()

	return cd
}

func (cd *CleanRows) Rows() *Rows {
	return cd.rows
}

func (cd *CleanRows) DeletedRows() *Rows {
	return cd.deleted
}

func (cd *CleanRows) Source() job.Table {
	return cd.raw
}

func (cd *CleanRows) Data() job.Table {
	return cd.rows
}

func (cd *CleanRows) Deleted() job.Table {
	return cd.deleted
}
