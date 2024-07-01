package config

import (
	"flag"
	"os"
	"sync"
)

type Config struct {
	Address         string
	ResultAddress   string
	FlagLogLevel    string
	FileStoragePath string
}

var (
	Instance *Config
	once     sync.Once
)

func Init() *Config {
	once.Do(func() {
		var address, baseURL, flagLogLevel, fileStoragePath string

		flag.StringVar(&address, "a", "localhost:8080", "Address for the server")
		flag.StringVar(&baseURL, "b", "http://localhost:8080", "Port for the server")
		flag.StringVar(&flagLogLevel, "l", "info", "log level")
		flag.StringVar(&fileStoragePath, "f", "/tmp/short-url-db.json", "Path for file storage")

		flag.Parse()

		if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
			address = envRunAddr
		}

		if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
			baseURL = envBaseURL
		}

		if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
			flagLogLevel = envLogLevel
		}

		if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
			fileStoragePath = envFileStoragePath
		}

		Instance = &Config{
			Address:         address,
			ResultAddress:   baseURL,
			FlagLogLevel:    flagLogLevel,
			FileStoragePath: fileStoragePath,
		}
	})
	return Instance
}

func (c *Config) GetAddress() string {
	return c.Address
}

func (c *Config) GetResultAddress() string {
	return c.ResultAddress
}
