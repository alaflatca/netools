package main

import (
	"context"
	"fmt"
	"netools/internal/db"
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

	// database Init
	dbname := "./tmp/netools.db"
	db, err := db.New(dbname)
	if err != nil {
		return err
	}
	defer db.Stop()

	// logging Init
	// logging.New()

	// tui Init
	if err := tui.Start(ctx, db); err != nil {
		return err
	}

	// *mods --> internal/*tool/api

	return nil
}
