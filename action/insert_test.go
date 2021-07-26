package action

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/auho/go-simple-db/mysql"
)

var dsn = "test:test@tcp(127.0.0.1:3306)/test"
var scheme = "test"
var ruleTableName = "rule_a"
var dataTableName = "data"
var tagTableName = "tag_data_a"
var pkName = "id"
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
	err = db.Drop(dataTableName)
	if err != nil {
		panic(err)
	}

	query = "CREATE TABLE `" + dataTableName + "` (" +
		"`did` int(11) unsigned NOT NULL AUTO_INCREMENT," +
		"`name` text," +
		"PRIMARY KEY (`did`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4"
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}

	err = db.Drop(ruleTableName)
	if err != nil {
		panic(err)
	}

	query = "CREATE TABLE `" + ruleTableName + "` (" +
		"`id` int(11) unsigned NOT NULL AUTO_INCREMENT," +
		"`a` varchar(30) NOT NULL DEFAULT ''," +
		"`ab` varchar(30) NOT NULL," +
		"`a_keyword` varchar(30) NOT NULL DEFAULT ''," +
		"`keyword_num` int(11) NOT NULL DEFAULT '0'," +
		"PRIMARY KEY (`id`)" +
		") ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4"

	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}

	err = db.Drop(tagTableName)
	if err != nil {
		panic(err)
	}

	query = "CREATE TABLE `" + tagTableName + "` (" +
		"`id` int(11) unsigned NOT NULL AUTO_INCREMENT," +
		"`did` int(11) NOT NULL DEFAULT '0'," +
		"`a` varchar(30) NOT NULL DEFAULT ''," +
		"`ab` varchar(30) NOT NULL DEFAULT ''," +
		"`a_keyword` varchar(30) NOT NULL DEFAULT ''," +
		"`a_keyword_num` int(11) NOT NULL DEFAULT '0'," +
		"PRIMARY KEY (`id`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4"
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}

	items := []interface{}{
		"b一ab一bc一abc一123b一b123一123一0123一1234一01234",
		`中文一b中文123一123中文b一中bb文一中123文一中00文一中aa文一中00文一中aa文一中中文文一中二二文一
123一一`,
		`文中ba321b--#$%^&*()_`,
	}

	rows := make([][]interface{}, 0)
	for i := 0; i < rand.Intn(100)+100; i++ {
		rows = append(rows, []interface{}{items[i%3]})
	}

	for i := 0; i < rand.Intn(100)+100; i++ {
		_, err = db.BulkInsertFromSliceSlice(dataTableName, []string{"name"}, rows)
		if err != nil {
			panic(err)
		}
	}
}

func tearDown() {
	_ = db.Truncate(ruleTableName)
	_ = db.Truncate(dataTableName)

	_ = db.Drop(ruleTableName)
	_ = db.Drop(dataTableName)
}

func Test_Tag(t *testing.T) {
}
