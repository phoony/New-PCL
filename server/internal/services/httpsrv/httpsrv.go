package httpsrv

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

type HTTPService struct {
	Port       int
	httpServer *http.Server
}

func (s *HTTPService) StartHTTPServer(r *mux.Router) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		return fmt.Errorf("failed to create listener: %w", err)
	}

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.Port),
		Handler: r,
	}

	go func() {
		log.Printf("HTTP server listening on :%d", s.Port)
		if err := s.httpServer.Serve(lis); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	return nil
}

func (s *HTTPService) Shutdown() error {
	var wg sync.WaitGroup
	var err error

	// Shutdown HTTP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		if s.httpServer != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 30)
			defer cancel()
			if shutdownErr := s.httpServer.Shutdown(ctx); shutdownErr != nil {
				log.Printf("HTTP server shutdown error: %v", shutdownErr)
				err = shutdownErr
			} else {
				log.Println("HTTP server stopped")
			}
		}
	}()

	wg.Wait()
	if err != nil {
		return fmt.Errorf("failed to shutdown HTTP service: %w", err)
	}
	return nil
}
