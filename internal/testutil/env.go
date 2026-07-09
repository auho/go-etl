package testutil

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

// ProjectRoot walks up from this file's location to find the directory
// containing go.mod.
func ProjectRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

// LoadEnv loads environment variables from the .env file at the project root.
// It is a no-op if the file does not exist; variables already set in the
// environment take precedence and are not overwritten.
func LoadEnv() (string, error) {
	projectRoot := ProjectRoot()
	err := godotenv.Load(filepath.Join(projectRoot, ".env.test"))
	if err != nil {
		return "", fmt.Errorf("load[env.test]: %w", err)
	}

	return projectRoot, nil
}
