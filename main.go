package main

import (
	"YH-FireWall/core"
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: yfw core|cli")
	}
	switch os.Args[1] {
	case "core":
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
	case "cli":
		// 运行客户端
		conn, err := net.Dial("unix", "/tmp/yfw.sock")
		if err != nil {
			log.Println("Core service is not running.")
			log.Println("Please start it with: yfw")
			os.Exit(1)
		}

		reader := bufio.NewReader(os.Stdin)
		server := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

		for {
			fmt.Print("> ")
			input, err := reader.ReadBytes('\n') // 读到回车为止
			if err != nil {
				log.Fatal("Failed to read command:", err)
			}
			// 这里会把换行符写入，作为服务端读取的结束符
			if _, err = server.Write(input); err != nil {
				log.Fatal("Failed to send command:", err)
			}
			if err = server.Flush(); err != nil {
				log.Fatal("Failed to send command:", err)
			}
			// 0 作为客户端读取结束符，服务端会返回一个0作为结束符，这里会读取到这个结束符，然后结束读取
			result, err := server.ReadString(0)
			if err != nil {
				log.Fatal("Failed to read result:", err)
			}
			fmt.Print(result)
		}
	default:
		log.Println("Usage: yfw core|cli")
	}
}
