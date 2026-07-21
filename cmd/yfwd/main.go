package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofrs/flock"

	"YH-FireWall/internal/handler"
)

func main() {
	configPath := flag.String("c", "/etc/yfw/config.toml", "Path to the configuration file")
	flag.Parse()

	// 防止多实例
	lock := flock.New("/var/run/yfwd.lock")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	locked, err := lock.TryLockContext(ctx, 100*time.Millisecond)
	if err != nil {
		log.Fatalf("Failed to acquire lock: %v", err)
	}
	if !locked {
		log.Fatal("Another instance is already running (cannot acquire /var/run/yfwd.lock)")
	}

	ctx, cancel = signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	h, err := handler.New(*configPath)
	if err != nil {
		log.Fatalf("Core service failed to start: %v\n", err)
	} else {
		log.Println("Core service started successfully")
	}

	<-ctx.Done()

	if err := h.Close(); err != nil {
		log.Printf("Core service failed to stop: %v", err)
	} else {
		log.Printf("Core service stopped successfully")
	}
}
