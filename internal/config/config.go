package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	HTTPConfig HTTPConfig
}

type HTTPConfig struct {
	Host             string
	Port             string
	TimeoutInSeconds int
}

func Load() (*Config, error) {
	godotenv.Load()

	const minHTTPTimeout = 1
	httpTimeoutInSeconds, err := strconv.Atoi(os.Getenv("HTTP_TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("config http timeout value misconfigured: %w", err)
	}
	if httpTimeoutInSeconds < minHTTPTimeout {
		return nil, fmt.Errorf("http timeout must be >= %d, got %d", minHTTPTimeout, httpTimeoutInSeconds)
	}

	httpConfig := HTTPConfig{
		Host:             os.Getenv("HTTP_HOST"),
		Port:             os.Getenv("HTTP_PORT"),
		TimeoutInSeconds: httpTimeoutInSeconds,
	}

	if httpConfig.Host == "" || httpConfig.Port == "" {
		return nil, fmt.Errorf("missing environment variables")
	}

	return &Config{
		HTTPConfig: httpConfig,
	}, nil
}
