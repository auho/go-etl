package task

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	simpledb "github.com/auho/go-simple-db/v3"
	"gorm.io/gorm"

	"github.com/auho/go-etl/v3/internal/testutil"
	"github.com/auho/go-etl/v3/internal/testutil/mysql"
)

var _ruleName = "a"
var _ruleTable = "rule_" + _ruleName
var _dataTable = "data"                              // master data source (read-only, copied per test)
var _updateAndTransferTable = "data_update_transfer" // for update and transfer
var _transferTable = "data_transfer"                 // for transfer
var _cleanDataTable = "clean_data"                   // for clean data
var _deletedDataTable = "deleted_data"               // for clean deleted
var _tagATable = "tag_data_a"
var _pkName = "did"
var _keyName = "name"
var _simpleDB *simpledb.SimpleDB
var _gormDB *gorm.DB
var _rule = &ruleTest{}
var _targetTagA = &targetTagATest{}
var _targetTransfer = &targetTransferTest{}
var _targetUpdateTransfer = &targetUpdateTransferTest{}

// test data dimensions (set in setUp, used in test assertions)
var _maxA int
var _maxB int
var _pageSize int64

// per-keyword expected counts derived from _maxA/_maxB
// base rows (10 items, cycled by i%10):
//   #0,#1: a一a一b一ab一123一中文      -> MostKey=a(amt=2),  ScanKey=5 rows (a,b,ab,123,中文)
//   #2,#3: 中文一中文一中文             -> MostKey=中文(amt=3), ScanKey=1 row (中文)
//   #4,#5,#6: xyz_no_match             -> no match
//   #7:     b一b一b一a                 -> MostKey=b(amt=3),  ScanKey=2 rows (b,a)
//   #8:     123一123一ab               -> MostKey=123(amt=2), ScanKey=2 rows (123,ab)
//   #9:     ab一ab一ab                 -> MostKey=ab(amt=3),  ScanKey=1 row (ab)

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setUp() {
	testutil.LoadEnv()
	_simpleDB, _gormDB = mysql.NewDB()

	query := ""

	// transfer table
	err := _simpleDB.Drop(_transferTable)
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
		"`a_keyword_amount` int(11) NOT NULL DEFAULT '0'," +
		"`xyz` varchar(30) NOT NULL DEFAULT ''," +
		"PRIMARY KEY (`did`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4"
	err = _gormDB.Exec(query).Error
	if err != nil {
		panic(err)
	}

	// data table (master)
	err = _simpleDB.Drop(_dataTable)
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
		"`a_keyword_amount` int(11) NOT NULL DEFAULT '0'," +
		"`xyz` varchar(30) NOT NULL DEFAULT ''," +
		"PRIMARY KEY (`did`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4"
	err = _gormDB.Exec(query).Error
	if err != nil {
		panic(err)
	}

	// 10 base rows, cycled by i%10:
	// #0,#1: a一a一b一ab一123一中文      -> MostKey=a(amt=2),  ScanKey=5 rows
	// #2,#3: 中文一中文一中文             -> MostKey=中文(amt=3), ScanKey=1 row
	// #4,#5,#6: xyz_no_match             -> no match
	// #7:     b一b一b一a                 -> MostKey=b(amt=3),  ScanKey=2 rows
	// #8:     123一123一ab               -> MostKey=123(amt=2), ScanKey=2 rows
	// #9:     ab一ab一ab                 -> MostKey=ab(amt=3),  ScanKey=1 row
	items := []any{
		"a一a一b一ab一123一中文",
		"a一a一b一ab一123一中文",
		"中文一中文一中文",
		"中文一中文一中文",
		"xyz_no_match",
		"xyz_no_match",
		"xyz_no_match",
		"b一b一b一a",
		"123一123一ab",
		"ab一ab一ab",
	}

	// maxA: 20~100, multiple of 10; total rows = maxA*maxB (200~1000)
	_maxA = (rand.Intn(9) + 2) * 10
	minB := (200 + _maxA - 1) / _maxA
	maxB := 1000 / _maxA
	_maxB = minB + rand.Intn(maxB-minB+1)
	_pageSize = int64(rand.Intn(31) + 20) // 20~50

	rows := make([][]any, 0)
	for i := 0; i < _maxA; i++ {
		rows = append(rows, []any{items[i%10]})
	}

	for i := 0; i < _maxB; i++ {
		err = _simpleDB.BulkInsertFromSliceSlice(_dataTable, []string{"name"}, rows, 2000)
		if err != nil {
			panic(err)
		}
	}

	var count int64
	err = _gormDB.Table(_dataTable).Count(&count).Error
	if err != nil {
		panic(err)
	}

	if count != int64(_maxA*_maxB) {
		panic(fmt.Sprintf("got: %d, want: %d", count, _maxA*_maxB))
	}

	err = _simpleDB.Drop(_updateAndTransferTable)
	if err != nil {
		panic(err)
	}

	err = _simpleDB.CopyStructure(_dataTable, _updateAndTransferTable)
	if err != nil {
		panic(err)
	}

	err = _simpleDB.Drop(_cleanDataTable)
	if err != nil {
		panic(err)
	}

	err = _simpleDB.CopyStructure(_dataTable, _cleanDataTable)
	if err != nil {
		panic(err)
	}

	err = _simpleDB.Drop(_deletedDataTable)
	if err != nil {
		panic(err)
	}

	err = _simpleDB.CopyStructure(_dataTable, _deletedDataTable)
	if err != nil {
		panic(err)
	}

	err = _simpleDB.Drop(_ruleTable)
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

	err = _gormDB.Exec(query).Error
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
	err = _gormDB.Exec(query).Error
	if err != nil {
		panic(err)
	}

	err = _simpleDB.Drop(_tagATable)
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
		"`a_keyword_amount` int(11) NOT NULL DEFAULT '0'," +
		"PRIMARY KEY (`id`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4"
	err = _gormDB.Exec(query).Error
	if err != nil {
		panic(err)
	}

}

func tearDown() {
	_ = _simpleDB.Drop(_ruleTable)
	_ = _simpleDB.Drop(_dataTable)
	_ = _simpleDB.Drop(_updateAndTransferTable)
	_ = _simpleDB.Drop(_transferTable)
	_ = _simpleDB.Drop(_cleanDataTable)
	_ = _simpleDB.Drop(_deletedDataTable)
	_ = _simpleDB.Drop(_tagATable)
}

// newTestSource creates a fresh copy of the master data table for a specific test,
// avoiding data competition between tests.
func newTestSource(t *testing.T, suffix string) *sourceTest {
	t.Helper()
	tableName := _dataTable + "_" + suffix
	_ = _simpleDB.Drop(tableName)
	err := _gormDB.Exec(fmt.Sprintf("CREATE TABLE `%s` LIKE `%s`", tableName, _dataTable)).Error
	if err != nil {
		t.Fatal(err)
	}
	err = _gormDB.Exec(fmt.Sprintf("INSERT INTO `%s` SELECT * FROM `%s`", tableName, _dataTable)).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = _simpleDB.Drop(tableName) })
	return &sourceTest{tableName: tableName}
}
