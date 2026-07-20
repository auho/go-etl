package app

import (
	"context"
	"fmt"
	"os"
	"testing"

	testutil "github.com/auho/go-toolkit-testutil"
)

func TestApp(t *testing.T) {
	pr := testutil.ProjectRoot()

	app = NewApplication(pr)
	app.Build("develop")
	if app.ConfName != "develop" {
		t.Error("app error")
	}

	err := app.DB.Ping(context.TODO())
	if err != nil {
		t.Error(err)
	}

	app.ConfName = "test"

	checkDirFn := func(dir string) (bool, error) {
		fi, err := os.Stat(dir)
		if err == nil {
			return fi.IsDir(), nil
		}

		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	removeDirFn := func(t *testing.T, dir string) {
		ok, err := checkDirFn(dir)
		if err != nil {
			t.Error(fmt.Errorf("checkDirFn: %w", err))
		}

		if ok {
			err = os.RemoveAll(dir)
			if err != nil {
				t.Error(fmt.Errorf("os.RemoveAll: %w", err))
			}
		}
	}

	removeDirFn(t, app.DataDir)
	removeDirFn(t, app.XlsxDir)
	removeDirFn(t, app.ConfDir)
}
