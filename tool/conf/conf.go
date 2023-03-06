package conf

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	goSimpleDb "github.com/auho/go-simple-db/v2"
	"github.com/pelletier/go-toml"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Db *DbConfig
}

type DbConfig struct {
	Driver string
	Dsn    string
}

func (dc *DbConfig) BuildDB() (*goSimpleDb.SimpleDB, error) {
	var db *goSimpleDb.SimpleDB
	var err error

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second,  // 慢 SQL 阈值
			LogLevel:                  logger.Error, // 日志级别
			IgnoreRecordNotFoundError: true,         // 忽略ErrRecordNotFound（记录未找到）错误
		},
	)

	dbc := &gorm.Config{
		Logger: newLogger,
	}

	switch dc.Driver {
	case "mysql":
		db, err = goSimpleDb.NewMysql(dc.Dsn, dbc)
	case "clickhouse":
		db, err = goSimpleDb.NewClickhouse(dc.Dsn, dbc)
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
