package config

import "os"

type DBConfig struct {
	Client   string
	Host     string
	User     string
	Password string
	Name     string
}

func NewDatabaseConfig() *DBConfig {
	LoadEnv()

	return &DBConfig{
		Client:   os.Getenv("DB_CLIENT"),
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}
}
