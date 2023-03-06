package app

import (
	"testing"
)

func Test_App(t *testing.T) {
	app = NewApp("develop")
	if app.ConfName != "develop" {
		t.Error("app error")
	}

	err := app.DB.Ping()
	if err != nil {
		t.Error(err)
	}

	app.ConfName = "test"
}
