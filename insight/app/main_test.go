package app

import (
	"os"
	"testing"
)

var testConfigContent = `[db]
dsn = "test:Test123$@tcp(127.0.0.1:3306)/test"
driver = "mysql"
`
var app *Application

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setup() {
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
