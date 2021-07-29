package go_etl

import (
	"flag"
	"fmt"
	"io/fs"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/auho/go-simple-db/simple"
)

type App struct {
	DbConfig *DbConfig
	Db       simple.Driver
	Demand   *AppDemand
	WorkDir  string
	ConfName string
}

func NewApp() *App {
	a := &App{}
	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	workDir, err = filepath.Abs(workDir)
	if err != nil {
		panic(err)
	}

	a.WorkDir = workDir

	name := flag.String("config", "office", "")
	flag.Parse()

	a.ConfName = *name

	config, err := LoadConfig(a.ConfName)
	if err != nil {
		panic(err)
	}

	a.DbConfig = config.Db

	a.Db, err = simple.NewDriver(a.DbConfig.Driver, a.DbConfig.Dsn)
	if err != nil {
		panic(err)
	}

	a.Demand = newAppDemand(a)

	return a
}

type demandState struct {
	file      string
	startTime time.Time
	endTime   time.Time
}

func newDemandState(file string) demandState {
	s := demandState{}
	s.startTime = time.Now()
	return s
}

func (ds demandState) end() {
	ds.endTime = time.Now()
}

func (ds demandState) state() string {
	d := ds.endTime.Sub(ds.startTime)
	return fmt.Sprintf("%s %f 分， %f 秒", ds.file, math.Floor(d.Minutes()), math.Floor(d.Seconds()))
}

type AppDemand struct {
	app    *App
	states []demandState
}

func newAppDemand(app *App) *AppDemand {
	ad := &AppDemand{}
	ad.app = app
	ad.states = make([]demandState, 0)

	return ad
}

func (ad *AppDemand) RunPaths(paths []string) error {
	for _, path := range paths {
		err := ad.RunPath(path)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ad *AppDemand) RunPath(path string) error {
	path = ad.absDir(path)

	ad.log("PATH:: " + path)

	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		out, err := ad.RunFile(path)
		fmt.Println(out)

		if err != nil {
			return nil
		}

		return nil
	})

	return err
}

func (ad *AppDemand) RunPathFiles(path string, files []string) error {
	path = ad.absDir(path)
	for _, file := range files {
		file = path + string(filepath.Separator) + file

		out, err := ad.RunFile(file)
		fmt.Println(out)

		if err != nil {
			return err
		}
	}

	return nil
}

func (ad *AppDemand) RunFile(file string) (string, error) {
	return ad.run(file)
}

func (ad *AppDemand) run(file string) (string, error) {
	if file[len(file)-3:] != ".go" {
		file = file + ".go"
	}

	s := newDemandState(file)
	out, err := exec.Command("go", "run", file, "--config="+ad.app.ConfName).Output()
	s.end()

	return string(out), err
}

func (ad *AppDemand) absDir(path string) string {
	if filepath.IsAbs(path) {
		return path
	}

	return ad.app.WorkDir + string(filepath.Separator) + path
}

func (ad *AppDemand) log(s string) {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05\n"), s)
}

func (ad *AppDemand) State() {
	for _, s := range ad.states {
		fmt.Println(s.state())
	}
}
