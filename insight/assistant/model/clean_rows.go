package model

import (
	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/job"
)

var _ job.CleanResource = (*CleanRows)(nil)

type CleanRows struct {
	model

	raw     assistant.Moder
	rows    *Rows
	deleted *Rows
}

func NewCleanRows(newName string, raw assistant.Moder) *CleanRows {
	cd := &CleanRows{}
	cd.raw = raw
	cd.rows = NewRows(newName, raw.GetIDName(), raw.GetDB())
	cd.deleted = cd.rows.ToDeletedRows()

	return cd
}

func (cd *CleanRows) Rows() *Rows {
	return cd.rows
}

func (cd *CleanRows) DeletedRows() *Rows {
	return cd.deleted
}

func (cd *CleanRows) Source() job.Target {
	return cd.raw
}

func (cd *CleanRows) Data() job.Target {
	return cd.rows
}

func (cd *CleanRows) Deleted() job.Target {
	return cd.deleted
}
