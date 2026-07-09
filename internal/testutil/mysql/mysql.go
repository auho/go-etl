package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"

	simpledb "github.com/auho/go-simple-db/v3"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dbName = "_test_etl"

// NewDB creates a *simpledb.SimpleDB and *gorm.DB from the MYSQL_DSN
// environment variable. It panics if MYSQL_DSN is unset or the connection fails.
// MYSQL_DSN must not include a database name; the test database is created
// automatically if it does not exist.
func NewDB() (*simpledb.SimpleDB, *gorm.DB) {
	// Build the full DSN with the database name.
	fullDSN, err := SetupDSN()
	if err != nil {
		panic(fmt.Errorf("SetupDSN: %w", err))
	}

	// Ensure the test database exists.
	initDB, err := sql.Open("mysql", fullDSN)
	if err != nil {
		panic(err)
	}
	defer initDB.Close()

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

	return simpleDB, gormDB
}

func LoadDSN() (string, error) {
	baseDSN := os.Getenv("MYSQL_DSN")
	if baseDSN == "" {
		panic("MYSQL_DSN environment variable is not set")
	}

	// Ensure baseDSN ends with "/" for MySQL driver compatibility.
	if !strings.HasSuffix(baseDSN, "/") {
		baseDSN += "/"
	}

	return baseDSN, nil
}

func SetupDSN() (string, error) {
	dsn, err := LoadDSN()
	if err != nil {
		return "", fmt.Errorf("LoadDSN: %w", err)
	}

	if dsn == "" {
		return "", errors.New("MYSQL_DSN not set; create .env.test with MYSQL_DSN")
	}

	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		return "", fmt.Errorf("ParseDSN: %w", err)
	}

	cfg.DBName = dbName

	return cfg.FormatDSN(), nil
}
