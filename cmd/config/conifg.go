package config

import (
	"flag"
	"os"
	"sync"
)

type Config struct {
	Address       string
	ResultAddress string
}

var (
	ConfigInstance *Config
	once           sync.Once
)

func Init() *Config {
	once.Do(func() {
		address := flag.String("a", "localhost:8080", "Address for the server")
		baseURL := flag.String("b", "http://localhost:8000", "Port for the server")

		flag.Parse()

		if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
			address = &envRunAddr
		}

		if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
			baseURL = &envBaseURL
		}

		ConfigInstance = &Config{
			Address:       *address,
			ResultAddress: *baseURL,
		}
	})
	return ConfigInstance
}

func (c *Config) GetAddress() string {
	return c.Address
}

func (c *Config) GetResultAddress() string {
	return c.Address
}
