package load

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/auho/go-etl/v3/insight/assistant/excel/writer"
	"github.com/auho/go-etl/v3/internal/testutil"
	"github.com/auho/go-etl/v3/internal/testutil/mysql"
	simpledb "github.com/auho/go-simple-db/v3"
)

var _simpleDB *simpledb.SimpleDB

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setUp() {
	err := testutil.LoadEnv()
	if err != nil {
		fmt.Println(fmt.Errorf("LoadEnv: %w", err))
		return
	}

	// mysql.NewDB panics on connection failure; recover so unit tests can still run
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("MySQL not available:", r)
			_simpleDB = nil
		}
	}()

	_simpleDB, _, err = mysql.NewDB()
	if err != nil {
		fmt.Println("MySQL not available:", err)
		_simpleDB = nil
	}
}

func tearDown() {
	if _simpleDB == nil {
		return
	}
}

func skipNoDB(t *testing.T) {
	t.Helper()
	if _simpleDB == nil {
		t.Skip("MySQL not available")
	}
}

// createTestExcel creates a temporary xlsx file with the given sheet data.
func createTestExcel(t *testing.T, sheetName string, rows [][]any) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "test.xlsx")

	e, err := writer.NewExcel(path)
	if err != nil {
		t.Fatalf("writer.NewExcel: %v", err)
	}

	_, err = e.NewSheetWithData(sheetName, rows)
	if err != nil {
		t.Fatalf("NewSheetWithData: %v", err)
	}

	err = e.SaveAs()
	if err != nil {
		t.Fatalf("SaveAs: %v", err)
	}

	err = e.Close()
	if err != nil {
		t.Fatalf("Close: %v", err)
	}

	return path
}

func dropTable(t *testing.T, tableName string) {
	t.Helper()
	if _simpleDB == nil {
		return
	}
	err := _simpleDB.Drop(context.TODO(), tableName)
	if err != nil {
		t.Logf("drop table %s: %v", tableName, err)
	}
}
