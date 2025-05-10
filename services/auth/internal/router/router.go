package router

import (
	"MessangerServerAuth/internal/service"

	"github.com/gin-gonic/gin"
)

func InitRouter(service service.IAuthService) *gin.Engine {
	router := gin.Default()
	auth := router.Group("/api/v1/auth")
	auth.GET("/health", service.HealthHandler)
	auth.POST("/register", service.RegisterHandler)
	auth.POST("/login", service.LoginHandler)
	auth.POST("/refresh", service.RefreshHandler)
	auth.POST("/logout", service.LogoutHandler)
	auth.GET("/me", service.GetMyselfHandler)
	auth.GET("/validate", service.Validate)
	// auth.POST("/password-reset")
	// auth.POST("/password-reset/confirm")

	return router
}
