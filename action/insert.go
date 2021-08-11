package action

import (
	"fmt"

	goEtl "github.com/auho/go-etl"
	"github.com/auho/go-etl/mode"
	"github.com/auho/go-etl/storage/database"
)

type Insert struct {
	action
	target       *database.DbTargetSlice
	mode         mode.InsertModer
	tagTableName string
	affixFields  []string
}

func NewInsert(config goEtl.DbConfig, tagTableName string, moder mode.InsertModer, affixFields []string) *Insert {
	ia := &Insert{}
	ia.tagTableName = tagTableName
	ia.affixFields = affixFields
	ia.mode = moder

	ia.init()

	targetConfig := ia.targetConfig(config, ia.tagTableName)
	ia.target = database.NewDbTargetInsertSliceSlice(targetConfig, ia.getKeys(), database.WithDbTargetSliceTruncate())

	return ia
}

func (ia *Insert) Start() {
	ia.target.Start()

	for i := 0; i < ia.concurrent; i++ {
		ia.wg.Add(1)
		go ia.doSource()
	}
}

func (ia *Insert) Done() {
	if ia.isDone {
		return
	}

	ia.isDone = true

	close(ia.itemsChan)
}

func (ia *Insert) Close() {
	ia.wg.Wait()

	ia.mode.Close()

	ia.target.Done()
	ia.target.Close()
}

func (ia *Insert) GetFields() []string {
	return append(ia.mode.GetFields(), ia.affixFields...)
}

func (ia *Insert) Receive(items []map[string]interface{}) {
	ia.itemsChan <- items
}

func (ia *Insert) GetStatus() string {
	return ia.target.State.GetStatus()
}

func (ia *Insert) GetTitle() string {
	return fmt.Sprintf("Insert[%s] {%s}", ia.tagTableName, ia.mode.GetTitle())
}

func (ia *Insert) getKeys() []string {
	return append(ia.mode.GetKeys(), ia.affixFields...)
}

func (ia *Insert) doSource() {
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

func (ia *Insert) doItem(item map[string]interface{}) [][]interface{} {
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
