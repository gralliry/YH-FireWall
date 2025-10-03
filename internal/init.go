package internal

import (
	"YH-FireWall/internal/core"
	"YH-FireWall/internal/repo"
	"YH-FireWall/internal/web"
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Start() {
	// 参数
	var (
		address  = flag.String("a", "0.0.0.0:80", "Web server address")
		username = flag.String("u", "", "Web server username")
		password = flag.String("p", "", "Web server password")
	)
	flag.Parse()
	// 检测当前用户是否为 root 用户
	if os.Geteuid() != 0 {
		log.Fatal("当前用户非 root 用户")
	}
	// syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGHUP
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGHUP)
	defer stop()

	if err := repo.Start(); err == nil {
		log.Println("配置数据加载成功")
	} else {
		log.Fatal("配置数据加载失败：", err)
	}
	if err := core.Start(); err == nil {
		log.Println("核心服务启动成功")
	} else {
		log.Fatal("核心服务启动失败：", err)
	}
	if err := web.Start(*address, *username, *password); err == nil {
		log.Println("接口服务启动成功")
	} else {
		log.Fatal("接口服务启动失败：", err)
	}

	// 阻塞主进程，等待信号
	<-ctx.Done()
	//
	log.Println("正在关闭服务...")
	// 关闭服务
	if err := web.Close(); err != nil {
		log.Println("接口服务关闭失败：", err)
	}
	if err := core.Close(); err != nil {
		log.Println("核心服务关闭失败：", err)
	}
	if err := repo.Close(); err != nil {
		log.Println("配置数据写回失败：", err)
	}
}
