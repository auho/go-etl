package runner

import (
	"fmt"
	"math/rand"
	"testing"
)

// baseNames are the 10 source rows cycled by i%10.
//
// Per-pattern mapping (drives expected counts in run_test.go):
//
//	#0,#1:    a一a一b一ab一123一中文   -> MostKey=a(amt=2),  ScanKey=5 rows (a,b,ab,123,中文)
//	#2,#3:    中文一中文一中文          -> MostKey=中文(amt=3), ScanKey=1 row  (中文)
//	#4,#5,#6: xyz_no_match            -> no match
//	#7:       b一b一b一a              -> MostKey=b(amt=3),  ScanKey=2 rows (b,a)
//	#8:       123一123一ab            -> MostKey=123(amt=2), ScanKey=2 rows (123,ab)
//	#9:       ab一ab一ab              -> MostKey=ab(amt=3),  ScanKey=1 row  (ab)
//
// Base rows per pattern: a=2, 中文=2, noMatch=3, b=1, 123=1, ab=1 (out of 10).
// _maxA is a multiple of 10, so each pattern appears (_maxA/10) times per batch,
// repeated _maxB batches => total = patternCount * (_maxA/10) * _maxB.
var baseNames = []any{
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

// --- DDL templates (table name injected via fmt.Sprintf) ---

const dataTableDDL = "CREATE TABLE `%s` (" +
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

// transferTableDDL mirrors dataTableDDL but renames a->a1, ab->ab1 (alias target).
const transferTableDDL = "CREATE TABLE `%s` (" +
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

const tagATableDDL = "CREATE TABLE `%s` (" +
	"`id` int(11) unsigned NOT NULL AUTO_INCREMENT," +
	"`did` int(11) NOT NULL DEFAULT '0'," +
	"`a` varchar(30) NOT NULL DEFAULT ''," +
	"`ab` varchar(30) NOT NULL DEFAULT ''," +
	"`a_keyword` varchar(30) NOT NULL DEFAULT ''," +
	"`a_keyword_num` int(11) NOT NULL DEFAULT '0'," +
	"`a_keyword_amount` int(11) NOT NULL DEFAULT '0'," +
	"PRIMARY KEY (`id`)" +
	") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4"

const ruleTableDDL = "CREATE TABLE `%s` (" +
	"`id` int(11) unsigned NOT NULL AUTO_INCREMENT," +
	"`a` varchar(30) NOT NULL DEFAULT ''," +
	"`ab` varchar(30) NOT NULL DEFAULT ''," +
	"`a_keyword` varchar(30) NOT NULL DEFAULT ''," +
	"`a_keyword_len` int(11) NOT NULL DEFAULT '0'," +
	"PRIMARY KEY (`id`)" +
	") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4"

const ruleDataInsert = "INSERT INTO `%s` (`a`, `ab`, `a_keyword`, `a_keyword_len`) VALUES" +
	"('a','a1','a',1)," +
	"('a','a1','b',1)," +
	"('ab','ab1','ab',1)," +
	"('123','123','123',3)," +
	"('中文','中文1','中文',2)"

// --- error helpers (setup-only; panicking in setUp is acceptable) ---

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func mustExec(query string) {
	must(_gormDB.Exec(query).Error)
}

func mustDrop(table string) {
	must(_simpleDB.Drop(table))
}

// --- setup steps, called from setUp() in main_test.go ---

// setupDataDimensions picks random _maxA/_maxB/_pageSize so that
// total rows = maxA*maxB stays within [200, 1000] and each page is 20~50.
func setupDataDimensions() {
	_maxA = (rand.Intn(9) + 2) * 10 // 20~100, multiple of 10
	minB := (200 + _maxA - 1) / _maxA
	maxB := 1000 / _maxA
	_maxB = minB + rand.Intn(maxB-minB+1)
	_pageSize = int64(rand.Intn(31) + 20) // 20~50
}

// createMasterDataTable builds the read-only master data table.
func createMasterDataTable() {
	mustDrop(_dataTable)
	mustExec(fmt.Sprintf(dataTableDDL, _dataTable))
}

// generateMasterData fills the master table with _maxA*_maxB rows.
func generateMasterData() {
	rows := make([][]any, 0, _maxA)
	for i := 0; i < _maxA; i++ {
		rows = append(rows, []any{baseNames[i%10]})
	}

	for i := 0; i < _maxB; i++ {
		must(_simpleDB.BulkInsertFromSliceSlice(_dataTable, []string{"name"}, rows, 2000))
	}

	var count int64
	must(_gormDB.Table(_dataTable).Count(&count).Error)
	if count != int64(_maxA*_maxB) {
		panic(fmt.Sprintf("got: %d, want: %d", count, _maxA*_maxB))
	}
}

// createTransferTable builds the transfer destination table.
func createTransferTable() {
	mustDrop(_transferTable)
	mustExec(fmt.Sprintf(transferTableDDL, _transferTable))
}

// createUpdateAndTransferTable clones the master structure for update+transfer.
func createUpdateAndTransferTable() {
	mustDrop(_updateAndTransferTable)
	must(_simpleDB.CopyStructure(_dataTable, _updateAndTransferTable))
}

// createCleanTables clones the master structure for clean data + deleted data.
func createCleanTables() {
	mustDrop(_cleanDataTable)
	must(_simpleDB.CopyStructure(_dataTable, _cleanDataTable))
	mustDrop(_deletedDataTable)
	must(_simpleDB.CopyStructure(_dataTable, _deletedDataTable))
}

// createTagATable builds the tag destination table.
func createTagATable() {
	mustDrop(_tagATable)
	mustExec(fmt.Sprintf(tagATableDDL, _tagATable))
}

// createRuleTable builds and populates the rule lookup table.
func createRuleTable() {
	mustDrop(_ruleTable)
	mustExec(fmt.Sprintf(ruleTableDDL, _ruleTable))
	mustExec(fmt.Sprintf(ruleDataInsert, _ruleTable))
}

// newTestSource creates a fresh copy of the master data table for a specific
// test, avoiding data competition between tests. The copy is dropped via
// t.Cleanup when the test finishes.
func newTestSource(t *testing.T, suffix string) *sourceTest {
	t.Helper()
	tableName := _dataTable + "_" + suffix
	mustDrop(tableName)
	mustExec(fmt.Sprintf("CREATE TABLE `%s` LIKE `%s`", tableName, _dataTable))
	mustExec(fmt.Sprintf("INSERT INTO `%s` SELECT * FROM `%s`", tableName, _dataTable))
	t.Cleanup(func() { _ = _simpleDB.Drop(tableName) })
	return &sourceTest{tableName: tableName}
}
