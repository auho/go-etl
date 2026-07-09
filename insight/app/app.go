package app

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/auho/go-etl/v3/insight/app/conf"
	simpledb "github.com/auho/go-simple-db/v3"
	"gorm.io/gorm"
)

var App *Application

func NewApp() {
	App = NewApplication("")
}

type Application struct {
	Run
	Xlsx
	DB       *simpledb.SimpleDB
	GormDB   *gorm.DB
	Name     string
	ConfName string
	WorkDir  string
	DataDir  string
	XlsxDir  string
	ConfDir  string
}

func NewApplication(dir string) *Application {
	a := &Application{}

	a.buildWorkDir(dir)
	a.checkDir()

	return a
}

func (a *Application) buildWorkDir(dir string) {
	workDir := dir
	if workDir == "" {
		_wd, err := os.Getwd()
		if err != nil {
			panic(fmt.Errorf("os.Getwd: %w", err))
		}

		workDir, err = filepath.Abs(_wd)
		if err != nil {
			panic(fmt.Errorf("filepath.Abs: %w", err))
		}
	}

	a.Name = filepath.Base(workDir)
	a.WorkDir = workDir
	a.DataDir = path.Join(a.WorkDir, "data")
	a.XlsxDir = path.Join(a.WorkDir, "xlsx")
	a.ConfDir = path.Join(a.WorkDir, "conf")

	a.Xlsx.XlsxDir = a.XlsxDir
}

func (a *Application) checkDir() {
	for _, _dir := range []string{a.DataDir, a.XlsxDir, a.ConfDir} {
		_, err := os.Stat(_dir)
		if os.IsNotExist(err) {
			err = os.Mkdir(_dir, 0744)
			if err != nil {
				panic(fmt.Errorf("mkdir[%s]: %w", _dir, err))
			}
		}
	}
}

func (a *Application) Build(cn string) {
	a.ConfName = cn

	config, err := conf.LoadConfig(a.ConfDir, a.ConfName)
	if err != nil {
		a.PrintlnState()

		panic(err)
	}

	a.DB, a.GormDB, err = config.DB.BuildWithGorm()
	if err != nil {
		a.PrintlnState()

		panic(err)
	}
}

// DataFilePath
// name with file suffix
func (a *Application) DataFilePath(name string) string {
	return path.Join(a.DataDir, name)
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
