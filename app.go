package go_etl

import (
	"flag"

	"github.com/auho/go-simple-db/simple"
)

type App struct {
	DbConfig *DbConfig
	Db       simple.Driver
}

func NewApp() *App {
	a := &App{}

	name := flag.String("config", "office", "")
	flag.Parse()

	config := LoadConfig(*name)
	a.DbConfig = config.dbConfig

	var err error
	a.Db, err = simple.NewDriver(a.DbConfig.Driver, a.DbConfig.Dsn)
	if err != nil {
		panic(err)
	}

	return a
}
