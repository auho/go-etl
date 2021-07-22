package conf

import (
	"fmt"
	"io/ioutil"

	"github.com/pelletier/go-toml"
)

type TaskConfig struct {
	Name     string
	Action   ActionConfig
	DbSource DbSourceConfig
	DbTarget DbTargetConfig
}

type ActionConfig struct {
	Name          string
	MaxConcurrent int
}

type DbSourceConfig struct {
	MaxConcurrent int
	Size          int
	Page          int
	Driver        string
	Dsn           string
	Scheme        string
	Table         string
	PKeyName      string
	Fields        []string
}

type DbTargetConfig struct {
	MaxConcurrent int
	Size          int
	Driver        string
	Dsn           string
	Scheme        string
	Table         string
}

func LoadConfig(name string) TaskConfig {
	filePath := fmt.Sprintf("conf/%s.toml", name)
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var tc TaskConfig
	err = toml.Unmarshal(fileContent, &tc)
	if err != nil {
		panic(err)
	}

	return tc
}
