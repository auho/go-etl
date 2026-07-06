package runner

import (
	"fmt"
	"os"
	"testing"

	simpledb "github.com/auho/go-simple-db/v3"
	"gorm.io/gorm"

	"github.com/auho/go-etl/v3/internal/testutil"
	"github.com/auho/go-etl/v3/internal/testutil/mysql"
)

// --- table names ---
var _ruleName = "runner"
var _ruleTable = "rule_" + _ruleName
var _dataTable = "runner_data"                              // master data source (read-only, copied per test)
var _updateAndTransferTable = "runner_data_update_transfer" // for update and transfer
var _transferTable = "runner_data_transfer"                 // for transfer
var _cleanDataTable = "runner_clean_data"                   // for clean data
var _deletedDataTable = "runner_deleted_data"               // for clean deleted
var _tagATable = "runner_tag_data_a"

// --- field names ---
var _pkName = "did"
var _keyName = "name"

// --- shared db handles (set in setUp) ---
var _simpleDB *simpledb.SimpleDB
var _gormDB *gorm.DB

// --- shared fixtures ---
var _rule = &ruleTest{}
var _targetTagA = &targetTagATest{}
var _targetTransfer = &targetTransferTest{}
var _targetUpdateTransfer = &targetUpdateTransferTest{}

// --- test data dimensions (set in setUp, used in test assertions) ---
// See setup_test.go (baseNames) for the per-keyword expected-counts mapping.
var _maxA int
var _maxB int
var _pageSize int64

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

// setUp prepares the shared schema and master data. The detailed DDL and data
// generation live in setup_test.go; this function only orchestrates the order.
func setUp() {
	testutil.LoadEnv()
	if os.Getenv("MYSQL_DSN") == "" {
		fmt.Println("skip: MYSQL_DSN not set")
		os.Exit(0)
	}
	_simpleDB, _gormDB = mysql.NewDB()

	setupDataDimensions()
	createMasterDataTable()
	generateMasterData()
	createTransferTable()
	createUpdateAndTransferTable()
	createCleanTables()
	createTagATable()
	createRuleTable()
}

// tearDown drops every shared table. Per-test copies are dropped via t.Cleanup.
func tearDown() {
	for _, table := range []string{
		_ruleTable,
		_dataTable,
		_updateAndTransferTable,
		_transferTable,
		_cleanDataTable,
		_deletedDataTable,
		_tagATable,
	} {
		_ = _simpleDB.Drop(table)
	}
}
