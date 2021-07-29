package main

import (
	"fmt"

	go_etl "github.com/auho/go-etl"
)

func main() {
	app := go_etl.NewApp()
	fmt.Println("b", app.ConfName)
}
