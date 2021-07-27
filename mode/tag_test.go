package mode

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/auho/go-etl/means/tager"
	"github.com/auho/go-simple-db/mysql"
)

var dsn = "test:test@tcp(127.0.0.1:3306)/test"
var ruleName = "a"
var ruleTableName = "rule_" + ruleName
var keyName = "name"
var db *mysql.Mysql

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setUp() {
	rand.Seed(time.Now().UnixNano())
	db = mysql.NewMysql(dsn)
	err := db.Connection()
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
	_, err = db.Exec(query)
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
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func tearDown() {
	_ = db.Drop(ruleTableName)
}

func Test_Tag(t *testing.T) {
	item := make(map[string]interface{})
	item[keyName] = "b一ab一bc一abc一ab一123b一b123一中文一123一中文一一0123一1234一01234-a-ab-123-中文一b一中文一a"

	ttm := tager.NewTagKeyMeans(ruleName, db)
	ti1 := NewTagInsert([]string{keyName}, ttm)
	results := ti1.Do(item)
	if len(results) <= 0 {
		t.Error("error")
	}
	fmt.Println(results)

	tmtm := tager.NewTagMostTextMeans(ruleName, db)
	ti2 := NewTagInsert([]string{keyName}, tmtm)
	results = ti2.Do(item)
	if len(results) <= 0 {
		t.Error("error")
	}
	fmt.Println(results)

	tmkm := tager.NewTagMostKeyMeans(ruleName, db)
	ti3 := NewTagInsert([]string{keyName}, tmkm)
	results = ti3.Do(item)
	if len(results) <= 0 {
		t.Error("error")
	}
	fmt.Println(results)
}
