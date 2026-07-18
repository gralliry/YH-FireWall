package cmdserver

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/google/shlex"
)

type Server struct {
	handler  Handler
	listener net.Listener
}

func New(config Config, handler Handler) (*Server, error) {
	server := &Server{
		handler: handler,
	}
	_ = os.Remove(config.SocketPath)
	listener, err := net.Listen("unix", config.SocketPath)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on socket: %w", err)
	}
	server.handler = handler
	server.listener = listener
	go server.acceptConn()
	return server, nil
}

func (s *Server) Running() bool {
	return s.listener != nil
}

func (s *Server) Close() error {
	if s.listener == nil {
		return nil
	}
	if err := s.listener.Close(); err != nil {
		return fmt.Errorf("failed to close cmd server: %w", err)
	}
	return nil
}

func (s *Server) acceptConn() {
	retry := time.Second
	for {
		conn, err := s.listener.Accept()
		if err == nil {
			go s.handleConn(conn)
			retry = time.Second
		} else if errors.Is(err, net.ErrClosed) {
			break
		} else {
			time.Sleep(retry)
			retry += time.Second
		}
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(30 * time.Second))

	server := bufio.NewReadWriter(
		bufio.NewReader(conn),
		bufio.NewWriter(conn),
	)
	defer server.Flush()
	// 写入终止符
	defer server.Write([]byte{0})

	command, err := server.ReadString('\n')
	if err != nil {
		server.WriteString(err.Error())
		return
	}

	args, err := shlex.Split(strings.TrimSpace(command))
	if err != nil {
		server.WriteString(err.Error())
		return
	}

	cmd := newCmd(s.handler)
	cmd.SetArgs(args)

	cmd.SetOut(server)
	cmd.SetErr(server)

	cmd.ExecuteC()
}
