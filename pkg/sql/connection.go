package sql

import (
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type sqlClient struct {
	cfg *SQLConfig
}

func NewClient(cfg *SQLConfig) (*sqlx.DB, error) {
	sqlCli := sqlClient{
		cfg: cfg,
	}

	return sqlCli.newConn()
}

func (s *sqlClient) newConn() (*sqlx.DB, error) {
	dsn := s.constructDSN()
	if s.cfg.DatabaseDriver == "postgres" {
		s.cfg.DatabaseDriver = "pgx"
	}

	db, err := sqlx.Open(s.cfg.DatabaseDriver, dsn)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// DatabaseSetMaxIdleConnection sets the maximum number of idle connections in the pool.
	// Default: 2. Helps reduce connection churn. Increase for better reuse under load.
	if s.cfg.DatabaseMaxIdleConnection != 0 {
		db.SetMaxIdleConns(s.cfg.DatabaseMaxIdleConnection)
	}

	// DatabaseSetMaxOpenConnection sets the maximum number of open connections to the database.
	// Default: 0 (unlimited). Recommended to set in production to prevent resource exhaustion.
	db.SetMaxOpenConns(s.cfg.DatabaseMaxOpenConnection)

	// DatabaseSetMaxIdleTimeConnection sets the maximum amount of time a connection may remain idle.
	// Default: 0 (no idle timeout). Recommended: ~5m (300 seconds) to clean up stale connections.
	// Value is expected in seconds.
	db.SetConnMaxIdleTime(time.Duration(s.cfg.DatabaseMaxIdleTimeConnection))

	// DatabaseSetMaxLifetimeConnection sets the maximum total lifetime of a connection.
	// Default: 0 (connections live forever). Recommended: ~30m (1800 seconds) to prevent DB timeouts or leaks.
	// Value is expected in seconds.
	db.SetConnMaxLifetime(time.Duration(s.cfg.DatabaseMaxLifetimeConnection))

	return db, nil
}

func (s *sqlClient) newMySQLDataSourceName() string {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%v",
		s.cfg.DatabaseUsername, s.cfg.DatabasePassword, s.cfg.DatabaseHost, s.cfg.DatabasePort,
		s.cfg.DatabaseName, s.cfg.DatabaseCharset, s.cfg.DatabaseParsetime,
	)

	return dsn
}

func (s *sqlClient) newPostgresDataSourceName() string {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		s.cfg.DatabaseHost, s.cfg.DatabaseUsername, s.cfg.DatabasePassword,
		s.cfg.DatabaseName, s.cfg.DatabasePort,
	)

	return dsn
}

func (s *sqlClient) constructDSN() string {
	switch strings.ToLower(s.cfg.DatabaseDriver) {
	case "postgres":
		return s.newPostgresDataSourceName()
	case "mysql":
		return s.newMySQLDataSourceName()
	default:
		panic("please specify your own db driver")
	}
}
