package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	APP_PORT   string
	JWT_SECRET string

	PG_HOST string
	PG_PORT string
	PG_USER string
	PG_PASS string

	MYSQL_HOST string
	MYSQL_PORT string
	MYSQL_USER string
	MYSQL_PASS string
}

// load env and return the config struct
func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Failed to load .env")
	}
	err = godotenv.Load("tenant.env")
	if err != nil {
		log.Println("Failed to load tenant.env")
	}

	return &Config{
		APP_PORT:   os.Getenv("APP_PORT"),
		JWT_SECRET: os.Getenv("JWT_SECRET"),
		PG_HOST:    os.Getenv("PG_HOST"),
		PG_USER:    os.Getenv("PG_USER"),
		PG_PORT:    os.Getenv("PG_PORT"),
		PG_PASS:    os.Getenv("PG_PASS"),
		MYSQL_HOST: os.Getenv("MYSQL_HOST"),
		MYSQL_PORT: os.Getenv("MYSQL_PORT"),
		MYSQL_USER: os.Getenv("MYSQL_USER"),
		MYSQL_PASS: os.Getenv("MYSQL_PASS"),
	}
}
