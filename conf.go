package go_etl

import (
	"fmt"
	"io/ioutil"

	"github.com/pelletier/go-toml"
)

type Config struct {
	dbConfig *DbConfig
}

type DbConfig struct {
	Driver string
	Dsn    string
}

func LoadConfig(name string) *Config {
	filePath := fmt.Sprintf("conf/%s.toml", name)
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var c Config
	err = toml.Unmarshal(fileContent, &c)
	if err != nil {
		panic(err)
	}

	return &c
}
