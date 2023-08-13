package server

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"sekawan-web/app/main/util"

	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgreSQLClientRepository struct {
	DB *gorm.DB
	TZ string
}

func NewPostgreSQLRepository(host, uname, pass, dbname string, port int, config *gorm.Config) (*PostgreSQLClientRepository, error) {
	tz := "Asia/Jakarta"

	if config == nil {
		config = &gorm.Config{}
	}

	maxIdleConn, errMaxIdleConn := strconv.Atoi(os.Getenv(util.CONFIG_DB_MAX_IDLE_CONNECTION))
	util.IsErrorDoPanic(errMaxIdleConn)
	maxOpenConn, errMaxOpenConn := strconv.Atoi(os.Getenv(util.CONFIG_DB_MAX_OPEN_CONNECTION))
	util.IsErrorDoPanic(errMaxOpenConn)
	lifetimeConn, errMaxLifetimeConn := time.ParseDuration(os.Getenv(util.CONFIG_DB_MAX_LIFETIME_CONNECTION))
	util.IsErrorDoPanic(errMaxLifetimeConn)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s", host, uname, pass, dbname, port, tz)
	sqlDB, err := sql.Open("pgx", dsn)
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(maxIdleConn)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(maxOpenConn)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(lifetimeConn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db.Use(otelgorm.NewPlugin())
	util.IsErrorDoPanicWithMessage("Cannot connect to PostgresSQL", err)

	if db == nil {
		panic("missing db")
	}

	return &PostgreSQLClientRepository{DB: db, TZ: tz}, nil
}
