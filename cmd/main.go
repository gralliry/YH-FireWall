package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: yfw core|cli")
	}
	switch os.Args[1] {
	case "core":
		// 运行核心服务
		startCore()
	case "cli":
		// 运行客户端
		startClient()
	default:
		log.Println("Usage: yfw core|cli")
	}
}
