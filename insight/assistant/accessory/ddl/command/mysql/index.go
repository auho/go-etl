package mysql

import (
	"fmt"
	"strings"
)

type Index struct {
	name   string
	_type  string
	fields []string
}

func newIndex(name, _type string) *Index {
	return &Index{
		name:  name,
		_type: _type,
	}
}

func (i *Index) AddField(field string, size int) *Index {
	s := fmt.Sprintf("`%s`", field)
	if size > 0 {
		s = fmt.Sprintf("%s(%d)", s, size)
	}

	i.fields = append(i.fields, s)

	return i
}

func (i *Index) SqlForCreateTable() string {
	var ss []string
	for _, field := range i.fields {
		ss = append(ss, fmt.Sprintf("%s", field))
	}

	return fmt.Sprintf("%s `%s` (%s)", i._type, i.name, strings.Join(ss, ","))
}

func (i *Index) SqlForAdd(tableName string) string {
	return fmt.Sprintf("ALTER TABLE `%s` ADD %s", tableName, i.SqlForCreateTable())
}
