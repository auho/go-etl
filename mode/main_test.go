package mode

import (
	"math/rand"
	"os"
	"testing"
	"time"

	goSimpleDb "github.com/auho/go-simple-db/v2"
)

var dsn = "test:Test123$@tcp(127.0.0.1:3306)/test"
var ruleName = "a"
var ruleTableName = "rule_" + ruleName
var keyName = "name"
var db *goSimpleDb.SimpleDB
var content = "b一ab一bc一abc一ab一123b一b123一中文一123一中文一一0123一1234一01234-a-ab-123-中文一b一中文一a"
var item = map[string]any{keyName: content}

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setUp() {
	rand.Seed(time.Now().UnixNano())

	var err error
	db, err = goSimpleDb.NewMysql(dsn)
	if err != nil {
		panic(err)
	}

	query := ""
	err = db.Drop(ruleTableName)
	if err != nil {
		panic(err)
	}

	query = "CREATE TABLE `" + ruleTableName + "` (" +
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

	query = "INSERT INTO `" + ruleTableName + "` (`a`, `ab`, `a_keyword`, `keyword_len`)" +
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
}

func tearDown() {
	_ = db.Drop(ruleTableName)
}
