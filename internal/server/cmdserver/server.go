package cmdserver

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/google/shlex"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
)

var (
	listener  net.Listener
	isRunning bool
)

type Config struct {
	Enable     bool   `json:"enable"`
	SocketPath string `json:"socket_path"`
}

func Start(h Handler, config Config) (err error) {
	// 删除残留的 cmdserver 文件
	_ = os.Remove(config.SocketPath)
	// 监听 Unix 域套接字
	listener, err = net.Listen("unix", config.SocketPath)
	if err != nil {
		return fmt.Errorf("failed to listen on socket: %w", err)
	}
	// 启动监听
	go acceptConn(h)
	//
	isRunning = true
	return nil
}

func Close() error {
	if !isRunning {
		return fmt.Errorf("cmdserver is not running")
	}
	return listener.Close()
}

func IsRunning() bool {
	return isRunning
}

func acceptConn(handler Handler) {
	for {
		conn, err := listener.Accept()
		if err == nil {
			go handleConn(conn, handler)
		} else if errors.Is(err, net.ErrClosed) {
			break
		} else {
			log.Printf("cmdserver accept error: %v", err)
			// 避免 Accept 持续失败导致忙等
			// 使用简单的 sleep 退避
			continue
		}
	}
}

func handleConn(conn net.Conn, handler Handler) {
	defer func() { _ = conn.Close() }()

	// 解析并执行命令
	cmder := newCommand(handler)
	var outBuf, errBuf bytes.Buffer
	cmder.SetOut(&outBuf)
	cmder.SetErr(&errBuf)
	//
	client := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	for {
		outBuf.Reset()
		errBuf.Reset()
		// 读取命令
		command, err := client.ReadString('\n')
		if err != nil {
			return
		}
		// 解析命令
		args, err := shlex.Split(command)
		if err != nil {
			errBuf.WriteString(err.Error())
			_, _ = client.Write(errBuf.Bytes())
			_ = client.WriteByte(0)
			_ = client.Flush()
			return
		}
		// todo 参数读取展示
		log.Printf("Args(%d): %v", len(args), args)
		// 设置命令参数
		cmder.SetArgs(args)
		if _, err = cmder.ExecuteC(); err != nil {
			if _, err = client.Write(errBuf.Bytes()); err != nil {
				return
			}
		} else {
			if _, err = client.Write(outBuf.Bytes()); err != nil {
				return
			}
		}
		if err = client.WriteByte(0); err != nil {
			return
		}
		if err = client.Flush(); err != nil {
			return
		}
	}
}
