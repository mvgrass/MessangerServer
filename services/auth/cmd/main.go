package main

import (
	"MessangerServer/services/auth/internal/config"
	"MessangerServer/services/auth/internal/repository"
	"MessangerServer/services/auth/internal/router"
	"MessangerServer/services/auth/internal/service"

	"github.com/ilyakaznacheev/cleanenv"
)

func main() {
	// init config
	var cfg config.Config
	cleanenv.ReadConfig("config/local.yml", &cfg)

	// init logger - later

	// init storage
	userRepository := repository.InitStorage(&cfg)
	userService := service.CreateUserService(userRepository)

	// init router
	router := router.InitRouter(userService)

	router.Run(cfg.App.Addr + ":" + cfg.App.Port)
}
