package config

import (
	"flag"
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
		port := flag.String("b", "http://localhost:8000", "Port for the server")

		flag.Parse()

		ConfigInstance = &Config{
			Address:       *address,
			ResultAddress: *port,
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
