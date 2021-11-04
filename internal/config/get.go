package config

import (
	"errors"

	"github.com/spf13/viper"
)

// Errors
var (
	ErrNoTokenDH = errors.New("no se encontró el token de DH")
	ErrNoToken   = errors.New("no se encontró el token")
	ErrNoUser    = errors.New("no se encontró el usuario")
	ErrNoCompany = errors.New("no se encontró la empresa")
)

type Configuration struct {
	Username string
	Token    string
	Company  string
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

	return &Configuration{
		Username: username,
		Token:    token,
		Company:  company,
	}, nil
}

func GetTokenDH() (string, error) {
	token := viper.GetString("tokendh")
	if token == "" {
		return "", ErrNoTokenDH
	}

	return token, nil
}
