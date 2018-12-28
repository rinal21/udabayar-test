package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	Key              string
	APIKey           string
	BaseUrl          string
	WebUrl           string
	AwsKey           string
	AwsSecret        string
	S3Bucket         string
	RajaOngkirAPIUrl string
	RajaOngkirAPIKey string
	BankAPIUrl       string
}

func LoadEnv() {
	if os.Getenv("ENVIRONMENT") == "testing" {
		if err := godotenv.Load("../.env"); err != nil {
			panic("Error loading .env file")
		}
	} else {
		if err := godotenv.Load(); err != nil {
			panic("Error loading .env file")
		}
	}
}

func NewConfig() *Config {
	LoadEnv()

	return &Config{
		Port: os.Getenv("APP_PORT"),
		Key:  os.Getenv("APP_KEY"),

		APIKey: os.Getenv("API_KEY"),

		BaseUrl: "http://localhost:5555/",
		WebUrl:  "http://dev.udabayar.com/",

		AwsKey:    os.Getenv("AWS_KEY"),
		AwsSecret: os.Getenv("AWS_SECRET"),

		S3Bucket: os.Getenv("S3_BUCKET"),

		RajaOngkirAPIUrl: os.Getenv("RAJAONGKIR_API_URL"),
		RajaOngkirAPIKey: os.Getenv("RAJAONGKIR_API_KEY"),

		BankAPIUrl: os.Getenv("BANK_API_URL"),
	}
}
