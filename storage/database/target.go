package database

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

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
}

func (dt *dbTarget) State() {
	fmt.Println(fmt.Sprintf("Max Concurrent: %d \nSize: %d\nAmount: %d", dt.state.maxConcurrent, dt.state.size, dt.state.itemAmount))
}

func (dt *dbTarget) doConfig(config *DbTargetConfig) {
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

type DbTargetSlice struct {
	dbTarget
	fields    []string
	itemsChan chan [][]interface{}
	sliceFunc func(t *DbTargetSlice, items [][]interface{}) error
}

func newDbTargetSlice(config *DbTargetConfig) *DbTargetSlice {
	t := &DbTargetSlice{}
	t.itemsChan = make(chan [][]interface{}, t.maxConcurrent)

	t.target = t.doTarget
	t.down = t.doDown

	t.doConfig(config)

	return t
}

func (d *DbTargetSlice) Send(items [][]interface{}) {
	d.itemsChan <- items
}

func (d *DbTargetSlice) doDown() {
	close(d.itemsChan)
}

func (d *DbTargetSlice) doTarget() {
	var startTime time.Time
	var endTime time.Time

	for {
		if items, ok := <-d.itemsChan; ok {
			if len(items) <= 0 {
				continue
			}

			startTime = time.Now()

			var insertItems [][]interface{}
			itemsAmount := len(items)

			for start := 0; start < itemsAmount; start += d.size {
				end := start + d.size
				if end >= itemsAmount {
					insertItems = items[start:]
				} else {
					insertItems = items[start:end]
				}

				err := d.sliceFunc(d, insertItems)
				if err != nil {
					panic(err)
				}
			}

			endTime = time.Now()
			stateDuration := uintptr(unsafe.Pointer(&d.state.duration))

			atomic.AddUintptr(&stateDuration, uintptr(endTime.Sub(startTime)))
			atomic.AddInt64(&d.state.itemAmount, int64(itemsAmount))
		} else {
			break
		}
	}
}

type DbTargetMap struct {
	dbTarget
	itemsChan chan []map[string]interface{}
	mapFunc   func(t *DbTargetMap, items []map[string]interface{}) error
}

func newDbTargetMap(config *DbTargetConfig) *DbTargetMap {
	d := &DbTargetMap{}
	d.itemsChan = make(chan []map[string]interface{}, d.maxConcurrent)

	d.target = d.doTarget
	d.down = d.doDown

	d.doConfig(config)

	return d
}

func (d *DbTargetMap) Send(items []map[string]interface{}) {
	d.itemsChan <- items
}

func (d *DbTargetMap) doDown() {
	close(d.itemsChan)
}

func (d *DbTargetMap) doTarget() {
	var startTime time.Time
	var endTime time.Time

	for {
		if items, ok := <-d.itemsChan; ok {
			if len(items) <= 0 {
				continue
			}

			startTime = time.Now()

			var insertItems []map[string]interface{}
			itemsAmount := len(items)

			for start := 0; start < itemsAmount; start += d.size {
				end := start + d.size
				if end >= itemsAmount {
					insertItems = items[start:]
				} else {
					insertItems = items[start:end]
				}

				err := d.mapFunc(d, insertItems)
				if err != nil {
					panic(err)
				}
			}

			endTime = time.Now()
			stateDuration := uintptr(unsafe.Pointer(&d.state.duration))

			atomic.AddUintptr(&stateDuration, uintptr(endTime.Sub(startTime)))
			atomic.AddInt64(&d.state.itemAmount, int64(itemsAmount))

		} else {
			break
		}
	}
}
