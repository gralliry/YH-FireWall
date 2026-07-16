package webserver

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	app *fiber.App
}

func New(config Config, handler Handler) (*Server, error) {
	server := &Server{}
	if !config.Enable {
		return server, nil
	}
	app := newApp(config, handler)
	// 启动监听
	server.app = app
	go func() {
		if err := app.Listen(config.Address); err != nil {
			server.app = nil
		}
	}()
	return server, nil
}

func (s *Server) Running() bool {
	return s.app != nil
}

func (s *Server) Close() error {
	if s.app == nil {
		return nil
	}
	if err := s.app.Shutdown(); err != nil {
		return fmt.Errorf("failed to close webserver: %w", err)
	}
	// 清理逻辑在 New 中
	return nil
}
