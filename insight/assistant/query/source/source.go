package source

import (
	"fmt"
	"strings"
	"time"

	"github.com/auho/go-etl/v2/insight/assistant/sqlbuilder/dml"
	"github.com/auho/go-etl/v2/insight/assistant/query/dataset"
	maps "github.com/auho/go-etl/v2/tool/mapx"
	simpledb "github.com/auho/go-simple-db/v2"
)

type Source interface {
	Dataset() (*dataset.Dataset, error)
}

// SourceBase
// select data from db
type SourceBase struct {
	HasNamePrefix bool // Add the name prefix before the item
	Name          string
	Table         dml.Tabler
	DB            *simpledb.SimpleDB
}

func (s *SourceBase) itemValuesToIdentification(itemValues []string) string {
	id := s.keysToIdentification(itemValues)
	if s.HasNamePrefix {
		id = fmt.Sprintf("%s_%s", s.Name, id)
	}

	return id
}

func (s *SourceBase) keysToIdentification(keys []string) string {
	return strings.Join(keys, "_")
}

func (s *SourceBase) queryItemsSet(fields, itemsId []string, itemsSql map[string]string) ([]dataset.Set, error) {
	var sets []dataset.Set

	for _, itemId := range itemsId {
		rows, _d, err := s.querySql(itemsSql[itemId], fields)
		if err != nil {
			return nil, fmt.Errorf("querySql error; %w", err)
		}

		sets = append(sets, dataset.NewSetWithQuery(itemId, itemsSql[itemId], _d, rows))
	}

	return sets, nil
}

func (s *SourceBase) querySql(sql string, fields []string) ([][]any, time.Duration, error) {
	var rows []map[string]any

	_start := time.Now()
	err := s.DB.GormDB().Raw(sql).Scan(&rows).Error
	if err != nil {
		return nil, 0, fmt.Errorf("raw error; %w", err)
	}

	return maps.SliceMapStringAnyToSliceSliceAny(rows, fields), time.Now().Sub(_start), nil
}
