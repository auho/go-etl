package app

import (
	"fmt"
	"os"
	"testing"

	"github.com/auho/go-etl/v3/internal/testutil"
)

var app *Application

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setup() {
	testutil.LoadEnv()
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		fmt.Println("skip: MYSQL_DSN not set")
		os.Exit(0)
	}
	testConfigContent := fmt.Sprintf(`[db]
dsn = "%s"
driver = "mysql"
`, dsn)
	_, err := os.Stat("conf")
	if err != nil {
		err = os.Mkdir("conf", 0700)
		if err != nil {
			panic(err)
		}
	}

	err = os.WriteFile("conf/office.toml", []byte(testConfigContent), 0600)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("conf/test.toml", []byte(testConfigContent), 0600)
	if err != nil {
		panic(err)
	}
}

func tearDown() {
	err := os.RemoveAll("conf")
	if err != nil {
		panic(err)
	}
}
