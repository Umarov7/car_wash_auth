package api

import (
	_ "auth-service/api/docs"
	"auth-service/api/handler"
	"auth-service/config"
	"auth-service/storage"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Authorazation
// @version 1.0
// @description Authorazation API of On-Demand Car Wash Service
// @host localhost:8081
// @BasePath /auth
func NewRouter(s storage.IStorage, cfg *config.Config) *gin.Engine {
	h := handler.NewHandler(s, cfg)

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := r.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/logout", h.Logout)
		auth.POST("/refresh", h.RefreshToken)
		auth.POST("/validate", h.ValidateToken)
	}

	return r
}
