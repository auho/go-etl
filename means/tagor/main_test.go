package tagor

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/auho/go-simple-db/mysql"
)

var dsn = "test:test@tcp(127.0.0.1:3306)/test"
var ruleName = "a"
var ruleTableName = "rule_" + ruleName
var dataRuleTableName = "rule_data_" + ruleName
var db *mysql.Mysql
var contents = []string{
	`b一ab一bc一abc一123b一b123一123一0123一1234一01234一`,
	`中文一b中文123一123中文b一中bb文一中123文一中00文一中aa文一中00文一中aa文一中中文文一中二二文一
123一一`,
}

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

	query = "CREATE TABLE IF NOT EXISTS`" + ruleTableName + "` (" +
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

	err = db.Copy(ruleTableName, dataRuleTableName)
	if err != nil {
		panic(err)
	}

	query = fmt.Sprintf("INSERT INTO `%s` SELECT * FROM `%s`", dataRuleTableName, ruleTableName)
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func tearDown() {
	_ = db.Drop(ruleTableName)
	_ = db.Drop(dataRuleTableName)
}
