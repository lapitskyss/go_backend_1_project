package config

type Config struct {
	ServerPort      int    `env:"SERVER_PORT" envDefault:"3001"`
	LinkServiceAddr string `env:"LINK_SERVICE_ADDR"`
}
