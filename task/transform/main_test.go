package transform

import (
	"context"
	"os"
	"testing"

	"github.com/auho/go-etl/v3/internal/testutil"
	"github.com/auho/go-etl/v3/internal/testutil/mysql"
	simpledb "github.com/auho/go-simple-db/v3"
	"gorm.io/gorm"
)

var _ruleName = "transform"
var _ruleTableName = "rule_" + _ruleName
var _keyName = "name"
var _simpleDB *simpledb.SimpleDB
var _gormDB *gorm.DB
var _content = "b一ab一bc一abc一ab一123b一b123一中文一123一中文一一0123一1234一01234-a-ab-123-中文一b一中文一a"
var _item = map[string]any{_keyName: _content}
var _rule = &ruleTest{}

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setUp() {
	err := testutil.LoadEnv()
	if err != nil {
		panic(err)
	}

	_simpleDB, _gormDB, err = mysql.NewDB()
	if err != nil {
		panic(err)
	}

	query := ""
	err = _simpleDB.Drop(context.TODO(), _ruleTableName)
	if err != nil {
		panic(err)
	}

	query = "CREATE TABLE `" + _ruleTableName + "` (" +
		"`id` int(11) unsigned NOT NULL AUTO_INCREMENT," +
		"`a` varchar(30) NOT NULL DEFAULT ''," +
		"`ab` varchar(30) NOT NULL DEFAULT ''," +
		"`a_keyword` varchar(30) NOT NULL DEFAULT ''," +
		"`a_keyword_len` int(11) NOT NULL DEFAULT '0'," +
		"PRIMARY KEY (`id`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4"
	err = _gormDB.Exec(query).Error
	if err != nil {
		panic(err)
	}

	query = "INSERT INTO `" + _ruleTableName + "` (`a`, `ab`, `a_keyword`, `a_keyword_len`)" +
		"VALUES" +
		"('a','a1','a',1)," +
		"('a','a1','b',1)," +
		"('ab','ab1','ab',1)," +
		"('123','123','123',3)," +
		"('中文','中文1','中文',2)"
	err = _gormDB.Exec(query).Error
	if err != nil {
		panic(err)
	}
}

func tearDown() {
	_ = _simpleDB.Drop(context.TODO(), _ruleTableName)
}
