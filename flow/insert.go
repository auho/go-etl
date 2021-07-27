package flow

import (
	"sync"

	goEtl "github.com/auho/go-etl"
	"github.com/auho/go-etl/action"
	"github.com/auho/go-etl/storage"
)

func RunInsertFlow(config goEtl.DbConfig, dataName string, idName string, actions []*action.InsertAction) {
	var wg sync.WaitGroup

	fields := []string{idName}
	for _, a := range actions {
		fields = append(fields, a.GetFields()...)
	}
	fields = goEtl.RemoveReplicaSliceString(fields)

	sourceConfig := storage.NewDbSourceConfig()
	sourceConfig.MaxConcurrent = 4
	sourceConfig.Size = 2000
	sourceConfig.Table = dataName
	sourceConfig.Driver = config.Driver
	sourceConfig.Dsn = config.Dsn
	sourceConfig.PKeyName = idName
	sourceConfig.Fields = fields

	source := storage.NewDbSource(sourceConfig)
	source.Start()

	for _, a := range actions {
		a.Start()
	}

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			for {
				if items, ok := source.Consume(); ok {
					for _, a := range actions {
						a.Receive(items)
					}
				} else {
					break
				}
			}

			wg.Done()
		}()
	}

	wg.Wait()

	for _, a := range actions {
		a.SourceDone()
		a.Close()
	}
}
