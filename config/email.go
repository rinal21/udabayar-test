package config

import "os"

type EmailConfig struct {
	Sender string
}

func NewEmailConfig() *EmailConfig {
	LoadEnv()

	return &EmailConfig{
		Sender: os.Getenv("EMAIL_SENDER"),
	}
}
