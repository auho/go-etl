package task

import (
	"runtime"

	"github.com/auho/go-etl/v2/job"
	"github.com/auho/go-etl/v2/job/action"
	"github.com/auho/go-etl/v2/tool/slices"
	"github.com/auho/go-toolkit/flow/action/singleton"
	"github.com/auho/go-toolkit/flow/flow"
	"github.com/auho/go-toolkit/flow/storage/database"
	"github.com/auho/go-toolkit/flow/storage/database/source"
)

func RunTask(aSource job.Source, actions []action.Actor) {
	fields := []string{aSource.GetIdName()}
	for _, a := range actions {
		fields = append(fields, a.GetFields()...)
	}

	fields = slices.SliceDropDuplicates(fields)

	dataSource, err := source.NewSectionSliceMap(&source.QueryConfig{
		Config: source.Config{
			Concurrency: runtime.NumCPU(),
			PageSize:    2000,
			TableName:   aSource.TableName(),
			IdName:      aSource.GetIdName(),
		},
		Fields: fields,
	}, func() (*database.DB, error) {
		return database.NewFromSimpleDb(aSource.GetDB()), nil
	})
	if err != nil {
		panic(err)
	}

	opts := []flow.Option[map[string]any]{
		flow.WithSource[map[string]any](dataSource),
	}

	for _, a := range actions {
		opts = append(opts, flow.WithActor[map[string]any](singleton.NewActor[map[string]any](a)))
	}

	err = flow.RunFlow[map[string]any](opts...)
	if err != nil {
		panic(err)
	}
}
