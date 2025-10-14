package cmdserver

import (
	"errors"
	"fmt"
	"github.com/google/shlex"
	"log"
	"net"
	"os"
	"strings"
)

var (
	SocketPath string
)

var (
	listener  net.Listener
	isRunning bool
)

func Start(h Handler) (err error) {
	// 删除残留的 cmdserver 文件
	_ = os.Remove(SocketPath)
	// 监听 Unix 域套接字
	listener, err = net.Listen("unix", SocketPath)
	if err != nil {
		return fmt.Errorf("failed to listen on socket: %w", err)
	}
	// 设置处理函数
	handler = h
	// 启动监听
	go acceptConn()
	//
	isRunning = true
	return nil
}

func Close() error {
	isRunning = false
	return listener.Close()
}

func IsRunning() bool {
	return isRunning
}

func acceptConn() {
	for {
		conn, err := listener.Accept()
		if err == nil {
			go handleConn(conn)
		} else if errors.Is(err, net.ErrClosed) {
			break
		}
	}
}

func handleConn(conn net.Conn) {
	defer func() { _ = conn.Close() }()
	// 创建缓冲区
	cmdBytes := make([]byte, 1024)
	// 读取命令
	n, err := conn.Read(cmdBytes)
	if err != nil {
		_, _ = conn.Write([]byte(err.Error()))
		return
	}
	// 处理命令(这里要截取，不如会取到后面的未写入字符)
	cmdStr := strings.TrimSpace(string(cmdBytes[:n]))
	// 解析命令
	args, err := shlex.Split(cmdStr)
	if err != nil {
		_, _ = conn.Write([]byte(err.Error()))
		return
	}
	//
	log.Printf("Args(%d): %v", len(args), args)
	// 解析并执行命令
	result := handleArgs(args) + "\n"
	// 返回结果
	_, _ = conn.Write([]byte(result))
}
