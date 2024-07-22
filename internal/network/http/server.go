package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/Saulius-Saulys/users-service/internal/config"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Server struct {
	router     *Router
	httpServer *http.Server
	logger     *zap.Logger
}

func NewHTTPServer(
	router *Router,
	config *config.Config,
	logger *zap.Logger,
) *Server {
	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%s", config.HTTPPort),
		Handler:           router.ginEngine,
		ReadHeaderTimeout: 3 * time.Second,
	}
	return &Server{
		router:     router,
		httpServer: httpServer,
		logger:     logger,
	}
}

func (s *Server) Serve() {
	shutdownWaitGrp := &sync.WaitGroup{}

	s.runHTTPServer(shutdownWaitGrp)

	s.logger.Info("Started server", zap.String("address", s.httpServer.Addr))

	shutdownWaitGrp.Wait()
}

func (s *Server) GracefulStop() {
	timeoutContext, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	<-done

	err := s.httpServer.Shutdown(timeoutContext)
	if err != nil {
		s.logger.Error("error when gracefully stopping http server", zap.Error(err))
	}
}

func (s *Server) runHTTPServer(shutdownWaitGrp *sync.WaitGroup) {
	shutdownWaitGrp.Add(1)
	go func() {
		defer shutdownWaitGrp.Done()
		if err := s.httpServer.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				s.logger.Info("Graceful shutdown of http server started")
				return
			}
			s.logger.Error("HTTP server encountered an error", zap.Error(err))
			return
		}
	}()
}
