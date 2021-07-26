package flow

import (
	goEtl "github.com/auho/go-etl"
	"github.com/auho/go-etl/action"
	"github.com/auho/go-etl/storage"
)

type InsertFlow struct {
	actions      []*action.InsertAction
	source       storage.DbSource
	sourceConfig *storage.DbSourceConfig
}

func RunInsertFlow(config goEtl.DbConfig, dataName string, idName string, actions []*action.InsertAction) {
	iFlow := &InsertFlow{}

	iFlow.actions = actions

	iFlow.sourceConfig = storage.NewDbSourceConfig()
	iFlow.sourceConfig.MaxConcurrent = 4
	iFlow.sourceConfig.Size = 2000
	iFlow.sourceConfig.Table = dataName
	iFlow.sourceConfig.Driver = config.Driver
	iFlow.sourceConfig.Scheme = config.Scheme
	iFlow.sourceConfig.Table = config.Table
	iFlow.sourceConfig.PKeyName = idName

	fields := []string{idName}
	for _, a := range actions {
		fields = append(fields, a.GetFields()...)
	}

	fields = goEtl.RemoveReplicaSliceString(fields)

	iFlow.sourceConfig.Fields = fields

	iFlow.source.Start(iFlow.sourceConfig)

	for i := 0; i < 4; i++ {
		for {
			if items, ok := iFlow.source.Receive(); ok {
				for _, a := range actions {
					a.Receive(items)
				}
			} else {
				break
			}
		}
	}

	for _, a := range actions {
		a.SourceDone()
		a.Close()
	}
}
