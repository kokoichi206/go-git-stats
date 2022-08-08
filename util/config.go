package util

import (
	"os"
)

// Configurations
type Config struct {
	Token string
}

func LoadConfig() Config {

	token := os.Getenv("GGS_TOKEN")

	return Config{
		Token: token,
	}
}
