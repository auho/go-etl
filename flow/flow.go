package flow

import (
	"runtime"

	"github.com/auho/go-etl/v2/action"
	"github.com/auho/go-etl/v2/tool"
	goSimpleDb "github.com/auho/go-simple-db/v2"
	"github.com/auho/go-toolkit/flow/action/singleton"
	"github.com/auho/go-toolkit/flow/flow"
	"github.com/auho/go-toolkit/flow/storage/database"
	"github.com/auho/go-toolkit/flow/storage/database/source"
)

func RunFlow(db *goSimpleDb.SimpleDB, dataTable string, idName string, actions []action.Actor) {
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

	var opts []flow.Option[map[string]any]
	opts = append(opts, flow.WithSource[map[string]any](dataSource))
	for _, a := range actions {
		opts = append(opts, flow.WithActor[map[string]any](singleton.NewActor[map[string]any](a)))
	}

	err = flow.RunFlow[map[string]any](opts...)
	if err != nil {
		panic(err)
	}
}
