package source

import (
	"fmt"
	"strings"
	"time"

	"github.com/auho/go-etl/v2/insight/assistant/accessory/dml"
	"github.com/auho/go-etl/v2/insight/assistant/query/dataset"
	"github.com/auho/go-etl/v2/tool/maps"
	simpleDb "github.com/auho/go-simple-db/v2"
)

type Sourcer interface {
	Dataset() (*dataset.Dataset, error)
}

// Source
// select data from db
type Source struct {
	HasNamePrefix bool // Add the name prefix before the item
	Name          string
	Table         dml.Tabler
	DB            *simpleDb.SimpleDB
}

func (s *Source) itemsNameToIdentification(itemsName []string) string {
	id := s.keysToIdentification(itemsName)
	if s.HasNamePrefix {
		id = fmt.Sprintf("%s_%s", s.Name, id)
	}

	return id
}

func (s *Source) keysToIdentification(keys []string) string {
	return strings.Join(keys, "_")
}

func (s *Source) queryItemsSet(fields, itemsName []string, itemsSql map[string]string) ([]dataset.Set, error) {
	var sets []dataset.Set

	for _, itemName := range itemsName {
		rows, _d, err := s.querySql(itemsSql[itemName], fields)
		if err != nil {
			return nil, fmt.Errorf("querySql error; %w", err)
		}

		sets = append(sets, dataset.NewSetWithQuery(itemName, itemsSql[itemName], _d, rows))
	}

	return sets, nil
}

func (s *Source) querySql(sql string, fields []string) ([][]any, time.Duration, error) {
	var rows []map[string]any

	_start := time.Now()
	err := s.DB.Raw(sql).Scan(&rows).Error
	if err != nil {
		return nil, 0, fmt.Errorf("raw error; %w", err)
	}

	return maps.SliceMapStringAnyToSliceSliceAny(rows, fields), time.Now().Sub(_start), nil
}
