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
	State  *DbTargetState

	target func()
	down   func()
}

func (t *dbTarget) Start() {
	t.State.status = "target start"

	for i := 0; i < t.maxConcurrent; i++ {
		t.wg.Add(1)

		go func() {
			t.target()

			t.wg.Done()
		}()
	}
}

func (t *dbTarget) Done() {
	if t.isDone {
		return
	}

	t.isDone = true

	t.down()
}

func (t *dbTarget) Close() {
	t.wg.Wait()

	t.db.Close()

	t.State.status = t.State.DoneStatus()
}

func (t *dbTarget) initDb(config *DbTargetConfig) {
	t.maxConcurrent = config.MaxConcurrent
	t.size = config.Size
	t.table = config.Table

	config.check()

	t.State = newDbTargetState()
	t.State.size = t.size
	t.State.maxConcurrent = t.maxConcurrent
	t.State.title = fmt.Sprintf("target[%s]", t.table)

	var err error
	t.db, err = simple.NewDriver(config.Driver, config.Dsn)
	if err != nil {
		panic(err)
	}

	err = t.db.Ping()
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
			stateDuration := uintptr(unsafe.Pointer(&t.State.duration))

			atomic.AddUintptr(&stateDuration, uintptr(endTime.Sub(startTime)))
			atomic.AddInt64(&t.State.itemAmount, int64(itemsAmount))

			t.State.status = fmt.Sprintf("target item amount[%d]", t.State.itemAmount)

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
			stateDuration := uintptr(unsafe.Pointer(&t.State.duration))

			atomic.AddUintptr(&stateDuration, uintptr(endTime.Sub(startTime)))
			atomic.AddInt64(&t.State.itemAmount, int64(itemsAmount))

			t.State.status = fmt.Sprintf("target item amount[%d]", t.State.itemAmount)

		} else {
			break
		}
	}
}
