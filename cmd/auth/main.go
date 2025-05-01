package main

import (
	"MessangerServer/internal/auth"
	authCfg "MessangerServer/internal/config/auth"

	"github.com/ilyakaznacheev/cleanenv"
)

func main() {
	// init config
	var cfg authCfg.Config
	cleanenv.ReadConfig("../../config/auth/local.yml", &cfg)

	// init logger - later

	// init storage
	userRepository := auth.InitStorage(&cfg)
	userService := auth.CreateUserService(userRepository)

	// init router
	router := auth.InitRouter(userService)

	router.Run(cfg.App.Addr + ":" + cfg.App.Port)
}
