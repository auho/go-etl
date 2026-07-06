package app

import (
	"context"
	"testing"
)

func Test_App(t *testing.T) {
	app = NewApplication()
	app.Build("develop")
	if app.ConfName != "develop" {
		t.Error("app error")
	}

	err := app.DB.Ping(context.TODO())
	if err != nil {
		t.Error(err)
	}

	app.ConfName = "test"
}
