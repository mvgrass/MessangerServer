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
	// load global_config and service config .env
	// might be useful to add some checks of critical
	globalEnvErr := godotenv.Load(".env")
	localEnvErr := godotenv.Load("config/.env")

	if globalEnvErr != nil && localEnvErr != nil {
		panic("Can't open nither global nor local .env files")
	}

	var cfg config.Config
	err := cleanenv.ReadConfig("config/local.yml", &cfg)
	if err != nil {
		panic(err)
	}

	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic(err)
	}

	// init logger - later

	// init storage
	AuthRepository := repository.InitStorage(&cfg)
	AuthService := service.CreateAuthService(AuthRepository, &cfg)

	// init router
	router := router.InitRouter(AuthService)

	router.Run(cfg.App.Addr + ":" + cfg.App.Port)

	//TODO: graceful shutdown
}
