package task

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/auho/go-etl/v2/insight/app/conf"
	simpleDb "github.com/auho/go-simple-db/v2"
)

var _dsn = "test:Test123$@tcp(127.0.0.1:3306)/test"
var _ruleName = "a"
var _ruleTable = "rule_" + _ruleName
var _dataTable = "data"                              // data source
var _updateAndTransferTable = "data_update_transfer" // for update and transfer
var _transferTable = "data_transfer"                 // for transfer
var _cleanTable = "data_clean"                       // for clean
var _tagATable = "tag_data_a"
var _pkName = "did"
var _keyName = "name"
var _db *simpleDb.SimpleDB
var _rule = &ruleTest{}
var _source = &sourceTest{}
var _targetTagA = &targetTagATest{}
var _targetTagA1 = &targetTagA1Test{}
var _targetTagA2 = &targetTagA2Test{}
var _targetTransfer = &targetTransferTest{}
var _targetUpdateTransfer = &targetUpdateTransferTest{}
var _targetClean = &targetCleanTest{}

var dbConfig conf.DbConfig

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setUp() {
	var err error
	query := ""
	dbConfig.Driver = "mysql"
	dbConfig.Dsn = _dsn

	_db, err = dbConfig.BuildDB()
	if err != nil {
		panic(err)
	}

	err = _db.Drop(_transferTable)
	if err != nil {
		panic(err)
	}

	query = "CREATE TABLE `" + _transferTable + "` (" +
		"`did` int(11) unsigned NOT NULL AUTO_INCREMENT," +
		"`name` text," +
		"`a1` varchar(30) NOT NULL DEFAULT ''," +
		"`ab1` varchar(30) NOT NULL DEFAULT ''," +
		"`a_keyword` varchar(30) NOT NULL DEFAULT ''," +
		"`a_keyword_num` int(11) NOT NULL DEFAULT '0'," +
		"`xyz` varchar(30) NOT NULL DEFAULT ''," +
		"PRIMARY KEY (`did`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4"
	err = _db.Exec(query).Error
	if err != nil {
		panic(err)
	}

	err = _db.Drop(_dataTable)
	if err != nil {
		panic(err)
	}

	query = "CREATE TABLE `" + _dataTable + "` (" +
		"`did` int(11) unsigned NOT NULL AUTO_INCREMENT," +
		"`name` text," +
		"`a` varchar(30) NOT NULL DEFAULT ''," +
		"`ab` varchar(30) NOT NULL DEFAULT ''," +
		"`a_keyword` varchar(30) NOT NULL DEFAULT ''," +
		"`a_keyword_num` int(11) NOT NULL DEFAULT '0'," +
		"`xyz` varchar(30) NOT NULL DEFAULT ''," +
		"PRIMARY KEY (`did`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4"
	err = _db.Exec(query).Error
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
		err = _db.BulkInsertFromSliceSlice(_dataTable, []string{"name"}, rows, 2000)
		if err != nil {
			panic(err)
		}
	}

	var count int64
	err = _db.Table(_dataTable).Count(&count).Error
	if err != nil {
		panic(err)
	}

	if count != int64(maxA*maxB) {
		panic(fmt.Sprintf("%d != %d", _db.RowsAffected, maxA))
	}

	err = _db.Drop(_updateAndTransferTable)
	if err != nil {
		panic(err)
	}

	err = _db.Copy(_dataTable, _updateAndTransferTable)
	if err != nil {
		panic(err)
	}

	err = _db.Drop(_cleanTable)
	if err != nil {
		panic(err)
	}

	err = _db.Copy(_dataTable, _cleanTable)
	if err != nil {
		panic(err)
	}

	err = _db.Drop(_ruleTable)
	if err != nil {
		panic(err)
	}

	query = "CREATE TABLE `" + _ruleTable + "` (" +
		"`id` int(11) unsigned NOT NULL AUTO_INCREMENT," +
		"`a` varchar(30) NOT NULL DEFAULT ''," +
		"`ab` varchar(30) NOT NULL DEFAULT ''," +
		"`a_keyword` varchar(30) NOT NULL DEFAULT ''," +
		"`a_keyword_len` int(11) NOT NULL DEFAULT '0'," +
		"PRIMARY KEY (`id`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4"

	err = _db.Exec(query).Error
	if err != nil {
		panic(err)
	}

	query = "INSERT INTO `" + _ruleTable + "` (`a`, `ab`, `a_keyword`, `a_keyword_len`)" +
		"VALUES" +
		"('a','a1','a',1)," +
		"('a','a1','b',1)," +
		"('ab','ab1','ab',1)," +
		"('123','123','123',3)," +
		"('中文','中文1','中文',2)"
	err = _db.Exec(query).Error
	if err != nil {
		panic(err)
	}

	err = _db.Drop(_tagATable)
	if err != nil {
		panic(err)
	}

	query = "CREATE TABLE `" + _tagATable + "` (" +
		"`id` int(11) unsigned NOT NULL AUTO_INCREMENT," +
		"`did` int(11) NOT NULL DEFAULT '0'," +
		"`a` varchar(30) NOT NULL DEFAULT ''," +
		"`ab` varchar(30) NOT NULL DEFAULT ''," +
		"`a_keyword` varchar(30) NOT NULL DEFAULT ''," +
		"`a_keyword_num` int(11) NOT NULL DEFAULT '0'," +
		"PRIMARY KEY (`id`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4"
	err = _db.Exec(query).Error
	if err != nil {
		panic(err)
	}

}

func tearDown() {
	_ = _db.Drop(_ruleTable)
	_ = _db.Drop(_dataTable)
	_ = _db.Drop(_updateAndTransferTable)
	_ = _db.Drop(_transferTable)
	_ = _db.Drop(_cleanTable)
	_ = _db.Drop(_tagATable)
}
