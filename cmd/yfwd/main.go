package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"YH-FireWall/internal/handler"
)

func main() {
	// 检测当前用户是否为 root 用户
	if os.Geteuid() != 0 {
		log.Fatalf("current user is not root")
	}

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
