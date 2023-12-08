package tag

import (
	"github.com/auho/go-etl/v2/job/explore/token"
)

var _ token.Tokenizer = (*ResultsToken)(nil)

var exporters map[string]token.Exporter

func RegisterExporter(exporter token.Exporter) {
	exporters[exporter.Name()] = exporter
}

type ResultsToken struct {
}

func (r *ResultsToken) ToExport(m string) token.Exporter {
	//TODO implement me
	panic("implement me")
}

func (r *ResultsToken) GetOk() bool {
	//TODO implement me
	panic("implement me")
}
