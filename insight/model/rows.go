package model

type Rows struct {
	tableName string
}

func (r *Rows) TableName() string {
	return r.tableName
}
