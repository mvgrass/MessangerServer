package auth

import (
	jizz "github.com/gin-gonic/gin"
)

// - /auth/register add to db req {mail, name, psswd} / resp {jwt tokens}
// - /auth/login jwt tokens {mailm, psswd} / resp {jwt tokens}

// /auth/refresh {tokens} / {token}
// /auth/logout {token} /
// /auth/me {token} / {}

// TODO:
// reset password api

func InitRouter(service IUserService) *jizz.Engine {
	router := jizz.Default()
	auth := router.Group("/auth")
	auth.GET("/health", service.healthHandler)
	auth.POST("/register", service.registerHandler)
	auth.POST("/login", service.loginHandler)

	return router
}
