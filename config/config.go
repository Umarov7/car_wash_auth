package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type Config struct {
	AUTH_ROUTER_PORT  string
	AUTH_SERVICE_PORT string
	ACCESS_TOKEN_KEY  string
	REFRESH_TOKEN_KEY string
	DB_HOST           string
	DB_PORT           string
	DB_USER           string
	DB_NAME           string
	DB_PASSWORD       string
	REDIS_ADDRESS     string
	REDIS_PASSWORD    string
	REDIS_DB          int
}

func Load() *Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("error while loading .env file: %v", err)
	}

	return &Config{
		DB_HOST:     cast.ToString(coalesce("DB_HOST", "localhost")),
		DB_PORT:     cast.ToString(coalesce("DB_PORT", "5432")),
		DB_USER:     cast.ToString(coalesce("DB_USER", "postgres")),
		DB_NAME:     cast.ToString(coalesce("DB_NAME", "postgres")),
		DB_PASSWORD: cast.ToString(coalesce("DB_PASSWORD", "password")),

		REDIS_ADDRESS:  cast.ToString(coalesce("REDIS_ADDRESS", "redis:6379")),
		REDIS_PASSWORD: cast.ToString(coalesce("REDIS_PASSWORD", "")),
		REDIS_DB:       cast.ToInt(coalesce("REDIS_DB", "0")),

		AUTH_ROUTER_PORT:  cast.ToString(coalesce("AUTH_ROUTER_PORT", ":8081")),
		AUTH_SERVICE_PORT: cast.ToString(coalesce("AUTH_SERVICE_PORT", ":50051")),
		ACCESS_TOKEN_KEY:  cast.ToString(coalesce("ACCESS_TOKEN_KEY", "key")),
		REFRESH_TOKEN_KEY: cast.ToString(coalesce("REFRESH_TOKEN_KEY", "key")),
	}
}

func coalesce(key string, value interface{}) interface{} {
	val, exists := os.LookupEnv(key)
	if exists {
		return val
	}
	return value
}
