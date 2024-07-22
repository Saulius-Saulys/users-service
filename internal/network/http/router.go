package http

import (
	"github.com/Saulius-Saulys/users-service/internal/config"
	"github.com/Saulius-Saulys/users-service/internal/network/http/controller"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"time"
)

type Router struct {
	ginEngine      *gin.Engine
	basePath       string
	userController *controller.User
}

func NewRouter(
	ginEngine *gin.Engine,
	logger *zap.Logger,
	config *config.Config,
	userController *controller.User,
) *Router {
	router := &Router{
		ginEngine:      ginEngine,
		basePath:       config.BasePath,
		userController: userController,
	}

	router.configureCORS()

	ginEngine.Use(errorLoggerMiddleware(logger))
	ginEngine.Use(requestLoggerMiddleware(logger, time.RFC3339))
	ginEngine.Use(ginzap.RecoveryWithZap(logger, false))

	router.configureEndpoints()

	return router
}

func (r *Router) configureEndpoints() {
	baseGroup := r.ginEngine.Group(r.basePath + "/users")
	baseGroup.POST("", r.userController.Create)
	baseGroup.PUT("/:id", r.userController.Update)
	baseGroup.DELETE("/:id", r.userController.Delete)
	baseGroup.GET("/:country", r.userController.GetByCountry)
}

func (r *Router) configureCORS() {
	r.ginEngine.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))
}

func requestLoggerMiddleware(logger *zap.Logger, timeFormat string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		status := c.Writer.Status()

		end := time.Now()
		latency := end.Sub(start)
		end = end.UTC()

		fields := []zapcore.Field{
			zap.Int("status", status),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("time", end.Format(timeFormat)),
			zap.Duration("latency", latency),
		}

		logger.Info(path,
			fields...,
		)
	}
}

func errorLoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 {
			return
		}
		ginErr := c.Errors[0].Err

		// On bigger project I would add custom error handling here
		logger.Error("API error appeared", zap.Error(ginErr))
		c.JSON(c.Writer.Status(), ginErr)
	}
}
