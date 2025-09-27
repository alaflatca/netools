package main

import (
	"context"
	"fmt"
	db "netools/internal/database"
	"netools/internal/tui"
	"os"
	"os/signal"
	"syscall"
)

// var dbPath flag.String  ("dbPath","./tmp/netools.db","sqlite storage path")

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "netools: %s", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// sqlite 생성
	dbname := "./tmp/netools.db"
	err := db.Init(ctx, dbname)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := tui.Start(ctx); err != nil {
		return err
	}

	return nil
}
