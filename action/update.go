package action

import (
	"sync"

	goEtl "github.com/auho/go-etl"
	"github.com/auho/go-etl/mode"
	"github.com/auho/go-etl/storage/database"
)

type UpdateAction struct {
	concurrent int
	dataName   string
	idName     string
	isDone     bool
	itemsChan  chan []map[string]interface{}
	modes      []mode.UpdateModer
	target     *database.DbTargetMap
	wg         sync.WaitGroup
}

func NewUpdateAction(config goEtl.DbConfig, dataName string, idName string, modes []mode.UpdateModer) *UpdateAction {
	ua := &UpdateAction{}
	ua.concurrent = 4
	ua.dataName = dataName
	ua.idName = idName
	ua.modes = modes
	ua.itemsChan = make(chan []map[string]interface{}, ua.concurrent)

	targetConfig := database.NewDbTargetConfig()
	targetConfig.MaxConcurrent = 4
	targetConfig.Size = 2000
	targetConfig.Driver = config.Driver
	targetConfig.Dsn = config.Dsn
	targetConfig.Table = ua.dataName

	ua.target = database.NewDbTargetUpdateSliceMap(targetConfig, idName)

	return ua
}

func (ua *UpdateAction) Start() {
	ua.target.Start()

	for i := 0; i < ua.concurrent; i++ {
		ua.wg.Add(1)
		go ua.doSource()
	}
}

func (ua *UpdateAction) Done() {
	if ua.isDone {
		return
	}

	ua.isDone = true

	close(ua.itemsChan)
}

func (ua *UpdateAction) Close() {
	ua.wg.Wait()

	for _, m := range ua.modes {
		m.Close()
	}

	ua.target.Done()
	ua.target.Close()
}

func (ua *UpdateAction) GetFields() []string {
	fields := make([]string, 0)
	for _, m := range ua.modes {
		fields = append(fields, m.GetFields()...)
	}

	fields = goEtl.RemoveReplicaSliceString(fields)

	return append(fields, ua.idName)
}

func (ua *UpdateAction) Receive(items []map[string]interface{}) {
	ua.itemsChan <- items
}

func (ua *UpdateAction) GetStatus() string {
	return ua.target.State.GetRealTimeStatus()
}

func (ua *UpdateAction) doSource() {
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
