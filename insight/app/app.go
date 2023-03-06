package app

import (
	"os"
	"path/filepath"

	"github.com/auho/go-etl/tool/conf"
	goSimpleDb "github.com/auho/go-simple-db/v2"
)

type App struct {
	DB       *goSimpleDb.SimpleDB
	WorkDir  string
	ConfName string
}

func NewApp(cn string) *App {
	a := &App{}
	a.ConfName = cn

	a.workDir()
	a.db()

	return a
}

func (a *App) db() {
	config, err := conf.LoadConfig(a.ConfName)
	if err != nil {
		panic(err)
	}

	a.DB, err = config.Db.BuildDB()
	if err != nil {
		panic(err)
	}
}

func (a *App) workDir() {
	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	workDir, err = filepath.Abs(workDir)
	if err != nil {
		panic(err)
	}

	a.WorkDir = workDir
}

func (a *App) state() []string {
	return []string{
		"conf name: " + a.ConfName,
		"work dir:" + a.WorkDir,
	}
}
