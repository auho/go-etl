package storage

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/auho/go-simple-db/simple"
)

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

func NewDbSourceConfig() *DbSourceConfig {
	c := &DbSourceConfig{}

	return c
}

func (sc *DbSourceConfig) check() {
	if sc.MaxConcurrent <= 0 {
		panic(fmt.Sprintf("db source config Max Concurrent is error %d", sc.MaxConcurrent))
	}

	if sc.Page <= 0 {
		panic(fmt.Sprintf("db source config Page is error %d", sc.Page))
	}

	if sc.Size <= 0 {
		panic(fmt.Sprintf("db source config Size is error %d", sc.Size))
	}

	if sc.PKeyName == "" {
		panic(fmt.Sprintf("db source config PKeyName is error %s", sc.PKeyName))
	}

	if sc.Scheme == "" {
		panic(fmt.Sprintf("db source config Scheme is error %s", sc.Scheme))
	}

	if sc.Table == "" {
		panic(fmt.Sprintf("db source config Table is error %s", sc.Table))
	}
}

type DbTargetConfig struct {
	MaxConcurrent int
	Size          int
	Driver        string
	Dsn           string
	Scheme        string
	Table         string
}

func NewDbTargetConfig() *DbTargetConfig {
	c := &DbTargetConfig{}

	return c
}

func (tc *DbTargetConfig) check() {
	if tc.MaxConcurrent <= 0 {
		panic(fmt.Sprintf("db target config Max Concurrent is error %d", tc.MaxConcurrent))
	}

	if tc.Size <= 0 {
		panic(fmt.Sprintf("db target config Size is error %d", tc.Size))
	}

	if tc.Table == "" {
		panic(fmt.Sprintf("db target config Table is error %s", tc.Table))
	}
}

type DbSourceState struct {
	duration      time.Duration
	maxConcurrent int
	page          int
	size          int
	itemAmount    int64
}

func newDbSourceState() *DbSourceState {
	s := &DbSourceState{}

	return s
}

type DbTargetState struct {
	duration      time.Duration
	maxConcurrent int
	size          int
	itemAmount    int64
}

func newDbTargetState() *DbTargetState {
	s := &DbTargetState{}

	return s
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

func NewDbSource() *DbSource {
	ds := &DbSource{}

	return ds
}

func (ds *DbSource) Start(config *DbSourceConfig) {
	ds.doStart(config)

	go ds.Close()
}

func (ds *DbSource) Receive() ([]map[string]interface{}, bool) {
	items, ok := <-ds.itemsChan

	return items, ok
}

func (ds *DbSource) Close() {
	ds.wg.Wait()

	close(ds.itemsChan)
	ds.db.Close()
}

func (ds *DbSource) doStart(config *DbSourceConfig) {
	ds.maxConcurrent = config.MaxConcurrent
	ds.size = config.Size
	ds.page = config.Page
	ds.scheme = config.Scheme
	ds.table = config.Table
	ds.pKeyName = config.PKeyName
	ds.fields = config.Fields

	config.check()

	var err error
	ds.db, err = simple.NewDriver(config.Driver, config.Dsn)
	if err != nil {
		panic(err)
	}

	err = ds.db.Ping()
	if err != nil {
		panic(err)
	}

	ds.pageChan = make(chan interface{}, ds.maxConcurrent)
	ds.itemsChan = make(chan []map[string]interface{}, ds.maxConcurrent)

	rowsQuery := ds.generateRowsQuery(ds.size)
	for i := 0; i < ds.maxConcurrent; i++ {
		ds.wg.Add(1)

		go ds.source(rowsQuery)
	}

	go ds.sourcePage()
}

func (ds *DbSource) sourcePage() {
	amount := ds.getTableAmount()
	if amount <= 0 {
		panic("db source table amount is  error 0")
	}

	maxPage := int(math.Ceil(float64(amount) / float64(ds.size)))
	if ds.page <= maxPage {
		maxPage = ds.page
	}

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
	query := fmt.Sprintf("SELECT `%s` FROM `%s`.`%s` ORDER BY `%s` ASC LIMIT 0, ?", ds.pKeyName, ds.scheme, ds.table, ds.pKeyName)
	minPk, err := ds.db.QueryFieldInterface(ds.pKeyName, query, 1)
	if err != nil {
		panic(err)
	}

	return minPk
}

func (ds *DbSource) getTableAmount() int64 {
	query := fmt.Sprintf("SELECT COUNT(*) AS 'amount' FROM `%s`.`%s`", ds.scheme, ds.table)
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
	return fmt.Sprintf("SELECT `%s` FROM `%s`.`%s` WHERE `%s` >= ? ORDER BY `%s` ASC LIMIT %d, %d", ds.pKeyName, ds.scheme, ds.table, ds.pKeyName, ds.pKeyName, size, 1)

}

func (ds *DbSource) generateRowsQuery(size int) string {
	fields := ""
	if len(ds.fields) <= 0 {
		fields = "*"
	} else {
		fields = "`" + strings.Join(ds.fields, "`, `") + "`"
	}

	return fmt.Sprintf("SELECT %s FROM `%s`.`%s` WHERE `%s` >= ? ORDER BY `%s` ASC LIMIT %d, %d", fields, ds.scheme, ds.table, ds.pKeyName, ds.pKeyName, 0, size)
}

type DbTarget struct {
	maxConcurrent int
	size          int
	scheme        string
	table         string

	isDone bool
	wg     sync.WaitGroup
	db     simple.Driver
	state  *DbTargetState

	target func()
	down   func()
}

func (dt *DbTarget) Start(config *DbTargetConfig) {
	dt.doStart(config)
}

func (dt *DbTarget) Done() {
	if dt.isDone {
		return
	}

	dt.isDone = true

	dt.down()
}

func (dt *DbTarget) State() {
	fmt.Println(fmt.Sprintf("Max Concurrent: %d \nSize: %d\nAmount: %d", dt.state.maxConcurrent, dt.state.size, dt.state.itemAmount))
}

func (dt *DbTarget) Close() {
	dt.wg.Wait()

	dt.db.Close()
}

func (dt *DbTarget) doStart(config *DbTargetConfig) {
	dt.maxConcurrent = config.MaxConcurrent
	dt.size = config.Size
	dt.scheme = config.Scheme
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

	for i := 0; i < dt.maxConcurrent; i++ {
		dt.wg.Add(1)

		go func() {
			dt.target()

			dt.wg.Done()
		}()
	}
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
	if len(d.fields) <= 0 {
		panic("fields is error")
	}

	var startTime time.Time
	var endTime time.Time

	for {
		if items, ok := <-d.itemsChan; ok {
			if len(items) <= 0 {
				continue
			}

			if len(items) > d.size {
				d.Send(items[:d.size])
				d.Send(items[d.size:])

				continue
			}

			startTime = time.Now()
			res, err := d.db.BulkInsertFromSliceSlice(d.table, d.fields, items)
			if err != nil {
				panic(err)
			}

			endTime = time.Now()

			num, err := res.RowsAffected()
			if err != nil {
				panic(err)
			}

			if num != int64(len(items)) {
				panic(fmt.Sprintf("target affected is error [%d != %d]", num, len(items)))
			}

			stateDuration := uintptr(unsafe.Pointer(&d.state.duration))

			atomic.AddUintptr(&stateDuration, uintptr(endTime.Sub(startTime)))
			atomic.AddInt64(&d.state.itemAmount, num)
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
	var startTime time.Time
	var endTime time.Time

	for {
		if items, ok := <-d.itemsChan; ok {
			if len(items) <= 0 {
				continue
			}

			if len(items) > d.size {
				d.Send(items[:d.size])
				d.Send(items[d.size:])

				continue
			}

			startTime = time.Now()
			res, err := d.db.BulkInsertFromSliceMap(d.table, items)
			if err != nil {
				panic(err)
			}

			endTime = time.Now()

			num, err := res.RowsAffected()
			if err != nil {
				panic(err)
			}

			if num != int64(len(items)) {
				panic(fmt.Sprintf("target affected is error [%d != %d]", num, len(items)))
			}

			stateDuration := uintptr(unsafe.Pointer(&d.state.duration))

			atomic.AddUintptr(&stateDuration, uintptr(endTime.Sub(startTime)))
			atomic.AddInt64(&d.state.itemAmount, num)

		} else {
			break
		}
	}
}
