package go_etl

import (
	"fmt"
	"testing"
)

func Test_App(t *testing.T) {
	app = NewApp()
	if app.DbConfig.Driver != "mysql" {
		t.Error("app error")
	}

	if app.ConfName != "office" {
		t.Error("app error")
	}

	err := app.Db.Ping()
	if err != nil {
		t.Error(err)
	}

	app.ConfName = "test"
}

func TestAppDemand_RunPaths(t *testing.T) {
	err := app.Demand.RunPaths([]string{
		"example/demand/a",
		"example/demand/a1",
	})

	if err != nil {
		t.Error(err)
	}
}

func TestAppDemand_RunPath(t *testing.T) {
	err := app.Demand.RunPath("example/demand")
	if err != nil {
		t.Error(err)
	}
}

func TestAppDemand_RunPathFiles(t *testing.T) {
	err := app.Demand.RunPathFiles("example/demand/a", []string{"one", "two"})
	if err != nil {
		t.Error(err)
	}
}

func TestAppDemand_Run(t *testing.T) {
	out, err := app.Demand.RunFile("example/demand/b")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(out)
}
