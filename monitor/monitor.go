package trace

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

func MonitorMemory(ctx context.Context) {
	var m runtime.MemStats
	var ticker = time.NewTicker(time.Second * 2)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("MonitorMemory 종료")
			return
		case <-ticker.C:
			runtime.ReadMemStats(&m)
			fmt.Printf("Alloc = %v MiB\n", m.Alloc/1024/1024)
			fmt.Printf("\tTotalAlloc = %v MiB\n", m.TotalAlloc/1024/1024)
			fmt.Printf("\tSys = %v MiB\n", m.Sys/1024/1024)
			fmt.Printf("\tNumGC = %v\n", m.NumGC)
		}
	}
}
