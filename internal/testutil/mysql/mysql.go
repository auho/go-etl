package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	simpledb "github.com/auho/go-simple-db/v3"
	testmysql "github.com/auho/go-toolkit-testutil/mysql"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dbName = "_test_etl"

// NewDB creates a *simpledb.SimpleDB and *gorm.DB from the TEST_MYSQL_DSN
// environment variable. It panics if TEST_MYSQL_DSN is unset or the connection fails.
// TEST_MYSQL_DSN must not include a database name; the test database is created
// automatically if it does not exist.
func NewDB() (*simpledb.SimpleDB, *gorm.DB, error) {
	// Build the full DSN with the database name.
	fullDSN, err := LoadDSN()
	if err != nil {
		return nil, nil, fmt.Errorf("LoadDSN: %w", err)
	}

	// Ensure the test database exists.
	initDB, err := sql.Open("mysql", fullDSN)
	if err != nil {
		return nil, nil, fmt.Errorf("sql.Open: %w", err)
	}

	defer func() { _ = initDB.Close() }()

	_, err = initDB.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	if err != nil {
		panic(err)
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Error,
			IgnoreRecordNotFoundError: true,
		},
	)

	simpleDB, gormDB, err := simpledb.NewMySQLGorm(fullDSN, &gorm.Config{Logger: newLogger})
	if err != nil {
		panic(err)
	}

	if sqlDB := simpleDB.SqlDB(); sqlDB != nil {
		conns := runtime.NumCPU() * 2
		sqlDB.SetMaxOpenConns(conns)
		sqlDB.SetMaxIdleConns(conns)
		sqlDB.SetConnMaxLifetime(5 * time.Minute)
	}

	return simpleDB, gormDB, nil
}

func LoadDSN() (string, error) {
	return testmysql.LoadDSN(dbName)
}
