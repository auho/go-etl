package go_etl

import (
	"io/ioutil"
	"os"
	"testing"
)

var testConfigContent = `[db]
dsn = "test:test@tcp(127.0.0.1:3306)/test"
driver = "mysql"
`
var app *App

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setup() {
	err := os.Mkdir("conf", 0700)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("conf/office.toml", []byte(testConfigContent), 0600)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("conf/test.toml", []byte(testConfigContent), 0600)
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
