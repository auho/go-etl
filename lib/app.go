package lib

import (
	"flag"

	"github.com/auho/go-simple-db/simple"
)

type App struct {
	config *Config
	db     simple.Driver
}

func NewApp() *App {
	a := &App{}

	return a
}

func (a *App) Run() {
	name := flag.String("name", "office", "")
	flag.Parse()

	a.config = LoadConfig(*name)
}
