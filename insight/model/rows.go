package model

var _ Rowsor = (*Rows)(nil)

type Rows struct {
	tableName string
}

func (r *Rows) TableName() string {
	return r.tableName
}
