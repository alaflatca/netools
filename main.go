package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sshtn/mods/prompt"
	"sshtn/mods/storage"
	"syscall"
	"time"
)

func main() {
	if err := run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "netools: %s", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	go monitorMemory(ctx)

	if err := storage.CreateStorage(); err != nil {
		return err
	}

	if err := prompt.Run(ctx); err != nil {
		return err
	}

	return nil
}

func monitorMemory(ctx context.Context) {
	var m runtime.MemStats
	var ticker = time.NewTicker(time.Second * 2)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("monitorMemory 종료")
			return
		case <-ticker.C:
			runtime.ReadMemStats(&m)
			fmt.Printf("Alloc = %v MiB\n", m.Alloc/1024/1024)
			fmt.Printf("\tTotalAlloc = %v MiB\n", m.TotalAlloc/1024/1024)
			fmt.Printf("\tSys = %v MiB\n", m.Sys/1024/1024)
			fmt.Printf("\tNumGC = %v \n", m.NumGC)
			fmt.Println()
		}
	}
}
