package mysql

import (
	"fmt"
)

type Field struct {
	name     string
	_type    string
	collate  string
	length   string
	_default string
	extra    string

	isUnsigned  bool
	isAllowNull bool

	after string
}

func newField(name, _type string, length string, _default, extra string) *Field {
	f := &Field{}

	return f.SetName(name).SetType(_type).SetLength(length).SetExtra(extra).SetDefault(_default)
}

func (f *Field) SetName(name string) *Field {
	f.name = name

	return f
}

func (f *Field) SetType(_type string) *Field {
	f._type = _type

	return f
}

func (f *Field) SetLength(length string) *Field {
	f.length = length

	return f
}

func (f *Field) SetExtra(extra string) *Field {
	f.extra = extra

	return f
}

func (f *Field) SetExtraAutoIncrement() *Field {
	f.extra = extraAutoIncrement

	return f
}

func (f *Field) SetDefault(_default string) *Field {
	f._default = _default

	return f
}

func (f *Field) SetUnsigned(b bool) *Field {
	f.isUnsigned = b

	return f
}

func (f *Field) SetCollate(collate string) *Field {
	f.collate = collate

	return f
}

func (f *Field) SetCollateUtf8mb4Bin() *Field {
	return f.SetCollate(collateUtf8mb4Bin)
}

func (f *Field) SetAllowNull(b bool) *Field {
	f.isAllowNull = b

	return f
}

func (f *Field) SetAfter(name string) *Field {
	f.after = name

	return f
}

func (f *Field) statement() string {
	name := fmt.Sprintf("`%s` ", f.name)

	length := ""
	if f.length != "" {
		length = fmt.Sprintf("(%s)", f.length)
	}
	_type := fmt.Sprintf("%s%s ", f._type, length)

	unsigned := ""
	switch f._type {
	case typeInt, typeBigInt, typeDecimal:
		if f.isUnsigned {
			unsigned = "unsigned "
		}
	}

	flag := fmt.Sprintf("%s", unsigned)

	collate := ""
	if f.collate != "" {
		collate = fmt.Sprintf("COLLATE %s ", f.collate)
	}

	null := "NOT NULL "
	if f.isAllowNull {
		null = "NULL "
	}

	_default := fmt.Sprintf("DEFAULT '%s' ", f._default)

	switch f._type {
	case typeText:
		_default = ""
		null = ""
	case typeTimestamp:
		if f._default == "" {
			_default = "DEFAULT NULL "
			null = "NULL "
		} else {
			_default = fmt.Sprintf("DEFAULT %s ", defaultCurrentTimestamp)
		}
	}

	extra := f.extra
	if f.extra == extraAutoIncrement {
		_default = ""
	}

	// name | type | flag | collate | null | default | extra
	return fmt.Sprintf("%s%s%s%s%s%s%s", name, _type, flag, collate, null, _default, extra)
}

func (f *Field) SqlForCreateTable() string {
	return f.statement()
}

func (f *Field) SqlForAdd(tableName string) string {
	_after := ""
	if f.after != "" {
		_after = fmt.Sprintf(" AFTER `%s`", f.after)
	}

	return fmt.Sprintf("ALTER TABLE `%s` ADD %s%s", tableName, f.statement(), _after)
}

func (f *Field) SqlForModify() string {
	// TODO
	return ""
}

func (f *Field) SqlForChange(tableName string) string {
	return fmt.Sprintf("ALTER TABLE `%s` MODIFY COLUMN `%s`", tableName, f.statement())
}
