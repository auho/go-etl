package database

import (
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/auho/go-simple-db/mysql"
)

var dsn = "test:Test123$@tcp(127.0.0.1:3306)/test"
var tableName = "test_mysql"
var pkName = "id"
var db *mysql.Mysql
var sourceConfig = &DbSourceConfig{}
var targetConfig = &DbTargetConfig{}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setup() {
	rand.Seed(time.Now().UnixNano())
	sourceConfig.MaxConcurrent = rand.Intn(4) + 1
	sourceConfig.Size = rand.Intn(500) + 1
	sourceConfig.Page = rand.Intn(1000) + 1
	sourceConfig.Driver = "mysql"
	sourceConfig.Dsn = dsn
	sourceConfig.Table = tableName
	sourceConfig.PKeyName = pkName

	targetConfig.MaxConcurrent = rand.Intn(4) + 1
	targetConfig.Size = rand.Intn(500) + 1
	targetConfig.Driver = "mysql"
	targetConfig.Dsn = dsn
	targetConfig.Table = tableName

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
