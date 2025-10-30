package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func startClient() {
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
}
