package main

import (
	"YH-FireWall/internal/core"
	"YH-FireWall/internal/server/unix"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	if len(os.Args) <= 1 {
		// 启动核心服务
		// syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGHUP
		ctx, stop := signal.NotifyContext(context.Background(),
			syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGHUP)
		defer stop()

		if err := core.Start(ctx); err != nil {
			log.Fatal("Core service failed to start:", err)
		} else {
			log.Println("Core service started successfully")
		}
		// 启动unix监听
		if err := unix.Start(); err != nil {
			log.Println("Unix service failed to start:", err)
		} else {
			log.Println("Unix service started successfully")
		}
		// 阻塞主进程，等待信号
		<-core.Done()
		// 关闭服务
		if core.IsRunning() {
			if err := core.Close(); err != nil {
				log.Println("Core service failed to stop:", err)
			}
		}
		if err := unix.Close(); err != nil {
			log.Println("Unix service failed to stop:", err)
		}
	} else {
		// 运行客户端
		conn, err := net.Dial("unix", "/tmp/firewall.sock")
		if err != nil {
			log.Println("Core service is not running.")
			log.Println("Please start it with: yfw start")
			os.Exit(1)
		}
		cmd := strings.Join(os.Args, " ")
		if _, err = conn.Write([]byte(cmd)); err != nil {
			log.Fatal("Failed to send command:", err)
		}
		result, err := io.ReadAll(conn)
		if err != nil {
			log.Fatal("Failed to read result:", err)
		}
		_ = conn.Close()

		fmt.Print(string(result))
	}
}
