package storage

import (
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/auho/go-simple-db/simple"
)

type Source interface {
}

type Target interface {
}

type DbSourceConfig struct {
	MaxConcurrent int
	Size          int
	Page          int
	Driver        string
	Dsn           string
	Scheme        string
	Table         string
	PKeyName      string
	Fields        []string
}

type DbTargetConfig struct {
	MaxConcurrent int
	Size          int
	Driver        string
	Dsn           string
	Scheme        string
	Table         string
}

type DbSource struct {
	maxConcurrent int
	size          int
	page          int
	scheme        string
	table         string
	pKeyName      string
	fields        []string

	retryPageMap sync.Map
	wg           sync.WaitGroup
	pageChan     chan interface{}
	itemsChan    chan []map[string]interface{}

	db simple.Driver
}

func NewDbSource(dsConfig DbSourceConfig) *DbSource {
	ds := &DbSource{}
	ds.maxConcurrent = dsConfig.MaxConcurrent
	ds.size = dsConfig.Size
	ds.page = dsConfig.Page
	ds.scheme = dsConfig.Scheme
	ds.pKeyName = dsConfig.PKeyName
	ds.fields = dsConfig.Fields

	var err error
	ds.db, err = simple.NewDriver(dsConfig.Driver, dsConfig.Dsn)
	if err != nil {
		panic(err)
	}

	ds.pageChan = make(chan interface{}, ds.maxConcurrent)
	ds.itemsChan = make(chan []map[string]interface{}, ds.maxConcurrent)

	return ds
}

func (ds *DbSource) Start() {
	maxPage := 0
	amount := ds.getTableAmount()
	if ds.page > 0 {
		maxPage = ds.page
	} else {
		maxPage = int(math.Ceil(float64(amount) / float64(ds.size)))
	}

	rowsQuery := ds.generateRowsQuery(ds.size)
	nextPkQuery := ds.generateNextPkQuery(ds.size)

	for i := 0; i < ds.maxConcurrent; i++ {
		ds.wg.Add(1)

		go ds.rows(rowsQuery)
	}

	go func() {
		minPk := ds.getMinPk()
		ds.pageChan <- minPk

		prePk := minPk
		for page := 1; page < maxPage; page++ {
			nextPk := ds.getNextPk(prePk, nextPkQuery)
			ds.pageChan <- nextPk
			prePk = nextPk
		}

		close(ds.pageChan)
	}()
}

func (ds *DbSource) Receive() ([]map[string]interface{}, bool) {
	items, ok := <-ds.itemsChan

	return items, ok
}

func (ds *DbSource) Close() {
	ds.wg.Wait()

	close(ds.itemsChan)
}

func (ds *DbSource) rows(query string) {
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
	query := fmt.Sprintf("SELECT `%s` FROM `%s`.`%s` ORDER BY `%s` ASC LIMIT 0, 1", ds.pKeyName, ds.scheme, ds.table, ds.pKeyName)
	minPk, err := ds.db.QueryFieldInterface(ds.pKeyName, query)
	if err != nil {
		panic(err)
	}

	return minPk
}

func (ds *DbSource) getTableAmount() int64 {
	query := fmt.Sprintf("SELECT COUNT(*) AS 'amount' FROM `%s`.`%s`", ds.scheme, ds.table)
	amount, err := ds.db.QueryFieldInterface("amount", query)
	if err != nil {
		panic(err)
	}

	return amount.(int64)
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
	return fmt.Sprintf("SELECT %s FROM `%s`.`%s` WHERE `%s` > ? ORDER BY `%s` ASC LIMIT %d, %d", ds.pKeyName, ds.scheme, ds.table, ds.pKeyName, ds.pKeyName, size, 1)

}

func (ds *DbSource) generateRowsQuery(size int) string {
	fields := strings.Join(ds.fields, "`, `")
	return fmt.Sprintf("SELECT `%s` FROM `%s`.`%s` WHERE `%s` >= ? ORDER BY `%s` ASC LIMIT %d, %d", fields, ds.scheme, ds.table, ds.pKeyName, ds.pKeyName, 0, size)
}

type DbTarget struct {
	maxConcurrent int
	size          int
	scheme        string
	table         string

	isDone bool
	wg     sync.WaitGroup
	db     simple.Driver

	target func()
	down   func()
}

func (dt *DbTarget) Start(config *DbTargetConfig) {
	dt.maxConcurrent = config.MaxConcurrent
	dt.size = config.Size
	dt.scheme = config.Scheme
	dt.table = config.Table

	var err error
	dt.db, err = simple.NewDriver(config.Driver, config.Dsn)
	if err != nil {
		panic(err)
	}

	for i := 0; i < dt.maxConcurrent; i++ {
		dt.wg.Add(1)

		go func() {
			dt.target()

			dt.wg.Done()
		}()
	}
}

func (dt *DbTarget) Done() {
	if dt.isDone {
		return
	}

	dt.isDone = true
	dt.down()
}

func (dt *DbTarget) Close() {
	dt.wg.Wait()
}

type DbTargetInsertInterface struct {
	DbTarget
	fields    []string
	itemsChan chan [][]interface{}
}

func NewDbTargetInsertInterface() *DbTargetInsertInterface {
	d := &DbTargetInsertInterface{}
	d.itemsChan = make(chan [][]interface{}, d.maxConcurrent)

	d.target = d.doTarget
	d.down = d.doDown

	return d
}

func (d *DbTargetInsertInterface) SetFields(fields []string) {
	d.fields = fields
}

func (d *DbTargetInsertInterface) Send(items [][]interface{}) {
	d.itemsChan <- items
}

func (d *DbTargetInsertInterface) doDown() {
	close(d.itemsChan)
}

func (d *DbTargetInsertInterface) doTarget() {
	for {
		if items, ok := <-d.itemsChan; ok {
			res, err := d.db.BulkInsertFromSliceSlice(d.table, d.fields, items)
			if err != nil {
				panic(err)
			}

			num, err := res.RowsAffected()
			if err != nil {
				panic(err)
			}

			if num != int64(len(items)) {
				panic(fmt.Sprintf("target affected is error [%d != %d]", num, len(items)))
			}
		} else {
			break
		}
	}
}

type DbTargetInsertMap struct {
	DbTarget
	itemsChan chan []map[string]interface{}
}

func NewDbTargetInsertMap() *DbTargetInsertMap {
	d := &DbTargetInsertMap{}
	d.itemsChan = make(chan []map[string]interface{}, d.maxConcurrent)

	d.target = d.doTarget
	d.down = d.doDown

	return d
}

func (d *DbTargetInsertMap) Send(items []map[string]interface{}) {
	d.itemsChan <- items
}

func (d *DbTargetInsertMap) doDown() {
	close(d.itemsChan)
}

func (d *DbTargetInsertMap) doTarget() {
	for {
		if items, ok := <-d.itemsChan; ok {
			res, err := d.db.BulkInsertFromSliceMap(d.table, items)
			if err != nil {
				panic(err)
			}

			num, err := res.RowsAffected()
			if err != nil {
				panic(err)
			}

			if num != int64(len(items)) {
				panic(fmt.Sprintf("target affected is error [%d != %d]", num, len(items)))
			}
		} else {
			break
		}
	}
}
