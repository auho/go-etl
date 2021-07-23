package storage

import (
	"fmt"
	"log"
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

	sw       sync.WaitGroup
	pageChan chan int

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

	return ds
}

func (ds *DbSource) Start() {
	page := 0
	maxRows := ds.maxRows()
	if ds.page > 0 {
		page = ds.page
	} else {
		page = int(math.Ceil(float64(maxRows) / float64(ds.size)))
	}

	for i := 0; i < ds.maxConcurrent; i++ {
		ds.sw.Add(1)

		go ds.rows()
	}

	for i := 0; i < page; i++ {
		ds.pageChan <- i
	}

	close(ds.pageChan)

	ds.sw.Wait()
}

func (ds *DbSource) rows() {
	for i := range ds.pageChan {
		query := ds.generateRowsSql(ds.size)

		rows, err := ds.db.QueryInterface(query)
		if err != nil {
			ds.pageChan <- i
		}

		if rows == nil {
			continue
		}

	}

	ds.sw.Done()
}

func (ds *DbSource) maxRows() int64 {
	query := ds.generateMaxSql()
	value, err := ds.db.QueryFieldInterface(ds.pKeyName, query)
	if err != nil {
		log.Fatalln(err)
	}

	return value.(int64)
}

func (ds *DbSource) close() {

}

func (ds *DbSource) generateRowsSql(size int) string {
	fields := "`" + strings.Join(ds.fields, "`, `") + "`"
	return fmt.Sprintf("SELECT %s FROM `%s`.`%s` WHERE `%s` > ? ORDER BY %s ASC LIMIT %d, %d", fields, ds.scheme, ds.table, ds.pKeyName, ds.pKeyName, 0, size)
}

func (ds *DbSource) generateMaxSql() string {
	return fmt.Sprintf("SELECT `%s` FROM `%s`.`%s` ORDER BY %s DESC LIMIT %d", ds.pKeyName, ds.scheme, ds.table, ds.pKeyName, 1)
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
