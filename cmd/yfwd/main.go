package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"YH-FireWall/internal/handler"
)

func main() {
	configPath := flag.String("c", "/etc/yfw/config.toml", "Path to the configuration file")
	flag.Parse()

	ctx, cancel := signal.NotifyContext(
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
