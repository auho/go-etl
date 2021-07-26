package lib

import (
	"fmt"
	"io/ioutil"

	"etl/lib/storage"
	"github.com/pelletier/go-toml"
)

type Config struct {
	DbSource storage.DbSourceConfig
	DbTarget storage.DbTargetConfig
}

type DbConfig struct {
	Driver string
	Dsn    string
	Scheme string
	Table  string
}

func LoadConfig(name string) *Config {
	filePath := fmt.Sprintf("conf/%s.toml", name)
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var tc Config
	err = toml.Unmarshal(fileContent, &tc)
	if err != nil {
		panic(err)
	}

	return &tc
}
