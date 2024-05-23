package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host      string
	Port      string
	User      string
	Password  string
	Name      string
	DBAddress string
	JWTExpire int64
	JWTSecret string
}

var Envs = initConfig()

func initConfig() DBConfig {
	godotenv.Load()

	return DBConfig{
		Host:      getEnv("DB_HOST", "http://localhost"),
		Port:      getEnv("DB_PORT", "8080"),
		User:      getEnv("DB_USER", "root"),
		Password:  getEnv("DB_PASSWORD", "root"),
		Name:      getEnv("DB_NAME", "go_ecommerce"),
		DBAddress: fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		JWTSecret: getEnv("JWT_SECRET", "v3rY-s3creT-Jw7"),
		JWTExpire: getEnvAsInt("JWT_EXP", 3600*24*7),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return i
		}
	}
	return fallback
}
