package main

import (
	"YH-FireWall/core"
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	configPath := flag.String("c", "/etc/yfw/config.yaml", "Path to the configuration file")
	flag.Parse()
	// syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGHUP
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGHUP)
	defer cancel()
	// 初始化核心服务
	core.Context = ctx
	core.Cancel = cancel
	// 启动核心服务
	if err := core.Start(*configPath); err != nil {
		log.Fatalf("Core service failed to start: %v", err)
	} else {
		log.Println("Core service started successfully")
	}
	// 阻塞主进程，等待信号
	<-ctx.Done()
	// 输出结束信息
	log.Println()
	// 关闭服务 // 必须是阻塞的，不然可能没清除就守护线程被关闭
	if err := core.Close(); err != nil {
		log.Printf("Core service failed to stop: %v", err)
	} else {
		log.Printf("Core service stopped successfully")
	}
}
