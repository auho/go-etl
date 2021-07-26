package action

import (
	"sync"

	"github.com/auho/go-etl/flow/mode"
	"github.com/auho/go-etl/storage"
)

type InsertAction struct {
	source      *storage.DbSource
	target      *storage.DbTargetInsertInterface
	mode        mode.InsertMode
	idName      string
	dataName    string
	affixFields []string
	wg          sync.WaitGroup
}

func NewInsertAction() *InsertAction {
	ia := &InsertAction{}

	return ia
}

func (ia *InsertAction) Start(dataName string, idName string, affixFields []string, tagTableName string) {
	for i := 0; i < 4; i++ {
		ia.wg.Add(1)
		go ia.doSource()
	}

	sourceConfig := storage.NewDbSourceConfig()
	sourceConfig.MaxConcurrent = 4
	sourceConfig.Size = 2000
	sourceConfig.Table = dataName

	targetConfig := storage.NewDbTargetConfig()
	ia.source.Start(sourceConfig)
	ia.target.Start(targetConfig)
}

func (ia *InsertAction) Close() {
	ia.wg.Wait()

	ia.target.Done()
	ia.target.Close()
}

func (ia *InsertAction) GetKeys() []string {
	return append(ia.mode.GetKeys(), ia.affixFields...)
}

func (ia *InsertAction) doItem(item map[string]interface{}) [][]interface{} {
	items := ia.mode.Do(item)
	if items == nil {
		return nil
	}

	if len(ia.affixFields) > 0 {
		for index := range items {
			for _, field := range ia.affixFields {
				items[index] = append(items[index], item[field])
			}
		}
	}

	return items
}

func (ia *InsertAction) doSource() {
	for {
		sourceItems, ok := ia.source.Receive()
		if ok == false {
			break
		}

		targetItems := make([][]interface{}, 0)

		for _, sourceItem := range sourceItems {
			items := ia.doItem(sourceItem)
			if items == nil {
				continue
			}

			targetItems = append(targetItems, items...)
		}

		ia.target.Send(targetItems)
	}

	ia.wg.Done()
}
