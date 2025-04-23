package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sshtn/mods/prompt"
	"sshtn/mods/storage"
	"syscall"
)

func main() {
	if err := run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	if err := storage.CreateStorage(); err != nil {
		return err
	}

	prompt := prompt.New()
	if err := prompt.Run(); err != nil {
		return err
	}
	return nil
}
