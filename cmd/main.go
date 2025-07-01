package main

import (
	"context"
	"fmt"
	"netools/internal/db"
	"netools/tui"
	"os"
	"os/signal"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
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
	db, err := db.New(ctx, dbname)
	if err != nil {
		return err
	}
	defer db.Close()

	// sqlite 테이블 생성
	if err := db.CreateTable(ctx); err != nil {
		return err
	}

	// tui 시작
	p := tea.NewProgram(tui.NewProgramModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}
