package webserver

import (
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

type Server struct {
	mu  sync.Mutex
	app *fiber.App
}

func New(config Config, handler Handler) (*Server, error) {
	server := &Server{}
	if !config.Enable {
		return server, nil
	}
	app, err := newApp(config, handler)
	if err != nil {
		return nil, err
	}
	// 启动监听
	server.app = app
	go func() {
		if err := app.Listen(config.Address, fiber.ListenConfig{
			DisableStartupMessage: true,
		}); err != nil {
			log.Error(err)
			server.mu.Lock()
			server.app = nil
			server.mu.Unlock()
		}
	}()
	return server, nil
}

func (s *Server) Running() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.app != nil
}

func (s *Server) Close() error {
	s.mu.Lock()
	app := s.app
	s.app = nil
	s.mu.Unlock()
	if app == nil {
		return nil
	}
	if err := app.ShutdownWithTimeout(5 * time.Second); err != nil {
		return fmt.Errorf("failed to close webserver: %w", err)
	}
	return nil
}
