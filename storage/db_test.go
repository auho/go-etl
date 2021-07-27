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

func Test_DbTargetInsertInterface(t *testing.T) {
	d := NewDbTargetInsertInterface(targetConfig)
	d.Start()

	d.SetFields([]string{"name", "value"})

	maxI := rand.Intn(100) + 1
	maxJ := rand.Intn(100) + 1

	nameValue := fmt.Sprintf("name slice %d", time.Now().Unix()/(rand.Int63n(999970)+1))
	for i := 0; i < maxI; i++ {
		items := make([][]interface{}, 0, 100)
		for j := 0; j < maxJ; j++ {
			items = append(items, []interface{}{nameValue, time.Now().Unix()})
		}

		d.Send(items)
	}

	d.Done()
	d.Close()
	//d.State()

	query := fmt.Sprintf("SELECT COUNT(*) AS _count FROM `%s` WHERE `%s` = ?", tableName, "name")
	res, err := db.QueryFieldInterface("_count", query, nameValue)
	if err != nil {
		t.Error(err)
	}

	count := int(res.(int64))
	if maxI*maxJ != count {
		t.Error(fmt.Sprintf("count is error. %d != %d", maxI*maxJ, res))
	}
}
func Test_DbTargetInsertMap(t *testing.T) {
	d := NewDbTargetInsertMap(targetConfig)
	d.Start()

	maxI := rand.Intn(100) + 1
	maxJ := rand.Intn(100) + 1

	nameValue := fmt.Sprintf("name map %d", time.Now().Unix()/(rand.Int63n(999970)+1))
	for i := 0; i < maxI; i++ {
		items := make([]map[string]interface{}, 0, 100)
		for j := 0; j < maxJ; j++ {
			items = append(items, map[string]interface{}{"name": nameValue, "value": time.Now().Unix()})
		}

		d.Send(items)
	}

	d.Done()
	d.Close()
	//d.State()

	query := fmt.Sprintf("SELECT COUNT(*) AS _count FROM `%s` WHERE `%s` = ?", tableName, "name")
	res, err := db.QueryFieldInterface("_count", query, nameValue)
	if err != nil {
		t.Error(err)
	}

	count := int(res.(int64))
	if maxI*maxJ != count {
		t.Error(fmt.Sprintf("count is error. %d != %d", maxI*maxJ, count))
	}
}

func Test_DbSource(t *testing.T) {
	d := NewDbSource(sourceConfig)
	d.Start()

	amount := getAmount()

	count := 0
	page := 0
	for {
		items, ok := <-d.itemsChan
		if ok == false {
			break
		}

		if len(items) <= 0 {
			t.Error(fmt.Sprintf("items size is error %d != %d", len(items), sourceConfig.Size))
		}

		page += 1
		count += len(items)
	}

	if page > sourceConfig.Page {
		t.Error(fmt.Sprintf("page is error %d != %d", page, sourceConfig.Page))
	}

	if amount != count {
		t.Error(fmt.Sprintf("amount is error %d != %d", amount, count))
	}
}

func getAmount() int {
	query := fmt.Sprintf("SELECT COUNT(*) AS _count FROM `%s`", tableName)
	res, err := db.QueryFieldInterface("_count", query)
	if err != nil {
		return 0
	}

	n, err := strconv.Atoi(string(res.([]uint8)))
	if err != nil {
		return 0
	}

	return n
}
