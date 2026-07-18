package webserver

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
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
		if err := app.Listen(config.Address, fiber.ListenConfig{
			DisableStartupMessage: true,
		}); err != nil {
			log.Error(err)
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
	if err := s.app.ShutdownWithTimeout(5 * time.Second); err != nil {
		return fmt.Errorf("failed to close webserver: %w", err)
	}
	// 清理逻辑在 New 中
	return nil
}
