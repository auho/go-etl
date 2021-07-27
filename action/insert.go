package action

import (
	"sync"

	goEtl "github.com/auho/go-etl"
	"github.com/auho/go-etl/mode"
	"github.com/auho/go-etl/storage"
)

type InsertAction struct {
	concurrent   int
	target       *storage.DbTargetInsertInterface
	mode         mode.InsertMode
	idName       string
	dataName     string
	tagTableName string
	affixFields  []string
	isSourceEnd  bool
	wg           sync.WaitGroup
	itemsChan    chan []map[string]interface{}
}

func NewInsertAction(config goEtl.DbConfig, tagTableName string, moder mode.InsertMode, affixFields []string) *InsertAction {
	ia := &InsertAction{}
	ia.mode = moder
	ia.affixFields = affixFields
	ia.concurrent = 4
	ia.itemsChan = make(chan []map[string]interface{}, ia.concurrent)
	ia.tagTableName = tagTableName

	targetConfig := storage.NewDbTargetConfig()
	targetConfig.MaxConcurrent = 4
	targetConfig.Size = 2000
	targetConfig.Driver = config.Driver
	targetConfig.Dsn = config.Dsn
	targetConfig.Table = ia.tagTableName

	ia.target = storage.NewDbTargetInsertInterface(targetConfig)
	ia.target.SetFields(ia.GetKeys())
	ia.target.Start()

	return ia
}

func (ia *InsertAction) Start() {
	for i := 0; i < ia.concurrent; i++ {
		ia.wg.Add(1)
		go ia.doSource()
	}
}

func (ia *InsertAction) Close() {
	ia.wg.Wait()

	ia.target.Done()
	ia.target.Close()
}

func (ia *InsertAction) SourceDone() {
	if ia.isSourceEnd {
		return
	}

	ia.isSourceEnd = true

	close(ia.itemsChan)
}

func (ia *InsertAction) GetFields() []string {
	return append(ia.mode.GetFields(), ia.affixFields...)
}

func (ia *InsertAction) GetKeys() []string {
	return append(ia.mode.GetKeys(), ia.affixFields...)
}

func (ia *InsertAction) Receive(items []map[string]interface{}) {
	ia.itemsChan <- items
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
		sourceItems, ok := <-ia.itemsChan
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
