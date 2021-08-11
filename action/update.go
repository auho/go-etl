package action

import (
	"fmt"
	"strings"

	goEtl "github.com/auho/go-etl"
	"github.com/auho/go-etl/mode"
	"github.com/auho/go-etl/storage/database"
)

type Update struct {
	action
	target   *database.DbTargetMap
	modes    []mode.UpdateModer
	idName   string
	dataName string
}

func NewUpdate(config goEtl.DbConfig, dataName string, idName string, modes []mode.UpdateModer) *Update {
	ua := &Update{}
	ua.dataName = dataName
	ua.idName = idName
	ua.modes = modes

	ua.init()

	targetConfig := ua.targetConfig(config, ua.dataName)
	ua.target = database.NewDbTargetUpdateSliceMap(targetConfig, idName)

	return ua
}

func (ua *Update) Start() {
	ua.target.Start()

	for i := 0; i < ua.concurrent; i++ {
		ua.wg.Add(1)
		go ua.doSource()
	}
}

func (ua *Update) Done() {
	if ua.isDone {
		return
	}

	ua.isDone = true

	close(ua.itemsChan)
}

func (ua *Update) Close() {
	ua.wg.Wait()

	for _, m := range ua.modes {
		m.Close()
	}

	ua.target.Done()
	ua.target.Close()
}

func (ua *Update) GetFields() []string {
	fields := make([]string, 0)
	for _, m := range ua.modes {
		fields = append(fields, m.GetFields()...)
	}

	fields = goEtl.RemoveReplicaSliceString(fields)

	return append(fields, ua.idName)
}

func (ua *Update) Receive(items []map[string]interface{}) {
	ua.itemsChan <- items
}

func (ua *Update) GetStatus() string {
	return ua.target.State.GetStatus()
}

func (ua *Update) GetTitle() string {
	s := make([]string, 0)
	for _, m := range ua.modes {
		s = append(s, m.GetTitle())
	}

	return fmt.Sprintf("Update[%s] {%s}", ua.dataName, strings.Join(s, ", "))
}

func (ua *Update) doSource() {
	for {
		sourceItems, ok := <-ua.itemsChan
		if ok == false {
			break
		}

		targetItems := make([]map[string]interface{}, 0)

		for _, sourceItem := range sourceItems {
			item := make(map[string]interface{})
			for _, m := range ua.modes {
				res := m.Do(sourceItem)
				for k, v := range res {
					item[k] = v
				}
			}

			if len(item) <= 0 {
				continue
			}

			item[ua.idName] = sourceItem[ua.idName]
			targetItems = append(targetItems, item)
		}

		ua.target.Send(targetItems)
	}

	ua.wg.Done()
}
