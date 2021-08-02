package action

import (
	"fmt"

	goEtl "github.com/auho/go-etl"
	"github.com/auho/go-etl/storage/database"
	"github.com/auho/go-simple-db/simple"
)

type TransferAction struct {
	action
	target         *database.DbTargetSlice
	targetDataName string
	fields         []string
	keys           []string
	fixedValues    []interface{}
}

func NewTransferAction(db simple.Driver, config goEtl.DbConfig, targetDataName string, alias map[string]string, fixedData map[string]interface{}) *TransferAction {
	ta := &TransferAction{}
	ta.targetDataName = targetDataName

	ta.init()

	if len(alias) >= 0 {
		for k, v := range alias {
			ta.fields = append(ta.fields, k)
			ta.keys = append(ta.keys, v)
		}
	} else {
		var err error
		ta.fields, err = db.GetTableColumns(targetDataName)
		if err != nil {
			panic(err)
		}

		for _, field := range ta.fields {
			ta.keys = append(ta.keys, field)
		}
	}

	if len(fixedData) > 0 {
		for k, v := range fixedData {
			ta.keys = append(ta.keys, k)
			ta.fixedValues = append(ta.fixedValues, v)
		}
	}

	targetConfig := ta.targetConfig(config, targetDataName)
	ta.target = database.NewDbTargetInsertSliceSlice(targetConfig, ta.keys, database.WithDbTargetSliceTruncate())

	return ta
}

func (ta *TransferAction) Start() {
	ta.target.Start()

	for i := 0; i < ta.concurrent; i++ {
		ta.wg.Add(1)
		go ta.doSource()
	}
}

func (ta *TransferAction) Done() {
	if ta.isDone {
		return
	}

	ta.isDone = true

	close(ta.itemsChan)
}

func (ta *TransferAction) Close() {
	ta.wg.Wait()

	ta.target.Done()
	ta.target.Close()
}

func (ta *TransferAction) GetFields() []string {
	return ta.fields
}

func (ta *TransferAction) Receive(items []map[string]interface{}) {
	ta.itemsChan <- items
}

func (ta *TransferAction) GetStatus() string {
	return ta.target.State.GetStatus()
}

func (ta *TransferAction) GetTitle() string {
	return fmt.Sprintf("Transfer[%s]", ta.targetDataName)
}

func (ta *TransferAction) doSource() {
	for {
		sourceItems, ok := <-ta.itemsChan
		if ok == false {
			break
		}

		targetItems := make([][]interface{}, 0)

		for _, sourceItem := range sourceItems {
			result := make([]interface{}, len(ta.fields), len(ta.fields))
			for k, field := range ta.fields {
				result[k] = sourceItem[field]
			}

			if len(ta.fixedValues) > 0 {
				result = append(result, ta.fixedValues...)
			}

			targetItems = append(targetItems, result)
		}

		ta.target.Send(targetItems)
	}

	ta.wg.Done()
}
