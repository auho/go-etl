package database

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/auho/go-simple-db/simple"
)

type DbSource struct {
	maxConcurrent int
	size          int
	page          int
	table         string
	pKeyName      string
	fields        []string

	retryPageMap sync.Map
	wg           sync.WaitGroup
	pageChan     chan interface{}
	itemsChan    chan []map[string]interface{}
	State        *DbSourceState
	lastPagePk   interface{}

	db simple.Driver
}

func NewDbSource(config *DbSourceConfig) *DbSource {
	ds := &DbSource{}
	ds.doConfig(config)

	return ds
}

func (ds *DbSource) Start() {
	ds.State.realtimeStatus = "source start"

	ds.pageChan = make(chan interface{}, ds.maxConcurrent)
	ds.itemsChan = make(chan []map[string]interface{}, ds.maxConcurrent)

	rowsQuery := ds.generateRowsQuery(ds.size)
	for i := 0; i < ds.maxConcurrent; i++ {
		ds.wg.Add(1)

		go ds.source(rowsQuery)
	}

	go ds.sourcePage()

	go ds.Close()
}

func (ds *DbSource) Consume() ([]map[string]interface{}, bool) {
	items, ok := <-ds.itemsChan

	return items, ok
}

func (ds *DbSource) Close() {
	ds.wg.Wait()

	close(ds.itemsChan)
	ds.db.Close()

	ds.State.realtimeStatus = ds.State.DoneStatus()
}

func (ds *DbSource) doConfig(config *DbSourceConfig) {
	ds.maxConcurrent = config.MaxConcurrent
	ds.size = config.Size
	ds.page = config.Page
	ds.table = config.Table
	ds.pKeyName = config.PKeyName
	ds.fields = config.Fields

	config.check()

	ds.State = newDbSourceState()
	ds.State.size = ds.size
	ds.State.maxConcurrent = ds.maxConcurrent

	var err error
	ds.db, err = simple.NewDriver(config.Driver, config.Dsn)
	if err != nil {
		panic(err)
	}

	err = ds.db.Ping()
	if err != nil {
		panic(err)
	}
}

func (ds *DbSource) sourcePage() {
	amount := ds.getTableAmount()
	if amount <= 0 {
		panic("db source table amount is  error 0")
	}

	maxPage := int(math.Ceil(float64(amount) / float64(ds.size)))
	if ds.page > 0 && ds.page <= maxPage {
		maxPage = ds.page
	}

	if maxPage <= 0 {
		panic(fmt.Sprintf("max Page is error %d", maxPage))
	}

	ds.State.page = maxPage

	nextPkQuery := ds.generateNextPkQuery(ds.size)
	go func() {
		minPk := ds.getMinPk()
		ds.pageChan <- minPk

		prePk := minPk
		for page := 1; page < maxPage; page++ {
			nextPk := ds.getNextPk(prePk, nextPkQuery)
			if nextPk == nil {
				break
			}

			ds.pageChan <- nextPk
			prePk = nextPk
			ds.lastPagePk = nextPk
		}

		close(ds.pageChan)
	}()
}

func (ds *DbSource) source(query string) {
	for {
		pk, ok := <-ds.pageChan
		if ok == false {
			break
		}

		rows, err := ds.db.QueryInterface(query, pk)
		if err != nil {
			ds.retryPage(pk)
		}

		ds.retryPageMap.Delete(pk)

		if rows == nil {
			continue
		}

		ds.itemsChan <- rows

		atomic.AddInt64(&ds.State.itemAmount, int64(len(rows)))

		ds.State.realtimeStatus = fmt.Sprintf("source last page pk:: %v; item amount:: %d", ds.lastPagePk, ds.State.itemAmount)
	}

	ds.wg.Done()
}

func (ds *DbSource) retryPage(pk interface{}) {
	if n, ok := ds.retryPageMap.Load(pk); ok {
		num := n.(int64)
		if num >= 2 {

		} else {
			ds.retryPageMap.Store(pk, num+1)
		}
	} else {
		ds.retryPageMap.Store(pk, 1)
	}
}

func (ds *DbSource) getMinPk() interface{} {
	query := fmt.Sprintf("SELECT `%s` FROM `%s` ORDER BY `%s` ASC LIMIT 0, ?", ds.pKeyName, ds.table, ds.pKeyName)
	minPk, err := ds.db.QueryFieldInterface(ds.pKeyName, query, 1)
	if err != nil {
		panic(err)
	}

	return minPk
}

func (ds *DbSource) getTableAmount() int64 {
	query := fmt.Sprintf("SELECT COUNT(*) AS 'amount' FROM `%s`", ds.table)
	res, err := ds.db.QueryFieldInterface("amount", query)
	if err != nil {
		panic(err)
	}

	amount, err := strconv.ParseInt(string(res.([]uint8)), 10, 64)
	if err != nil {
		panic(err)
	}

	return amount
}

func (ds *DbSource) getNextPk(pk interface{}, query string) interface{} {
	nextPk, err := ds.db.QueryFieldInterface(ds.pKeyName, query, pk)
	if err != nil {
		panic(err)
	} else {
		return nextPk
	}
}

func (ds *DbSource) generateNextPkQuery(size int) string {
	return fmt.Sprintf("SELECT `%s` FROM `%s` WHERE `%s` >= ? ORDER BY `%s` ASC LIMIT %d, %d", ds.pKeyName, ds.table, ds.pKeyName, ds.pKeyName, size, 1)
}

func (ds *DbSource) generateRowsQuery(size int) string {
	fields := ""
	if len(ds.fields) <= 0 {
		fields = "*"
	} else {
		fields = "`" + strings.Join(ds.fields, "`, `") + "`"
	}

	return fmt.Sprintf("SELECT %s FROM `%s` WHERE `%s` >= ? ORDER BY `%s` ASC LIMIT %d, %d", fields, ds.table, ds.pKeyName, ds.pKeyName, 0, size)
}
