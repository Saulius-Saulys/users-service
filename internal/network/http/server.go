package http

import (
	"fmt"
	"github.com/telia-company/convhub-lmm-communication-service/internal/config"
	"go.uber.org/zap"
	"net/http"
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
		ReadHeaderTimeout: 3 * time.Second,
	}
	return &Server{
		router:     router,
		httpServer: httpServer,
		logger:     logger,
	}
}
