package unix

import (
	"bufio"
	"errors"
	"net"
	"os"
	"strings"
)

var (
	listener   net.Listener
	socketPath string = "/tmp/firewall.sock"
)

func Start(path string) (err error) {
	if path != "" {
		socketPath = path
	}
	// 删除残留的 socket 文件
	_ = os.Remove(socketPath)
	// 监听 Unix 域套接字
	listener, err = net.Listen("unix", socketPath)
	if err != nil {
		return err
	}
	go acceptConn()
	return nil
}

func Close() error {
	if listener == nil {
		return nil
	}
	return listener.Close()
}

func acceptConn() {
	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				break
			} else {
				continue
			}
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	cmdLine, _ := reader.ReadString(byte(0))
	cmd := strings.TrimSpace(cmdLine)

	switch cmd {
	case "status":
		_, _ = conn.Write([]byte("running\n"))
	case "stop":
		_, _ = conn.Write([]byte("stopping...\n"))
		os.Remove(socketPath)
		os.Exit(0)
	default:
		conn.Write([]byte("unknown command\n"))
	}
}
