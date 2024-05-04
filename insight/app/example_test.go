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
	_ = _app.RunCommandE(nil, []string{}, &cobra.Command{})
}
