package mysql

import (
	"log"
	"os"
	"runtime"
	"time"

	simpledb "github.com/auho/go-simple-db/v3"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDB creates a *simpledb.SimpleDB and *gorm.DB from the MYSQL_DSN
// environment variable. It panics if MYSQL_DSN is unset or the connection fails.
func NewDB() (*simpledb.SimpleDB, *gorm.DB) {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		panic("MYSQL_DSN environment variable is not set")
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Error,
			IgnoreRecordNotFoundError: true,
		},
	)

	simpleDB, gormDB, err := simpledb.NewMySQLGorm(dsn, &gorm.Config{Logger: newLogger})
	if err != nil {
		panic(err)
	}

	if sqlDB := simpleDB.SqlDB(); sqlDB != nil {
		conns := runtime.NumCPU() * 2
		sqlDB.SetMaxOpenConns(conns)
		sqlDB.SetMaxIdleConns(conns)
		sqlDB.SetConnMaxLifetime(5 * time.Minute)
	}

	return simpleDB, gormDB
}
