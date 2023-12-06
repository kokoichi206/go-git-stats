package util

import (
	"fmt"
	"os"
	"regexp"
)

// Configurations
type Config struct {
	Token      string
	ApiBaseURL string
}

// Load configurations for actual usecase.
func LoadConfig() (Config, error) {

	token := os.Getenv("GGS_TOKEN")

	if !isValidFormat(token) {
		return Config{}, fmt.Errorf("Your token: '%s' is invalid format.\nPlease check your environment variable [GGS_TOKEN].", token)
	}

	return Config{
		Token:      token,
		ApiBaseURL: "https://api.github.com",
	}, nil
}

// Returns whether the string matches the GitHub personal access token format.
// see: https://github.blog/2021-04-05-behind-githubs-new-authentication-token-formats/
func isValidFormat(token string) bool {
	// Token format
	const TokenRegex = `^ghp_[a-zA-Z0-9]{36}$`

	// empty token means that the token has been not set, so it should be true
	if token == "" {
		return true
	}

	r := regexp.MustCompile(TokenRegex)
	return r.MatchString(token)
}
