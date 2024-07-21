package main

import (
	"github.com/telia-company/convhub-lmm-communication-service/internal/network/http"
	"go.uber.org/zap"
)

type userService struct {
	httpServer *http.Server
	logger     *zap.Logger
}

func newUserService(
	httpServer *http.Server,
	logger *zap.Logger,
) userService {
	return userService{
		httpServer: httpServer,
		logger:     logger,
	}
}

// @title Users Server API
// @version 1.0.0
// @description this API for endpoints related to user.

// @url http://localhost:8088
func main() {
	//cancelableCtx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	//wiredServer, err := inject(cancelableCtx)
	//if err != nil {
	//	panic("unable to construct DI")
	//}
	//defer wiredServer.cleanup()
	//
	//go func() {
	//	defer cancel()
	//	wiredServer.httpServer.Serve()
	//}()
	//fmt.Println("Hello, World!")
}
