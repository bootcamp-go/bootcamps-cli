package config

import (
	"errors"

	"github.com/spf13/viper"
)

// Errors
var (
	ErrNoToken   = errors.New("no se encontró el token")
	ErrNoUser    = errors.New("no se encontró el usuario")
	ErrNoCompany = errors.New("no se encontró la empresa")
)

var (
	GoRepoNameFormat = "%s_bootcamp_go_w%s-%s"
)

type Configuration struct {
	Username string
	Token    string
	Company  string
	ApiUrl   string
	ApiKey   string
}

func GetConfiguration() (*Configuration, error) {
	token := viper.GetString("token")
	if token == "" {
		return nil, ErrNoToken
	}

	username := viper.GetString("username")
	if username == "" {
		return nil, ErrNoUser
	}

	company := viper.GetString("company")
	if company == "" {
		return nil, ErrNoCompany
	}

	apiUrl := viper.GetString("api_url")

	apiKey := viper.GetString("api_key")

	return &Configuration{
		Username: username,
		Token:    token,
		Company:  company,
		ApiUrl:   apiUrl,
		ApiKey:   apiKey,
	}, nil
}
