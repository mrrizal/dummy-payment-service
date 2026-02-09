package config

import "os"

type databaseConfig struct {
	DSN string
}

type appConfig struct {
	Port string
}

type Config struct {
	Database databaseConfig
	App      appConfig
}

func LoadConfig() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "file:payments.db"
	}

	return Config{
		Database: databaseConfig{DSN: dsn},
		App:      appConfig{Port: port},
	}
}
