package conf

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"time"

	simpledb "github.com/auho/go-simple-db/v2"
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

func (dc *DbConfig) BuildDB() (*simpledb.SimpleDB, *gorm.DB, error) {
	var simpleDB *simpledb.SimpleDB
	var gromDB *gorm.DB
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
		simpleDB, gromDB, err = simpledb.NewMySQLGorm(dc.Dsn, dbc)
		if err != nil {
			err = fmt.Errorf("NewMySQLGorm: %w", err)
		}
	case "clickhouse":
		simpleDB, gromDB, err = simpledb.NewClickHouseGorm(dc.Dsn, dbc)
		if err != nil {
			err = fmt.Errorf("NewClickHouseGorm: %w", err)
		}
	default:
		err = fmt.Errorf("driver[%s] not found", dc.Driver)
	}

	if simpleDB != nil {
		sqldb := simpleDB.SqlDB()
		if sqldb != nil {
			conns := runtime.NumCPU() * 2
			sqldb.SetMaxOpenConns(conns)
			sqldb.SetMaxIdleConns(conns)
			sqldb.SetConnMaxLifetime(5 * time.Minute)
		}
	}

	return simpleDB, gromDB, err
}

func LoadConfig(dir string, name string) (*Config, error) {
	filePath := path.Join(dir, fmt.Sprintf("%s.toml", name))
	fileContent, err := os.ReadFile(filePath)
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
