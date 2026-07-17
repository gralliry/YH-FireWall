package main

import (
	"YH-FireWall/internal/handler"
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "core" {
		runCore()
	} else {
		runClient()
	}
}

func runCore() {
	configPath := flag.String("c", "/etc/yfw/config.toml", "Path to the configuration file")
	flag.CommandLine.Parse(os.Args[2:])

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

func runClient() {
	socketPath := flag.String("s", "/tmp/yfw.sock", "Path to the socket file")
	flag.Parse()

	conn, err := net.Dial("unix", *socketPath)
	if err != nil {
		log.Println("Core service is not running.")
		log.Println("Please start it with: yfw core")
		os.Exit(1)
	}
	defer conn.Close()

	client := bufio.NewReadWriter(
		bufio.NewReader(conn),
		bufio.NewWriter(conn),
	)

	// 非交互模式：直接发送参数并退出
	args := flag.Args()

	cmd := strings.Join(args, " ") + "\n"
	if _, err = client.WriteString(cmd); err != nil {
		log.Fatal("Failed to send command:", err)
	}
	if err = client.Flush(); err != nil {
		log.Fatal("Failed to send command:", err)
	}
	result, err := client.ReadString(0)
	if err != nil {
		log.Fatal("Failed to read result:", err)
	}
	fmt.Print(result)
}
