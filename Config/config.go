package config

import (
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"strings"
)

type Config struct {
	SQL   *SQL
	Kafka *Kafka
}

type SQL struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

type Kafka struct {
	Port    string
	Brokers []string
	Topic   string
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file")
		return nil
	}
	slog.Info("Loading fine .env file")
	return &Config{
		SQL: &SQL{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Database: os.Getenv("DB_NAME"),
			Username: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
		},
		Kafka: &Kafka{
			Port:    os.Getenv("KAFKA_PORT"),
			Brokers: strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
			Topic:   os.Getenv("KAFKA_TOPIC"),
		},
	}
}
