package query

import (
	simpleDb "github.com/auho/go-simple-db/v2"
)

type sourcer interface {
	GetSheetName() string
	Rows() ([][]any, error)
}

type Source struct {
	SheetName string
	DB        *simpleDb.SimpleDB
}

func (s *Source) GetSheetName() string {
	return s.SheetName
}
