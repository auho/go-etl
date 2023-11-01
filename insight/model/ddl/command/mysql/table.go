package mysql

import (
	"fmt"
	"strconv"
	"strings"
)

type Table struct {
	name    string
	engine  string
	charset string
	collate string

	fields      []*Field
	primaryKeys []string
	indexes     []*Index
}

func (t *Table) GetName() string {
	return t.name
}

func (t *Table) setTable(name, engine, charset, collate string) *Table {
	return t.SetName(name).SetEngine(engine).SetCharset(charset, collate)
}

func (t *Table) SetName(name string) *Table {
	t.name = name

	return t
}

func (t *Table) SetEngine(engine string) *Table {
	if engine == "" {
		engine = engineMyISAM
	}

	t.engine = engine

	return t
}

func (t *Table) SetEngineMyISAM() *Table {
	t.engine = engineMyISAM

	return t
}

func (t *Table) SetEngineInnoDB() *Table {
	t.engine = engineInnoDB

	return t
}

func (t *Table) SetCharset(charset, collate string) *Table {
	if charset == "" || collate == "" {
		charset = charsetUtf8mb4
		collate = collateUtf8mb4GeneralCi
	}

	t.charset = charset
	t.collate = collate

	return t
}

func (t *Table) addInteger(name, _type string, length int, _default int, isUnsigned bool) *Field {
	if isUnsigned && _default <= 0 {
		_default = 0
	}

	_sLength := strconv.Itoa(length)
	_sDefault := strconv.Itoa(_default)

	f := newField(name, _type, _sLength, _sDefault, "")
	f.SetUnsigned(isUnsigned)

	t.AddField(f)
	return f
}

func (t *Table) addString(name, _type string, length int, _default string) *Field {
	_sLength := strconv.Itoa(length)

	f := newField(name, _type, _sLength, _default, "")

	t.AddField(f)
	return f
}

func (t *Table) AddPKInt(name string, length int) *Field {
	f := t.AddInt(name, length, 0, true).SetExtra(extraAutoIncrement)
	t.primaryKeys = append(t.primaryKeys, name)

	return f
}

func (t *Table) AddPkBigInt(name string, length int) *Field {

	f := t.AddBigInt(name, length, 0, true).SetExtra(extraAutoIncrement)
	t.primaryKeys = append(t.primaryKeys, name)

	return f
}

func (t *Table) AddPkVarchar(name string, length int) *Field {
	f := t.AddVarchar(name, length, "")
	t.primaryKeys = append(t.primaryKeys, name)

	return f
}

func (t *Table) AddBigInt(name string, length int, _default int, isUnsigned bool) *Field {
	return t.addInteger(name, typeBigInt, length, _default, isUnsigned)
}

func (t *Table) AddInt(name string, length int, _default int, isUnsigned bool) *Field {
	return t.addInteger(name, typeInt, length, _default, isUnsigned)
}

func (t *Table) AddVarchar(name string, length int, _default string) *Field {
	return t.addString(name, typeVarchar, length, _default)
}

func (t *Table) AddText(name string) *Field {
	f := newField(name, typeText, "", "", "")

	t.AddField(f)
	return f
}

func (t *Table) AddDecimal(name string, m, d int, _default float64) *Field {
	_sDefault := strconv.FormatFloat(_default, 'g', -1, 64)

	f := newField(name, typeDecimal, fmt.Sprintf("%d,%d", m, d), _sDefault, "")

	t.AddField(f)
	return f
}

func (t *Table) AddTimestamp(name string, onDefault, onUpdate bool) *Field {
	_default := ""
	if onDefault {
		_default = defaultCurrentTimestamp
	}

	extra := ""
	if onUpdate {
		extra = extraOnUpdateCurrentTimestamp
	}

	f := newField(name, typeTimestamp, "", _default, extra)
	t.AddField(f)
	return f
}

func (t *Table) AddKey(field string, size int) *Index {
	i := newIndex(field, indexTypeKey)
	i.AddField(field, size)

	t.AddIndex(i)
	return i
}

func (t *Table) AddUniqueKey(field string) *Index {
	i := newIndex(field, indexTypeUniqueKey)
	i.AddField(field, 0)

	t.AddIndex(i)
	return i
}

func (t *Table) AddIndex(index *Index) *Index {
	t.indexes = append(t.indexes, index)

	return index
}

func (t *Table) AddField(field *Field) *Table {
	t.fields = append(t.fields, field)

	return t
}

func (t *Table) SqlForCreate() string {
	var columnStringList []string
	for _, field := range t.fields {
		columnStringList = append(columnStringList, field.SqlForCreateTable())
	}

	if len(t.primaryKeys) > 0 {
		var pks []string
		for _, key := range t.primaryKeys {
			pks = append(pks, fmt.Sprintf("`%s`", key))
		}

		columnStringList = append(columnStringList, fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(pks, ",")))
	}

	for _, index := range t.indexes {
		columnStringList = append(columnStringList, index.SqlForCreateTable())
	}

	if t.charset == "" || t.collate == "" {
		t.charset = charsetUtf8mb4
		t.collate = collateUtf8mb4GeneralCi
	}

	sql := fmt.Sprintf("CREATE TABLE `%s`(\n%s\n)ENGINE=%s DEFAULT CHARSET=%s COLLATE=%s",
		t.name,
		strings.Join(columnStringList, ",\n"),
		t.engine,
		t.charset,
		t.collate,
	)

	return sql
}
