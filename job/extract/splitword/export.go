package splitword

import (
	"github.com/auho/go-etl/v3/job/extract"
)

func NewExportAll() *extract.Exporter[Results] {
	format := DefaultFormat
	return extract.NewExporter(
		map[string]any{format.WordName: ""},
		func(ctx extract.ExportContext[Results]) []map[string]any {
			return ctx.Results.ToAll(ctx.Format.(Format))
		},
		extract.WithFormat[Results](format),
	)
}

func NewExportLine() *extract.Exporter[Results] {
	format := DefaultFormat
	return extract.NewExporter(
		map[string]any{format.WordName: ""},
		func(ctx extract.ExportContext[Results]) []map[string]any {
			return ctx.Results.ToLine(ctx.Format.(Format))
		},
		extract.WithFormat[Results](format),
	)
}
