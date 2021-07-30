package database

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	go_etl "github.com/auho/go-etl"
	"github.com/auho/go-simple-db/simple"
)

type dbTarget struct {
	maxConcurrent int
	size          int
	table         string

	isDone bool
	wg     sync.WaitGroup
	db     simple.Driver
	state  *DbTargetState

	target func()
	down   func()
}

func (dt *dbTarget) Start() {
	fmt.Println("target start")

	for i := 0; i < dt.maxConcurrent; i++ {
		dt.wg.Add(1)

		go func() {
			dt.target()

			dt.wg.Done()
		}()
	}
}

func (dt *dbTarget) Done() {
	if dt.isDone {
		return
	}

	dt.isDone = true

	dt.down()
}

func (dt *dbTarget) Close() {
	dt.wg.Wait()

	dt.db.Close()

	fmt.Println("\ntarget done!")
}

func (dt *dbTarget) State() {
	fmt.Println(fmt.Sprintf("Max Concurrent: %d \nSize: %d\nAmount: %d", dt.state.maxConcurrent, dt.state.size, dt.state.itemAmount))
}

func (dt *dbTarget) initDb(config *DbTargetConfig) {
	dt.maxConcurrent = config.MaxConcurrent
	dt.size = config.Size
	dt.table = config.Table

	config.check()

	dt.state = newDbTargetState()
	dt.state.size = dt.size
	dt.state.maxConcurrent = dt.maxConcurrent

	var err error
	dt.db, err = simple.NewDriver(config.Driver, config.Dsn)
	if err != nil {
		panic(err)
	}

	err = dt.db.Ping()
	if err != nil {
		panic(err)
	}
}

func WithDbTargetSliceTruncate() func(*DbTargetSlice) error {
	return func(d *DbTargetSlice) error {
		return d.db.Truncate(d.table)
	}
}

type DbTargetSlicePrepareFunc func(d *DbTargetSlice) error

type DbTargetSlice struct {
	dbTarget
	fields    []string
	itemsChan chan [][]interface{}
	sliceFunc func(t *DbTargetSlice, items [][]interface{}) error
}

func newDbTargetSlice(config *DbTargetConfig, prepareFuncs ...DbTargetSlicePrepareFunc) *DbTargetSlice {
	t := &DbTargetSlice{}
	t.itemsChan = make(chan [][]interface{}, t.maxConcurrent)

	t.target = t.doTarget
	t.down = t.doDown

	t.initDb(config)

	for _, pf := range prepareFuncs {
		err := pf(t)
		if err != nil {
			panic(err)
		}
	}

	return t
}

func (t *DbTargetSlice) Send(items [][]interface{}) {
	t.itemsChan <- items
}

func (t *DbTargetSlice) doDown() {
	close(t.itemsChan)
}

func (t *DbTargetSlice) doTarget() {
	var startTime time.Time
	var endTime time.Time

	for {
		if items, ok := <-t.itemsChan; ok {
			if len(items) <= 0 {
				continue
			}

			startTime = time.Now()

			var insertItems [][]interface{}
			itemsAmount := len(items)

			for start := 0; start < itemsAmount; start += t.size {
				end := start + t.size
				if end >= itemsAmount {
					insertItems = items[start:]
				} else {
					insertItems = items[start:end]
				}

				err := t.sliceFunc(t, insertItems)
				if err != nil {
					panic(err)
				}
			}

			endTime = time.Now()
			stateDuration := uintptr(unsafe.Pointer(&t.state.duration))

			atomic.AddUintptr(&stateDuration, uintptr(endTime.Sub(startTime)))
			atomic.AddInt64(&t.state.itemAmount, int64(itemsAmount))

			go_etl.PrintColor(fmt.Sprintf("target item amount:: %d", t.state.itemAmount))

		} else {
			break
		}
	}
}

func WithDbTargetMapTruncate() func(*DbTargetMap) error {
	return func(d *DbTargetMap) error {
		return d.db.Truncate(d.table)
	}
}

type DbTargetMapPrepareFunc func(d *DbTargetMap) error

type DbTargetMap struct {
	dbTarget
	itemsChan chan []map[string]interface{}
	mapFunc   func(t *DbTargetMap, items []map[string]interface{}) error
}

func newDbTargetMap(config *DbTargetConfig, prepareFuncs ...DbTargetMapPrepareFunc) *DbTargetMap {
	t := &DbTargetMap{}
	t.itemsChan = make(chan []map[string]interface{}, t.maxConcurrent)

	t.target = t.doTarget
	t.down = t.doDown

	t.initDb(config)

	for _, pf := range prepareFuncs {
		err := pf(t)
		if err != nil {
			panic(err)
		}
	}

	return t
}

func (t *DbTargetMap) Send(items []map[string]interface{}) {
	t.itemsChan <- items
}

func (t *DbTargetMap) doDown() {
	close(t.itemsChan)
}

func (t *DbTargetMap) doTarget() {
	var startTime time.Time
	var endTime time.Time

	for {
		if items, ok := <-t.itemsChan; ok {
			if len(items) <= 0 {
				continue
			}

			startTime = time.Now()

			var insertItems []map[string]interface{}
			itemsAmount := len(items)

			for start := 0; start < itemsAmount; start += t.size {
				end := start + t.size
				if end >= itemsAmount {
					insertItems = items[start:]
				} else {
					insertItems = items[start:end]
				}

				err := t.mapFunc(t, insertItems)
				if err != nil {
					panic(err)
				}
			}

			endTime = time.Now()
			stateDuration := uintptr(unsafe.Pointer(&t.state.duration))

			atomic.AddUintptr(&stateDuration, uintptr(endTime.Sub(startTime)))
			atomic.AddInt64(&t.state.itemAmount, int64(itemsAmount))

			go_etl.PrintColor(fmt.Sprintf("target item amount:: %d", t.state.itemAmount))

		} else {
			break
		}
	}
}
