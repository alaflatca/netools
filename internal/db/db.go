package db

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

//go:embed schema/schema.sql
var schemaSQL string

type Querier interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

var conn *sql.DB

func Init(ctx context.Context, path string) error {
	if conn != nil {
		return nil
	}
	var err error
	conn, err = sql.Open("sqlite", path)
	if err != nil {
		return fmt.Errorf("failed to open db: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()
	if err := conn.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping db: %v", err)
	}
	return nil
}

func Get() Querier {
	if conn != nil {
		return conn
	}
	panic("DB Not Initialized")
}

func Close() error {
	if conn == nil {
		return nil
	}
	return conn.Close()
}

// 나중에 sqlite_ssh.go로 이동?
func Migrate(ctx context.Context) error {
	_, err := conn.ExecContext(ctx, schemaSQL)
	if err != nil {
		return fmt.Errorf("failed to execute schema: %v", err)
	}
	return nil
}
