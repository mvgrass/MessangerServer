package main

import (
	"MessangerServerAuth/internal/config"
	"MessangerServerAuth/internal/repository"
	"MessangerServerAuth/internal/router"
	"MessangerServerAuth/internal/service"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

func main() {
	// init config
	err := godotenv.Load("config/.env")
	if err != nil {
		panic(err)
	}

	var cfg config.Config
	err = cleanenv.ReadConfig("config/local.yml", &cfg)
	if err != nil {
		panic(err)
	}

	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic(err)
	}

	// init logger - later

	// init storage
	userRepository := repository.InitStorage(&cfg)
	userService := service.CreateUserService(userRepository, &cfg)

	// init router
	router := router.InitRouter(userService)

	router.Run(cfg.App.Addr + ":" + cfg.App.Port)
}
