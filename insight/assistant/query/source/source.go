package source

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/auho/go-etl/v3/insight/assistant/query/dataset"
	"github.com/auho/go-etl/v3/insight/assistant/sqlbuilder/dml"
	"github.com/auho/go-etl/v3/tool/mapx"
	simpledb "github.com/auho/go-simple-db/v3"
)

type Source interface {
	Dataset() (*dataset.Dataset, error)
}

// Base
// select data from db
type Base struct {
	HasNamePrefix bool // Add the name prefix before the item
	Name          string
	Table         dml.Tabler
	DB            *simpledb.SimpleDB
}

func (b *Base) itemValuesToIdentification(itemValues []string) string {
	id := b.valuesToIdentification(itemValues)
	if b.HasNamePrefix {
		id = fmt.Sprintf("%s_%s", b.Name, id)
	}

	return id
}

func (b *Base) valuesToIdentification(values []string) string {
	return strings.Join(values, "_")
}

func (b *Base) queryItemsSet(fields, itemsId []string, itemsSql map[string]string) ([]dataset.Set, error) {
	var sets []dataset.Set

	for _, itemId := range itemsId {
		rows, _d, err := b.querySql(itemsSql[itemId], fields)
		if err != nil {
			return nil, fmt.Errorf("querySql: %w", err)
		}

		sets = append(sets, dataset.NewSetWithQuery(itemId, itemsSql[itemId], _d, rows))
	}

	return sets, nil
}

func (b *Base) querySql(sql string, fields []string) ([][]any, time.Duration, error) {
	var rows []map[string]any

	_start := time.Now()
	err := b.DB.GormDB().WithContext(context.TODO()).Raw(sql).Scan(&rows).Error
	if err != nil {
		return nil, 0, fmt.Errorf("scan: %w", err)
	}

	return mapx.SliceMapStringAnyToSliceSliceAny(rows, fields), time.Since(_start), nil
}
