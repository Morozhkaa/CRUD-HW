package http

import (
	"user-service/pkg/infra/logger"

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
	r.Use(cors.New(config))
	r.Use(sloggin.New(log))

	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/health", a.health)

	r.POST("/user", a.createUser) // по сути register
	r.POST("/login", a.login)
	r.POST("/verify", a.verify)

	r.GET("/user/:username", a.getUser)
	r.PUT("/user/:username", a.updateUser)
	r.DELETE("/user/:username", a.deleteUser)
	return nil
}
