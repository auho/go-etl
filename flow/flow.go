package flow

import (
	"runtime"

	"github.com/auho/go-etl/action"
	"github.com/auho/go-etl/tool"
	goSimpleDb "github.com/auho/go-simple-db/v2"
	"github.com/auho/go-toolkit/flow"
	"github.com/auho/go-toolkit/flow/storage/database"
	"github.com/auho/go-toolkit/flow/storage/database/source"
)

func RunFlow(db *goSimpleDb.SimpleDB, dataTable string, idName string, actions []action.Actionor) {
	fields := []string{idName}
	for _, a := range actions {
		fields = append(fields, a.GetFields()...)
	}

	fields = tool.RemoveReplicaSliceString(fields)

	dataSource, err := source.NewSectionSliceMap(&source.QueryConfig{
		Config: source.Config{
			Concurrency: runtime.NumCPU(),
			PageSize:    2000,
			TableName:   dataTable,
			IdName:      idName,
		},
		Fields: fields,
	}, func() (*database.DB, error) {
		return database.NewFromSimpleDb(db), nil
	})
	if err != nil {
		panic(err)
	}

	opts := flow.Options[map[string]any]{}
	opts = append(opts, flow.WithSource[map[string]any](dataSource))
	for _, a := range actions {
		opts = append(opts, flow.WithTasker[map[string]any](a))
	}

	err = flow.RunFlow[map[string]any](opts...)
	if err != nil {
		panic(err)
	}
}
