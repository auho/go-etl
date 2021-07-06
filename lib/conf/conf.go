package conf

import (
	"fmt"
	"io/ioutil"

	"github.com/pelletier/go-toml"
)

type TaskConfig struct {
	Name string
}

type Action struct {
	Name          string
	MaxConcurrent int
}

type SourceDB struct {
	MaxConcurrent int
	Size          int
	Page          int
	Driver        string
	Dsn           string
	Scheme        string
	Table         string
}

type TargetDB struct {
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
