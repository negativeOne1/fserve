package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	// DefaultWriteTimeout for the http server
	DefaultWriteTimeout = 2 * time.Minute

	// DefaultReadTimeout for the http server
	DefaultReadTimeout = 2 * time.Minute

	// DefaultShutdownTimeout for the http server
	DefaultShutdownTimeout = 5 * time.Second

	// DefaultResponseTimeout for the http server
	DefaultResponseTimeout = 10 * time.Second
)

type HTTPServer struct {
	host    string
	port    int
	handler http.Handler
}

func NewHTTPServer(host string, port int, handler http.Handler) *HTTPServer {
	return &HTTPServer{
		host:    host,
		port:    port,
		handler: handler,
	}
}

func (s HTTPServer) Address() string {
	return fmt.Sprintf("%s:%d", s.host, s.port)
}

func (s HTTPServer) Run(ctx context.Context) error {
	srv := http.Server{
		Addr:         s.Address(),
		Handler:      s.handler,
		WriteTimeout: DefaultWriteTimeout,
		ReadTimeout:  DefaultReadTimeout,
	}

	go func() {
		log.Info().Str("Address", s.Address()).Msg("http server started")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Err(err).Msg("error in http server")
		}
	}()

	<-ctx.Done()

	ctxShutDown, cancel := context.WithTimeout(context.Background(), DefaultShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctxShutDown); err != nil {
		return fmt.Errorf("failed to shutdown http server properly: %s", err)
	}

	return nil
}
