package action

import (
	"fmt"
	"sync"

	goEtl "github.com/auho/go-etl"
	"github.com/auho/go-etl/mode"
	"github.com/auho/go-etl/storage/database"
)

type InsertAction struct {
	concurrent   int
	tagTableName string
	affixFields  []string
	isDone       bool
	itemsChan    chan []map[string]interface{}
	mode         mode.InsertModer
	target       *database.DbTargetSlice
	wg           sync.WaitGroup
}

func NewInsertAction(config goEtl.DbConfig, tagTableName string, moder mode.InsertModer, affixFields []string) *InsertAction {
	ia := &InsertAction{}
	ia.concurrent = 4
	ia.tagTableName = tagTableName
	ia.affixFields = affixFields
	ia.mode = moder
	ia.itemsChan = make(chan []map[string]interface{}, ia.concurrent)

	targetConfig := database.NewDbTargetConfig()
	targetConfig.MaxConcurrent = 4
	targetConfig.Size = 2000
	targetConfig.Driver = config.Driver
	targetConfig.Dsn = config.Dsn
	targetConfig.Table = ia.tagTableName

	ia.target = database.NewDbTargetInsertSliceSlice(targetConfig, ia.getKeys(), database.WithDbTargetSliceTruncate())

	return ia
}

func (ia *InsertAction) Start() {
	ia.target.Start()

	for i := 0; i < ia.concurrent; i++ {
		ia.wg.Add(1)
		go ia.doSource()
	}
}

func (ia *InsertAction) Done() {
	if ia.isDone {
		return
	}

	ia.isDone = true

	close(ia.itemsChan)
}

func (ia *InsertAction) Close() {
	ia.wg.Wait()

	ia.mode.Close()

	ia.target.Done()
	ia.target.Close()
}

func (ia *InsertAction) GetFields() []string {
	return append(ia.mode.GetFields(), ia.affixFields...)
}

func (ia *InsertAction) Receive(items []map[string]interface{}) {
	ia.itemsChan <- items
}

func (ia *InsertAction) GetStatus() string {
	return ia.target.State.GetStatus()
}

func (ia *InsertAction) GetTitle() string {
	return fmt.Sprintf("Insert[%s] {%s}", ia.tagTableName, ia.mode.GetTitle())
}

func (ia *InsertAction) getKeys() []string {
	return append(ia.mode.GetKeys(), ia.affixFields...)
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
