package app

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/auho/go-etl/v2/insight/app/conf"
	simpleDb "github.com/auho/go-simple-db/v2"
)

var APP *Application

func NewApp(cn string) {
	APP = NewApplication(cn)
}

type Application struct {
	Run
	DB       *simpleDb.SimpleDB
	Name     string
	ConfName string
	WorkDir  string
	DataDir  string
	XlsxDir  string
	ConfDir  string
}

func NewApplication(cn string) *Application {
	a := &Application{}
	a.ConfName = cn

	a.buildWorkDir()
	a.checkDir()
	a.buildDb()

	return a
}

func (a *Application) buildWorkDir() {
	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	workDir, err = filepath.Abs(workDir)
	if err != nil {
		panic(err)
	}

	a.Name = filepath.Base(workDir)
	a.WorkDir = workDir
	a.DataDir = path.Join(a.WorkDir, "data")
	a.XlsxDir = path.Join(a.WorkDir, "xlsx")
	a.ConfDir = path.Join(a.WorkDir, "conf")
}

func (a *Application) buildDb() {
	config, err := conf.LoadConfig(a.ConfDir, a.ConfName)
	if err != nil {
		panic(err)
	}

	a.DB, err = config.Db.BuildDB()
	if err != nil {
		panic(err)
	}
}

func (a *Application) checkDir() {
	for _, _dir := range []string{a.DataDir, a.XlsxDir} {
		_, err := os.Stat(_dir)
		if os.IsNotExist(err) {
			err = os.Mkdir(_dir, 0744)
			if err != nil {
				panic(fmt.Errorf("dir[%s]; %w", _dir, err))
			}
		}
	}
}

func (a *Application) State() []string {
	return []string{
		"name: " + a.Name,
		"conf name: " + a.ConfName,
		"data dir: " + a.DataDir,
		"xlsx dir: " + a.XlsxDir,
	}
}

func (a *Application) PrintlnState() {
	ss := a.State()
	for _, _s := range ss {
		fmt.Println(_s)
	}

	fmt.Println()
}
