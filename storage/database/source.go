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
	s := &DbSource{}
	s.initDb(config)

	return s
}

func (s *DbSource) Start() {
	s.State.status = "source start"

	s.pageChan = make(chan interface{}, s.maxConcurrent)
	s.itemsChan = make(chan []map[string]interface{}, s.maxConcurrent)

	rowsQuery := s.generateRowsQuery(s.size)
	for i := 0; i < s.maxConcurrent; i++ {
		s.wg.Add(1)

		go s.source(rowsQuery)
	}

	go s.sourcePage()

	go s.Close()
}

func (s *DbSource) Consume() ([]map[string]interface{}, bool) {
	items, ok := <-s.itemsChan

	return items, ok
}

func (s *DbSource) Close() {
	s.wg.Wait()

	close(s.itemsChan)
	s.db.Close()

	s.State.status = s.State.DoneStatus()
}

func (s *DbSource) initDb(config *DbSourceConfig) {
	s.maxConcurrent = config.MaxConcurrent
	s.size = config.Size
	s.page = config.Page
	s.table = config.Table
	s.pKeyName = config.PKeyName
	s.fields = config.Fields

	config.check()

	s.State = newDbSourceState()
	s.State.size = s.size
	s.State.maxConcurrent = s.maxConcurrent
	s.State.title = fmt.Sprintf("source[%s]", s.table)

	var err error
	s.db, err = simple.NewDriver(config.Driver, config.Dsn)
	if err != nil {
		panic(err)
	}

	err = s.db.Ping()
	if err != nil {
		panic(err)
	}
}

func (s *DbSource) sourcePage() {
	amount := s.getTableAmount()
	if amount <= 0 {
		panic("db source table amount is  error 0")
	}

	maxPage := int(math.Ceil(float64(amount) / float64(s.size)))
	if s.page > 0 && s.page <= maxPage {
		maxPage = s.page
	}

	if maxPage <= 0 {
		panic(fmt.Sprintf("max Page is error %d", maxPage))
	}

	s.State.page = maxPage

	nextPkQuery := s.generateNextPkQuery(s.size)
	go func() {
		minPk := s.getMinPk()
		s.pageChan <- minPk

		prePk := minPk
		for page := 1; page < maxPage; page++ {
			nextPk := s.getNextPk(prePk, nextPkQuery)
			if nextPk == nil {
				break
			}

			s.pageChan <- nextPk
			prePk = nextPk
			s.lastPagePk = nextPk
		}

		close(s.pageChan)
	}()
}

func (s *DbSource) source(query string) {
	for {
		pk, ok := <-s.pageChan
		if ok == false {
			break
		}

		rows := s.retryPage(query, pk)
		if rows == nil {
			continue
		}

		s.itemsChan <- rows

		atomic.AddInt64(&s.State.itemAmount, int64(len(rows)))

		s.State.status = fmt.Sprintf("source last page pk[%v]; item amount[%d]", s.lastPagePk, s.State.itemAmount)
	}

	s.wg.Done()
}

func (s *DbSource) retryPage(query string, pk interface{}) []map[string]interface{} {
	var rows []map[string]interface{}
	var err error

	for i := 0; i < 3; i++ {
		rows, err = s.db.QueryInterface(query, pk)
		if err != nil {
			continue
		} else {
			break
		}
	}

	if err != nil {
		panic(fmt.Sprintf("source pk[%v] error.  \n%v\n%s", pk, err, query))
	}

	return rows
}

func (s *DbSource) getMinPk() interface{} {
	query := fmt.Sprintf("SELECT `%s` FROM `%s` ORDER BY `%s` ASC LIMIT 0, ?", s.pKeyName, s.table, s.pKeyName)
	minPk, err := s.db.QueryFieldInterface(s.pKeyName, query, 1)
	if err != nil {
		panic(err)
	}

	return minPk
}

func (s *DbSource) getTableAmount() int64 {
	query := fmt.Sprintf("SELECT COUNT(*) AS 'amount' FROM `%s`", s.table)
	res, err := s.db.QueryFieldInterface("amount", query)
	if err != nil {
		panic(err)
	}

	amount, err := strconv.ParseInt(string(res.([]uint8)), 10, 64)
	if err != nil {
		panic(err)
	}

	return amount
}

func (s *DbSource) getNextPk(pk interface{}, query string) interface{} {
	nextPk, err := s.db.QueryFieldInterface(s.pKeyName, query, pk)
	if err != nil {
		panic(err)
	} else {
		return nextPk
	}
}

func (s *DbSource) generateNextPkQuery(size int) string {
	return fmt.Sprintf("SELECT `%s` FROM `%s` WHERE `%s` >= ? ORDER BY `%s` ASC LIMIT %d, %d", s.pKeyName, s.table, s.pKeyName, s.pKeyName, size, 1)
}

func (s *DbSource) generateRowsQuery(size int) string {
	fields := ""
	if len(s.fields) <= 0 {
		fields = "*"
	} else {
		fields = "`" + strings.Join(s.fields, "`, `") + "`"
	}

	return fmt.Sprintf("SELECT %s FROM `%s` WHERE `%s` >= ? ORDER BY `%s` ASC LIMIT %d, %d", fields, s.table, s.pKeyName, s.pKeyName, 0, size)
}
