package application

import (
	"github.com/auho/go-etl/v2/insight/app"
)

var App *app.Application

func NewApp(cn string) {
	App = app.NewApplication(cn)
}
