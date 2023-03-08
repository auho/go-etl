package flow

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/auho/go-etl/v2/tool/conf"
	goSimpleDb "github.com/auho/go-simple-db/v2"
)

var dsn = "test:Test123$@tcp(127.0.0.1:3306)/test"
var ruleName = "a"
var ruleTable = "rule_" + ruleName
var dataTable = "data"                              // data source
var updateAndTransferTable = "data_update_transfer" // for update and transfer
var transferTable = "data_transfer"                 // for transfer
var cleanTable = "data_clean"                       // for clean
var tagATable = "tag_data_a"
var pkName = "did"
var keyName = "name"
var db *goSimpleDb.SimpleDB

var dbConfig conf.DbConfig

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setUp() {
	rand.Seed(time.Now().UnixNano())

	var err error
	query := ""
	dbConfig.Driver = "mysql"
	dbConfig.Dsn = dsn

	db, err = dbConfig.BuildDB()
	if err != nil {
		panic(err)
	}

	err = db.Drop(transferTable)
	if err != nil {
		panic(err)
	}

	query = "CREATE TABLE `" + transferTable + "` (" +
		"`did` int(11) unsigned NOT NULL AUTO_INCREMENT," +
		"`name` text," +
		"`a1` varchar(30) NOT NULL DEFAULT ''," +
		"`ab1` varchar(30) NOT NULL DEFAULT ''," +
		"`a_keyword` varchar(30) NOT NULL DEFAULT ''," +
		"`a_keyword_num` int(11) NOT NULL DEFAULT '0'," +
		"`xyz` varchar(30) NOT NULL DEFAULT ''," +
		"PRIMARY KEY (`did`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4"
	err = db.Exec(query).Error
	if err != nil {
		panic(err)
	}

	err = db.Drop(dataTable)
	if err != nil {
		panic(err)
	}

	query = "CREATE TABLE `" + dataTable + "` (" +
		"`did` int(11) unsigned NOT NULL AUTO_INCREMENT," +
		"`name` text," +
		"`a` varchar(30) NOT NULL DEFAULT ''," +
		"`ab` varchar(30) NOT NULL DEFAULT ''," +
		"`a_keyword` varchar(30) NOT NULL DEFAULT ''," +
		"`a_keyword_num` int(11) NOT NULL DEFAULT '0'," +
		"`xyz` varchar(30) NOT NULL DEFAULT ''," +
		"PRIMARY KEY (`did`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4"
	err = db.Exec(query).Error
	if err != nil {
		panic(err)
	}

	items := []any{
		"b一ab一bc一abc一123b一b123一123一0123一1234一01234",
		`中文一b中文123一123中文b一中bb文一中123文一中00文一中aa文一中00文一中aa文一中中文文一中二二文一
123一一`,
		`文中ba321b--#$%^&*()_`,
	}

	maxA := (rand.Intn(100) + 10) * 3
	maxB := rand.Intn(100) + 10
	rows := make([][]any, 0)
	for i := 0; i < maxA; i++ {
		rows = append(rows, []any{items[i%3]})
	}

	for i := 0; i < maxB; i++ {
		err = db.BulkInsertFromSliceSlice(dataTable, []string{"name"}, rows, 2000)
		if err != nil {
			panic(err)
		}
	}

	var count int64
	err = db.Table(dataTable).Count(&count).Error
	if err != nil {
		panic(err)
	}

	if count != int64(maxA*maxB) {
		panic(fmt.Sprintf("%d != %d", db.RowsAffected, maxA))
	}

	err = db.Drop(updateAndTransferTable)
	if err != nil {
		panic(err)
	}

	err = db.Copy(dataTable, updateAndTransferTable)
	if err != nil {
		panic(err)
	}

	err = db.Drop(cleanTable)
	if err != nil {
		panic(err)
	}

	err = db.Copy(dataTable, cleanTable)
	if err != nil {
		panic(err)
	}

	err = db.Drop(ruleTable)
	if err != nil {
		panic(err)
	}

	query = "CREATE TABLE `" + ruleTable + "` (" +
		"`id` int(11) unsigned NOT NULL AUTO_INCREMENT," +
		"`a` varchar(30) NOT NULL DEFAULT ''," +
		"`ab` varchar(30) NOT NULL DEFAULT ''," +
		"`a_keyword` varchar(30) NOT NULL DEFAULT ''," +
		"`keyword_len` int(11) NOT NULL DEFAULT '0'," +
		"PRIMARY KEY (`id`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4"

	err = db.Exec(query).Error
	if err != nil {
		panic(err)
	}

	query = "INSERT INTO `" + ruleTable + "` (`a`, `ab`, `a_keyword`, `keyword_len`)" +
		"VALUES" +
		"('a','a1','a',1)," +
		"('a','a1','b',1)," +
		"('ab','ab1','ab',1)," +
		"('123','123','123',3)," +
		"('中文','中文1','中文',2)"
	err = db.Exec(query).Error
	if err != nil {
		panic(err)
	}

	err = db.Drop(tagATable)
	if err != nil {
		panic(err)
	}

	query = "CREATE TABLE `" + tagATable + "` (" +
		"`id` int(11) unsigned NOT NULL AUTO_INCREMENT," +
		"`did` int(11) NOT NULL DEFAULT '0'," +
		"`a` varchar(30) NOT NULL DEFAULT ''," +
		"`ab` varchar(30) NOT NULL DEFAULT ''," +
		"`a_keyword` varchar(30) NOT NULL DEFAULT ''," +
		"`a_keyword_num` int(11) NOT NULL DEFAULT '0'," +
		"PRIMARY KEY (`id`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4"
	err = db.Exec(query).Error
	if err != nil {
		panic(err)
	}

}

func tearDown() {
	_ = db.Drop(ruleTable)
	_ = db.Drop(dataTable)
	_ = db.Drop(updateAndTransferTable)
	_ = db.Drop(transferTable)
	_ = db.Drop(cleanTable)
	_ = db.Drop(tagATable)
}
