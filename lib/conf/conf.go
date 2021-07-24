package conf

import (
	"fmt"
	"io/ioutil"

	"etl/lib/storage"
	"github.com/pelletier/go-toml"
)

type TaskConfig struct {
	Name     string
	Action   ActionConfig
	DbSource storage.DbSourceConfig
	DbTarget storage.DbTargetConfig
}

type ActionConfig struct {
	Name          string
	MaxConcurrent int
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
