package action

import (
	"fmt"
	"strings"

	goEtl "github.com/auho/go-etl"
	"github.com/auho/go-etl/mode"
	"github.com/auho/go-etl/storage/database"
	"github.com/auho/go-simple-db/simple"
)

type CleanAction struct {
	action
	target         *database.DbTargetSlice
	modes          []mode.UpdateModer
	targetDataName string
	fields         []string
}

func NewCleanAction(db simple.Driver, config goEtl.DbConfig, targetDataName string, modes []mode.UpdateModer) *CleanAction {
	ca := &CleanAction{}
	ca.modes = modes
	ca.targetDataName = targetDataName

	ca.init()

	var err error
	ca.fields, err = db.GetTableColumns(ca.targetDataName)
	if err != nil {
		panic(err)
	}

	targetConfig := ca.targetConfig(config, targetDataName)
	ca.target = database.NewDbTargetInsertSliceSlice(targetConfig, ca.fields, database.WithDbTargetSliceTruncate())

	return ca
}

func (ca *CleanAction) Start() {
	ca.target.Start()

	for i := 0; i < ca.concurrent; i++ {
		ca.wg.Add(1)
		go ca.doSource()
	}
}

func (ca *CleanAction) Done() {
	if ca.isDone {
		return
	}

	ca.isDone = true

	close(ca.itemsChan)
}

func (ca *CleanAction) Close() {
	ca.wg.Wait()

	for _, m := range ca.modes {
		m.Close()
	}

	ca.target.Done()
	ca.target.Close()
}

func (ca *CleanAction) GetFields() []string {
	return ca.fields
}

func (ca *CleanAction) Receive(items []map[string]interface{}) {
	ca.itemsChan <- items
}

func (ca *CleanAction) GetStatus() string {
	return ca.target.State.GetStatus()
}

func (ca *CleanAction) GetTitle() string {
	s := make([]string, 0)
	for _, m := range ca.modes {
		s = append(s, m.GetTitle())
	}

	return fmt.Sprintf("Clean[%s] {%s}", ca.targetDataName, strings.Join(s, ", "))
}

func (ca *CleanAction) doSource() {
	for {
		sourceItems, ok := <-ca.itemsChan
		if ok == false {
			break
		}

		targetItems := make([][]interface{}, 0)

		for _, sourceItem := range sourceItems {
			isClean := false
			for _, m := range ca.modes {
				res := m.Do(sourceItem)
				if res != nil && len(res) > 0 {
					isClean = true
					break
				}
			}

			if isClean == true {
				continue
			}

			item := make([]interface{}, 0, len(ca.fields))
			for _, field := range ca.fields {
				item = append(item, sourceItem[field])
			}

			targetItems = append(targetItems, item)
		}

		ca.target.Send(targetItems)
	}

	ca.wg.Done()
}
