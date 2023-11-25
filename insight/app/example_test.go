package app

import (
	"github.com/spf13/cobra"
)

func ExampleNewApp() {
	NewApp()
}

func ExampleNewApplication() {
	_app := NewApplication()

	// build
	_app.Build("config name")

	// path
	_app.DataFilePath("file name")
	_app.XlsxFilePath("file name")

	// run command
	_app.AddCommand(&cobra.Command{})
	_ = _app.RunCommandE([]string{})
}
