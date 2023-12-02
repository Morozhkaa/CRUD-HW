package http

import (
	"billing-service/pkg/infra/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func initRouter(a *Adapter, r *gin.Engine) error {
	log := logger.Get()
	log.Info("initializing handlers and routes...")

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	r.Use(cors.New(config))
	r.Use(sloggin.New(log))

	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/health", a.health)

	v1_auth := r.Group("/").Use(authMiddleware(a))
	{
		v1_auth.GET("/balance", a.getBalance)
		v1_auth.POST("/balance", a.updateBalance)
	}
	return nil
}
