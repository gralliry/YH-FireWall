package main

import (
	"YH-FireWall/core"
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
		// syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGHUP
		ctx, stop := signal.NotifyContext(context.Background(),
			syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGHUP)
		defer stop()
		// 初始化核心服务
		core.Context = ctx
		core.Cancel = stop
		// 启动核心服务
		if err := core.Start(); err != nil {
			log.Fatalf("Core service failed to start: %v", err)
		} else {
			log.Println("Core service started successfully")
		}
		// 阻塞主进程，等待信号
		<-ctx.Done()
		// 输出结束信息
		fmt.Println()
		// 关闭服务 // 必须是阻塞的，不然可能没清除就守护线程被关闭
		if err := core.Close(); err != nil {
			log.Printf("Core service failed to stop: %v", err)
		} else {
			log.Printf("Core service stopped successfully")
		}
	} else {
		// 运行客户端
		conn, err := net.Dial("unix", "/tmp/yfw.sock")
		if err != nil {
			log.Println("Core service is not running.")
			log.Println("Please start it with: yfw")
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
