package config

type Config struct {
	App struct {
		Name string `yaml:"name" env:"APP_NAME"`
		Addr string `yaml:"addr" env:"APP_ADDR" env-default:"localhost"`
		Port string `yaml:"port" env:"APP_PORT" env-default:"8080"`
	} `yaml:app`

	Db struct {
		User     string `yaml:"user" env:"DB_USERNAME"`
		Password string `yaml:"password" env:"DB_PASSWORD"`
		DbName   string `yaml:"dbname" env:"DB_NAME" env-default:"test"`
		Addr     string `yaml:"addr" env:"DB_ADDR" env-default:"localhost"`
		Port     string `yaml:"port" env:"DB_PORT" env-default:"5432"`
	} `yaml:db`
}
