package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     string
	DBName   string
	Password string
	User     string
	SSLMode  string
}

func CreateDSN(cfg Config) (string, error) {
	if cfg.Host == "" || cfg.DBName == "" || cfg.Password == "" || cfg.Port == "" || cfg.User == "" {
		return "", fmt.Errorf("invalid db config")
	}

	// mysql
	// dsn := fmt.Sprintf(
	// 	"%s:%s@tcp(%s:%s)/%s?parseTime=true",
	// 	cfg.User,
	// 	cfg.Password,
	// 	cfg.Host,
	// 	cfg.Port,
	// 	cfg.DBName,
	// )

	// postgres
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
	)

	fmt.Println("dsn", dsn)

	return dsn, nil

}

func NewConnection(driver, dsn string) (*sql.DB, error) {
	conn, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(10)
	conn.SetConnMaxLifetime(10 * time.Minute)

	if err := conn.Ping(); err != nil {
		return nil, err
	}

	return conn, nil
}
