package entity

import (
	"github.com/auho/go-etl/v3/insight/assistant"
	"github.com/auho/go-etl/v3/task"
)

var _ task.CleanResource = (*CleanRows)(nil)

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

func (cd *CleanRows) Source() task.Table {
	return cd.raw
}

func (cd *CleanRows) Data() task.Table {
	return cd.rows
}

func (cd *CleanRows) Deleted() task.Table {
	return cd.deleted
}
