package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	socketPath := flag.String("s", "/tmp/yfw.sock", "Path to the socket file")
	flag.Parse()

	conn, err := net.Dial("unix", *socketPath)
	if err != nil {
		log.Println("Core service is not running. Please start it with: yfwd")
		os.Exit(1)
	}
	defer conn.Close()

	client := bufio.NewReadWriter(
		bufio.NewReader(conn),
		bufio.NewWriter(conn),
	)

	// 非交互模式：直接发送参数并退出
	args := flag.Args()

	buf, err := json.Marshal(args)
	if err != nil {
		log.Fatal("Failed to encode command:", err)
	}
	if _, err = client.WriteString(string(buf) + "\n"); err != nil {
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
