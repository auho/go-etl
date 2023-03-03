package go_etl

import (
	"fmt"
	"io/ioutil"

	go_simple_db "github.com/auho/go-simple-db/v2"
	"github.com/pelletier/go-toml"
)

type Config struct {
	Db *DbConfig
}

type DbConfig struct {
	Driver string
	Dsn    string
}

func (dc *DbConfig) BuildDB() (*go_simple_db.SimpleDB, error) {
	var db *go_simple_db.SimpleDB
	var err error
	switch dc.Driver {
	case "mysql":
		db, err = go_simple_db.NewMysql(dc.Dsn)
	case "clickhouse":
		db, err = go_simple_db.NewClickhouse(dc.Dsn)
	default:
		err = fmt.Errorf("driver[%s] not found", dc.Driver)
	}

	if err != nil {
		err = fmt.Errorf("driver[%s] [%s] build error", dc.Driver, dc.Dsn)
	}

	return db, err
}

func LoadConfig(name string) (*Config, error) {
	filePath := fmt.Sprintf("conf/%s.toml", name)
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var c Config
	err = toml.Unmarshal(fileContent, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
