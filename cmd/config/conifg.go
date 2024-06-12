package config

import (
	"flag"
	"sync"
)

type Config struct {
	Address string
	Port    string
}

var (
	ConfigInstance *Config
	once           sync.Once
)

func Init() *Config {
	once.Do(func() {
		address := flag.String("b", "localhost", "Address for the server")
		port := flag.String("a", "8080", "Port for the server")

		flag.Parse()

		ConfigInstance = &Config{
			Address: *address,
			Port:    *port,
		}
	})
	return ConfigInstance
}

func (c *Config) GetFullAddress() string {
	return c.Address + ":" + c.Port
}
