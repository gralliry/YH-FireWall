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
	server := &Server{}
	if !config.Enable {
		return server, nil
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

	reader := bufio.NewReader(conn)
	command, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	args, err := shlex.Split(strings.TrimSpace(command))
	if err != nil {
		conn.Write([]byte(err.Error()))
		return
	}

	cmd := newCmd(s.handler)
	cmd.SetOut(conn)
	cmd.SetErr(conn)

	cmd.SetArgs(args)
	cmd.ExecuteC()
}
