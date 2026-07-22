package app

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/auho/go-etl/v3/internal/testutil"
	"github.com/auho/go-etl/v3/internal/testutil/mysql"
)

var app *Application
var confDir string

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setup() {
	err := testutil.LoadEnv()
	if err != nil {
		panic(err)
	}

	dsn, err := mysql.LoadDSN()
	if err != nil {
		panic(err)
	}

	if dsn == "" {
		fmt.Println("skip: mysql dsn not set")
		os.Exit(0)
	}

	testConfigContent := fmt.Sprintf(`[db]
dsn = "%s"
driver = "mysql"
`, dsn)

	confDir = filepath.Join(testutil.ProjectRoot(), "/conf")
	_, err = os.Stat(confDir)
	if err != nil {
		err = os.Mkdir(confDir, 0700)
		if err != nil {
			panic(fmt.Sprintf("Mkdir[%s]: %v", confDir, err))
		}
	}

	confPath := filepath.Join(confDir, "develop.toml")
	err = os.WriteFile(confPath, []byte(testConfigContent), 0600)
	if err != nil {
		panic(err)
	}
}

func tearDown() {
	err := os.RemoveAll(confDir)
	if err != nil {
		panic(err)
	}
}
