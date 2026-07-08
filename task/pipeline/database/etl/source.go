package etl

import (
	"fmt"
	"runtime"

	"github.com/auho/go-toolkit-flow/v3/storage"
	"github.com/auho/go-toolkit-flow/v3/storage/database/source"
)

// SourceConfig configures the source reader (pagination + concurrency).
type SourceConfig struct {
	Concurrency int
	Maximum     int64
	PageSize    int64
}

func (sc *SourceConfig) Check() {
	if sc.Concurrency <= 0 {
		sc.Concurrency = runtime.NumCPU()
	}
	if sc.PageSize <= 0 {
		sc.PageSize = batchSize
	}
}

// newSource builds a Gorm paginated source over table (etl internal helper).
func newSource(table Table, fields []string, cfg SourceConfig) (storage.Source[storage.MapEntry], error) {
	s, err := source.NewSectionMapWithGorm(
		source.SectionConfig{
			Concurrency: cfg.Concurrency,
			MaxItems:    cfg.Maximum,
			StartID:     0,
			EndID:       0,
			PageSize:    cfg.PageSize,
		},
		source.ScanConfig{
			TableName:     table.TableName(),
			SegmentIDName: table.IDName(),
			Where:         "",
			Order:         "",
			SelectFields:  fields,
			WhereArgs:     nil,
		},
		table.GetDB().GormDB(),
	)
	if err != nil {
		return nil, fmt.Errorf("NewSectionMapWithGorm: %w", err)
	}
	return s, nil
}
