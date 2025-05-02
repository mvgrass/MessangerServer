package router

import (
	"MessangerServer/services/auth/internal/service"

	"github.com/gin-gonic/gin"
)

// - /auth/register add to db req {mail, name, psswd} / resp {jwt tokens}
// - /auth/login jwt tokens {mailm, psswd} / resp {jwt tokens}

// /auth/refresh {tokens} / {token}
// /auth/logout {token} /
// /auth/me {token} / {}

// TODO:
// reset password api

func InitRouter(service service.IUserService) *gin.Engine {
	router := gin.Default()
	auth := router.Group("/auth")
	auth.GET("/health", service.HealthHandler)
	auth.POST("/register", service.RegisterHandler)
	auth.POST("/login", service.LoginHandler)

	return router
}
