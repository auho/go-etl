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
	DB *DBConfig
}

type DBConfig struct {
	Driver string
	DSN    string
}

func (dc *DBConfig) BuildDB() (*simpledb.SimpleDB, *gorm.DB, error) {
	var simpleDB *simpledb.SimpleDB
	var gormDB *gorm.DB
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
		simpleDB, gormDB, err = simpledb.NewMySQLGorm(dc.DSN, dbc)
		if err != nil {
			err = fmt.Errorf("NewMySQLGorm: %w", err)
		}
	case "clickhouse":
		simpleDB, gormDB, err = simpledb.NewClickHouseGorm(dc.DSN, dbc)
		if err != nil {
			err = fmt.Errorf("NewClickHouseGorm: %w", err)
		}
	default:
		err = fmt.Errorf("driver[%s] not found", dc.Driver)
	}

	if simpleDB != nil {
		sqlDB := simpleDB.SqlDB()
		if sqlDB != nil {
			conns := runtime.NumCPU() * 2
			sqlDB.SetMaxOpenConns(conns)
			sqlDB.SetMaxIdleConns(conns)
			sqlDB.SetConnMaxLifetime(5 * time.Minute)
		}
	}

	return simpleDB, gormDB, err
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
