package storage

import (
	"fmt"
	"math"
	"strings"
	"sync"

	"etl/lib/conf"

	"github.com/auho/go-simple-db/simple"
)

type Source interface {
}

type Target interface {
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
	sw           sync.WaitGroup
	pageChan     chan interface{}
	itemsChan    chan []map[string]interface{}

	db simple.Driver
}

func NewDbSource(dsConfig conf.DbSourceConfig) *DbSource {
	ds := &DbSource{}
	ds.maxConcurrent = dsConfig.MaxConcurrent
	ds.size = dsConfig.Size
	ds.page = dsConfig.Page
	ds.scheme = dsConfig.Scheme
	ds.pKeyName = dsConfig.PKeyName
	ds.fields = dsConfig.Fields

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
		ds.sw.Add(1)

		go ds.rows(rowsQuery)
	}

	minPk := ds.getMinPk()
	ds.pageChan <- minPk

	prePk := minPk
	for page := 1; page < maxPage; page++ {
		nextPk := ds.getNextPk(prePk, nextPkQuery)
		ds.pageChan <- nextPk
		prePk = nextPk
	}

	close(ds.pageChan)

	ds.sw.Wait()
}

func (ds *DbSource) Receive() []map[string]interface{} {
	return <-ds.itemsChan
}

func (ds *DbSource) rows(query string) {
	for pk := range ds.pageChan {
		rows, err := ds.db.QueryInterface(query, pk)
		if err != nil {
			ds.retryPage(pk)
		}

		ds.retryPageMap.Delete(pk)

		if rows == nil {
			continue
		}
	}

	ds.sw.Done()
}

func (ds *DbSource) close() {

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
	maxConcurrent  int
	size           int
	driver         string
	dsn            string
	scheme         string
	table          string
	sqlTemplate    string
	maxSqlTemplate string
}
