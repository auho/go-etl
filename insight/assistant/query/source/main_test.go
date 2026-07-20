package source

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/auho/go-etl/v3/internal/testutil/mysql"
	simpledb "github.com/auho/go-simple-db/v3"
	testutil "github.com/auho/go-toolkit-testutil"
	"gorm.io/gorm"
)

var _simpleDB *simpledb.SimpleDB
var _gormDB *gorm.DB

const _testTable = "test_source"

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

// setup prepares the test table and data.
// When MySQL is unavailable _gormDB stays nil and the setup is skipped
// so that in-memory tests can still run.
func setup() {
	err := testutil.LoadEnv()
	if err != nil {
		fmt.Println(fmt.Errorf("LoadEnv: %w", err))
		return
	}

	_simpleDB, _gormDB, err = mysql.NewDB()
	if err != nil {
		fmt.Println("MySQL not available:", err)
		_simpleDB = nil
		_gormDB = nil
	}

	if _simpleDB == nil || _gormDB == nil {
		return
	}

	_ = _simpleDB.Drop(context.TODO(), _testTable)

	ddl := fmt.Sprintf(`CREATE TABLE %s (
		id int unsigned NOT NULL AUTO_INCREMENT,
		name varchar(50) NOT NULL DEFAULT '',
		category varchar(50) NOT NULL DEFAULT '',
		value int NOT NULL DEFAULT 0,
		PRIMARY KEY (id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`, _testTable)

	if err := _gormDB.Exec(ddl).Error; err != nil {
		panic(err)
	}

	if err := _gormDB.Exec(fmt.Sprintf(`INSERT INTO %s (name, category, value) VALUES
		('a', 'cat1', 10),
		('b', 'cat1', 20),
		('c', 'cat2', 30),
		('d', 'cat2', 40)`, _testTable)).Error; err != nil {
		panic(err)
	}
}

// teardown drops the test table.
func teardown() {
	if _simpleDB == nil {
		return
	}
	_ = _simpleDB.Drop(context.TODO(), _testTable)
}

func skipIfNoDB(t *testing.T) {
	t.Helper()
	if _gormDB == nil || _simpleDB == nil {
		t.Skip("MySQL not available")
	}
}
