package reader

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/auho/go-etl/v3/insight/assistant/schema/create"
	"github.com/auho/go-toolkit/v2/farmtools/sort/maps"
)

var (
	intRe     = regexp.MustCompile(`^\d{1,10}$`)
	int64Re   = regexp.MustCompile(`^\d{11,20}$`)
	float64Re = regexp.MustCompile(`^\d+\.\d+$`)
)

// Schema reads an Excel sheet and builds a database table schema by inferring column types from data.
type Schema struct {
	excel  *Excel
	table  create.Tabler
	config Config

	titleFunc []func(string) string
}

// NewSchemaWithPath creates a Schema by opening the xlsx file at xlsxPath.
func NewSchemaWithPath(xlsxPath string, table create.Tabler, config Config) (*Schema, error) {
	excel, err := NewExcel(xlsxPath)
	if err != nil {
		return nil, err
	}

	return NewSchema(excel, table, config)
}

// NewSchema creates a Schema from an existing Excel.
func NewSchema(excel *Excel, table create.Tabler, config Config) (*Schema, error) {
	return &Schema{
		excel:  excel,
		table:  table,
		config: config,
	}, nil
}

// WithTitleFunc appends a title transformation function.
// All registered functions are applied in order to each column title.
func (s *Schema) WithTitleFunc(fn func(string) string) *Schema {
	s.titleFunc = append(s.titleFunc, fn)

	return s
}

// WithTitleAlias registers a title alias map as a transformation function.
// Each title is replaced by its alias if one is defined.
func (s *Schema) WithTitleAlias(alias map[string]string) *Schema {
	s.WithTitleFunc(func(s string) string {
		if _a, ok := alias[s]; ok {
			return _a
		} else {
			return s
		}
	})

	return s
}

// WithTitleAliasByIndex is not yet implemented.
func (s *Schema) WithTitleAliasByIndex() {
	// TODO implement me
	panic("implement me")
}

// BuildTable reads the sheet, infers column types from the data, and adds columns to the table schema.
// If EndRow is unset, defaults to 100 rows for type detection.
func (s *Schema) BuildTable() (create.Tabler, error) {
	if s.config.EndRow <= 0 {
		s.config.EndRow = 100
	}

	rows, err := s.excel.readSheet(s.config)
	if err != nil {
		return nil, fmt.Errorf("excel.readSheet: %w", err)
	}

	s.buildTable(rows)

	return s.table, err
}

// buildTable processes the first row as titles and infers column types from remaining rows.
func (s *Schema) buildTable(rows [][]string) {
	titles := rows[0]
	rows = rows[1:]

	_command := s.table.GetCommand()

	for i, title := range titles {
		for _, _fn := range s.titleFunc {
			title = _fn(title)
		}

		_type, _len1, _ := s.detectColumnType(i, rows)
		switch _type {
		case reflect.String:
			if _len1 <= 30 {
				_command.AddString(title)
			} else if _len1 <= 255 {
				_command.AddStringWithLength(title, 255)
			} else if _len1 <= 2000 {
				_command.AddStringWithLength(title, 2000)
			} else {
				_command.AddText(title)
			}
		case reflect.Int:
			_command.AddInt(title)
		case reflect.Int64:
			_command.AddBigInt(title)
		case reflect.Float64:
			_command.AddDecimal(title, 11, 2)
		default:
			panic("type not found")
		}
	}
}

// detectColumnType infers the reflect.Kind, max integer length, and max decimal scale
// for the column at index by scanning all rows.
func (s *Schema) detectColumnType(index int, rows [][]string) (reflect.Kind, int, int) {
	_types := make(map[reflect.Kind]int, len(rows))
	_len1 := 0
	_len2 := 0

	var _type reflect.Kind

	_vLen2 := 0
	for _, row := range rows {
		if index >= len(row) {
			continue
		}

		_value := row[index]
		_valueLen := len(_value)

		if intRe.MatchString(_value) {
			_type = reflect.Int
		} else if int64Re.MatchString(_value) {
			_type = reflect.Int64
		} else if float64Re.MatchString(_value) {
			_type = reflect.Float64

			_dotPos := strings.Index(_value, ".")
			_vLen2 = _valueLen - _dotPos - 1
			_valueLen = _dotPos
		} else {
			_type = reflect.String
		}

		if _valueLen > _len1 {
			_len1 = _valueLen
		}

		if _vLen2 > _len2 {
			_len2 = _vLen2
		}

		_types[_type] += 1
	}

	if len(_types) <= 0 {
		return reflect.String, 0, 0
	}

	_kt, _ := maps.SortValueDesc(_types)
	return _kt[0], _len1, _len2
}
