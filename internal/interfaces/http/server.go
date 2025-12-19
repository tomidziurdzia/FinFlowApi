package http

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Server struct {
	httpServer *http.Server
}

type Config struct {
	Addr string
}

func NewServer(cfg Config) *Server {
	mux := http.NewServeMux()
	SetupRoutes(mux)

	srv := &http.Server{
		Addr:              normalizeAddr(cfg.Addr),
		Handler:           mux,
		ReadTimeout:       time.Second * 10,
		ReadHeaderTimeout: time.Second * 5,
		WriteTimeout:      time.Second * 30,
		IdleTimeout:       time.Minute,
	}

	return &Server{
		httpServer: srv,
	}
}

func (s *Server) Run() error {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	errChan := make(chan error, 1)
	go func() {
		log.Printf("server has started at %s", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		return err
	case <-stop:
		log.Println("shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := s.httpServer.Shutdown(ctx); err != nil {
			return err
		}
		log.Println("server stopped")
		return nil
	}
}

func normalizeAddr(addr string) string {
	if addr == "" {
		addr = "8080"
	}
	addr = strings.TrimLeft(addr, ":")
	return ":" + addr
}