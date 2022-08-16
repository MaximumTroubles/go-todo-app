package todo

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

// Run server
func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		TLSConfig:         nil,
		ReadTimeout:       10 * time.Second, // 10 seconds
		ReadHeaderTimeout: 0,
		WriteTimeout:      10 * time.Second, // 10 seconds
		IdleTimeout:       0,
		MaxHeaderBytes:    1 << 20, // "1 times 2, 20 times" or "1 to the extent 20" = 1 048 576 bytes = 8 000 000 bit = 1 MB
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
		BaseContext:       nil,
		ConnContext:       nil,
	}
	// under the hood launch endless for loop and listen all incoming request and handle them.
	return s.httpServer.ListenAndServe()
}

// Shutdown server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
