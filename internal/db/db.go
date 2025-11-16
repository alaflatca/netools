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

type DB struct {
	conn *sql.DB
}

func New(path string) (*DB, error) {
	db := &DB{}
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %v", err)
	}
	db.conn = conn

	return db, nil
}

func (d *DB) Start(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	if err := d.conn.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping db: %v", err)
	}

	if err := d.Migrate(ctx); err != nil {
		return fmt.Errorf("failed to migrate db: %v", err)
	}

	return nil
}

func (d *DB) Stop() {
	if d.conn != nil {
		if err := d.conn.Close(); err != nil {

		}
	}
}

// 나중에 sqlite_ssh.go로 이동?
func (d *DB) Migrate(ctx context.Context) error {
	_, err := d.conn.ExecContext(ctx, schemaSQL)
	if err != nil {
		return fmt.Errorf("failed to execute schema: %v", err)
	}
	return nil
}
