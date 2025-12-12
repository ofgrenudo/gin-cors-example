package main

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/ofgrenudo/gin-example/api/handlers/health"
	"github.com/ofgrenudo/gin-example/internal/config/env"
	"github.com/ofgrenudo/gin-example/internal/config/logging"
)

func main() {
	env.InitializeConfig()
	cleanup_logger := logging.InitializeGlobalLogger(env.GlobalConfig.LogPath, env.GlobalConfig.LogLevel)
	defer cleanup_logger()

	router := gin.Default()
	v1 := router.Group("/api/v1/")
	healthGroup := v1.Group("/health")
	healthGroup.GET("/ping", health.Ping)

	address := fmt.Sprintf(":%d", env.GlobalConfig.BackendPort)
	if err := router.Run(address); err != nil {
		slog.Error("Failed to start. Shutting down...", slog.Any("Error", err))
		log.Fatalf("Failed to start. Shutting down...")
	}
}
