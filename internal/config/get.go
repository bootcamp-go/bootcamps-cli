package config

import (
	"errors"

	"github.com/spf13/viper"
)

// Errors
var (
	ErrNoToken   = errors.New("No se encontró el token")
	ErrNoUser    = errors.New("No se encontró el usuario")
	ErrNoCompany = errors.New("No se encontró la empresa")
)

type Configuration struct {
	Username string
	Token    string
	Company  string
}

func GetConfiguration() (*Configuration, error) {
	// Get from config file using Viper
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

	return &Configuration{
		Username: username,
		Token:    token,
		Company:  company,
	}, nil
}
