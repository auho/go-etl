package testutil

import "github.com/auho/go-toolkit-testutil"

func LoadEnv() error {
	return testutil.LoadEnv()
}

func ProjectRoot() string {
	return testutil.ProjectRoot()
}
