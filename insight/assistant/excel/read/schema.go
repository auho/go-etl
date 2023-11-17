package read

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure/buildtable"
	"github.com/auho/go-toolkit/farmtools/sort/maps"
)

// Schema
// sheet to table schema
type Schema struct {
	excel  *Excel
	data   assistant.Rowsor
	config Config

	titleFunc func(string) string
}

func NewSchemaWithPath(xlsxPath string, data assistant.Rowsor, config Config) (*Schema, error) {
	excel, err := NewExcel(xlsxPath)
	if err != nil {
		return nil, err
	}

	return NewSchema(excel, data, config)
}

func NewSchema(excel *Excel, data assistant.Rowsor, config Config) (*Schema, error) {
	return &Schema{
		excel:  excel,
		data:   data,
		config: config,
	}, nil
}

func (s *Schema) WithFuncTitle(fn func(string) string) *Schema {
	s.titleFunc = fn

	return s
}

func (s *Schema) BuildTable() (*buildtable.RawTable, error) {
	if s.config.EndRow <= 0 {
		s.config.EndRow = 100
	}

	rows, err := s.excel.readSheet(s.config)
	if err != nil {
		return nil, fmt.Errorf("readSheet error; %w", err)
	}

	_table := buildtable.NewRawTable(s.data)
	s.buildTable(_table, rows)

	return _table, err
}

func (s *Schema) buildTable(table *buildtable.RawTable, rows [][]string) {
	titles := rows[0]
	rows = rows[1:]

	for i, title := range titles {
		if s.titleFunc != nil {
			title = s.titleFunc(title)
		}

		_type, _len1, _ := s.detectColumnType(i, rows)
		switch _type {
		case reflect.String:
			if _len1 <= 30 {
				table.AddString(title)
			} else if _len1 <= 255 {
				table.AddStringWithLength(title, 255)
			} else if _len1 <= 2000 {
				table.AddStringWithLength(title, 2000)
			} else {
				table.AddText(title)
			}
		case reflect.Int:
			table.AddInt(title)
		case reflect.Int64:
			table.AddBigInt(title)
		case reflect.Float64:
			table.AddDecimal(title, 11, 2)
		default:
			panic("type not found")
		}
	}
}

func (s *Schema) detectColumnType(index int, rows [][]string) (reflect.Kind, int, int) {
	intRe := regexp.MustCompile(`^\d{1,10}$`)
	int64Re := regexp.MustCompile(`^\d{11,20}$`)
	float64Re := regexp.MustCompile(`^\d+\.\d+$`)

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

	return maps.SorterByValueDesc(_types).Keys()[0], _len1, _len2
}
