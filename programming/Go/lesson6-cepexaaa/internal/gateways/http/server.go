package http

import (
	"context"
	"fmt"
	"homework/internal/usecase"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	host       string
	port       uint16
	router     *gin.Engine
	httpServer *http.Server
	wsHandler  *WebSocketHandler
}

type UseCases struct {
	Event  *usecase.Event
	Sensor *usecase.Sensor
	User   *usecase.User
}

func NewServer(useCases UseCases, options ...func(*Server)) *Server {
	r := gin.Default()
	wsHandler := NewWebSocketHandler(useCases)
	setupRouter(r, useCases, wsHandler)

	host := os.Getenv("HTTP_HOST")
	if host == "" {
		host = "localhost"
	}

	portStr := os.Getenv("HTTP_PORT")
	port := uint16(8080)
	if portStr != "" {
		if p, err := strconv.ParseUint(portStr, 10, 16); err == nil {
			port = uint16(p)
		}
	}

	s := &Server{router: r, host: host, port: port, wsHandler: wsHandler}
	for _, o := range options {
		o(s)
	}

	return s
}

func WithHost(host string) func(*Server) {
	return func(s *Server) {
		s.host = host
	}
}

func WithPort(port uint16) func(*Server) {
	return func(s *Server) {
		s.port = port
	}
}

func (s *Server) Run(ctx context.Context) error {
	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.host, s.port),
		Handler: s.router,
	}

	serverErr := make(chan error, 1)
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.wsHandler.Shutdown(); err != nil {
			log.Printf("websocket shutdown error: %v", err)
		}

		if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("http server shutdown error: %w", err)
		}

		return nil
	}
}
