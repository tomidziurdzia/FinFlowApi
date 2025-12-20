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
	httpServer     *http.Server
	shutdownTimeout time.Duration
}

type Config struct {
	Addr            string
	ReadTimeout     time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

func NewServer(cfg Config) *Server {
	mux := http.NewServeMux()
	SetupRoutes(mux)

	readTimeout := cfg.ReadTimeout
	if readTimeout == 0 {
		readTimeout = 10 * time.Second
	}
	readHeaderTimeout := cfg.ReadHeaderTimeout
	if readHeaderTimeout == 0 {
		readHeaderTimeout = 5 * time.Second
	}
	writeTimeout := cfg.WriteTimeout
	if writeTimeout == 0 {
		writeTimeout = 30 * time.Second
	}
	idleTimeout := cfg.IdleTimeout
	if idleTimeout == 0 {
		idleTimeout = 1 * time.Minute
	}

	srv := &http.Server{
		Addr:              normalizeAddr(cfg.Addr),
		Handler:           mux,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	shutdownTimeout := cfg.ShutdownTimeout
	if shutdownTimeout == 0 {
		shutdownTimeout = 10 * time.Second
	}

	return &Server{
		httpServer:     srv,
		shutdownTimeout: shutdownTimeout,
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
		ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
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