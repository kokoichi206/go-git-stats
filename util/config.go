package util

import (
	"os"
)

// Configurations
type Config struct {
	Token      string
	ApiBaseURL string
}

// Load configurations for actual usecase.
func LoadConfig() Config {

	token := os.Getenv("GGS_TOKEN")

	return Config{
		Token:      token,
		ApiBaseURL: "https://api.github.com",
	}
}
