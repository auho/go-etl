package storage

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/auho/go-simple-db/mysql"
	_ "github.com/go-sql-driver/mysql"
)

var dsn = "test:test@tcp(127.0.0.1:3306)/test"
var scheme = "test"
var tableName = "test_mysql"
var db *mysql.Mysql

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setup() {
	db = mysql.NewMysql(dsn)
	err := db.Connection()
	if err != nil {
		log.Fatalln(err)
	}

	err = db.Drop(tableName)
	if err != nil {
		log.Fatalln(err)
	}

	query := "CREATE TABLE IF NOT EXISTS `" + tableName + "` (" +
		"	`id` int(11) unsigned NOT NULL AUTO_INCREMENT," +
		"	`name` varchar(32) NOT NULL DEFAULT ''," +
		"	`value` int(11) NOT NULL DEFAULT '0'," +
		"	PRIMARY KEY (`id`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;"

	_, err = db.Exec(query)
	if err != nil {
		log.Fatalln(err)
	}
}

func tearDown() {
	err := db.Truncate(tableName)
	if err != nil {
		log.Fatalln(err)
	}

	err = db.Drop(tableName)
	if err != nil {
		log.Fatalln(err)
	}

	db.Close()
}

func TestNewDbTargetInsertInterface(t *testing.T) {
	config := &DbTargetConfig{}
	config.MaxConcurrent = 4
	config.Size = 100
	config.Driver = "mysql"
	config.Dsn = dsn
	config.Scheme = scheme
	config.Table = tableName

	d := NewDbTargetInsertInterface()
	d.Start(config)

	d.SetFields([]string{"name", "value"})

	rand.Seed(time.Now().UnixNano())
	maxI := rand.Intn(100)
	maxJ := rand.Intn(100)

	for i := 0; i < maxI; i++ {
		items := make([][]interface{}, 0, 100)
		for j := 0; j < maxJ; j++ {
			items = append(items, []interface{}{"name", time.Now().Unix()})
		}

		d.Send(items)
	}

	d.Done()
	d.Close()

	query := fmt.Sprintf("SELECT COUNT(*) AS _count FROM `%s`", tableName)
	res, err := db.QueryFieldInterface("_count", query)
	if err != nil {
		t.Error(err)
	}

	count, err := strconv.Atoi(string(res.([]uint8)))
	if err != nil {
		t.Error(err)
	}
	if maxI*maxJ != count {
		t.Error(fmt.Sprintf("count is error. %d != %d", maxI*maxJ, res))
	}
}
