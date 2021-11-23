package config

type Config struct {
	API             string `env:"API_URL,required"`
	ServerPort      int    `env:"SERVER_PORT" envDefault:"3001"`
	LinkServiceAddr string `env:"LINK_SERVICE_ADDR"`
}
