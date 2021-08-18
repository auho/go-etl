package action

import (
	"runtime"
	"sync"

	go_etl "github.com/auho/go-etl"
	"github.com/auho/go-etl/storage/database"
)

type Actionor interface {
	Start()
	Done()
	Close()
	GetFields() []string
	Receive([]map[string]interface{})
	GetStatus() string
	GetTitle() string
}

type action struct {
	wg         sync.WaitGroup
	isDone     bool
	concurrent int
	itemsChan  chan []map[string]interface{}
}

func (a *action) init() {
	a.concurrent = runtime.NumCPU()
	a.itemsChan = make(chan []map[string]interface{})
}

func (a *action) targetConfig(config go_etl.DbConfig, targetTableName string) *database.DbTargetConfig {
	targetConfig := database.NewDbTargetConfig()
	targetConfig.MaxConcurrent = a.concurrent
	targetConfig.Size = 2000
	targetConfig.Driver = config.Driver
	targetConfig.Dsn = config.Dsn
	targetConfig.Table = targetTableName

	return targetConfig
}
