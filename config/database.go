package config

// Database configuration

import "os"



type DatabaseConfig struct {
	Host     string
	Name     string
	Password string
	User     string
	Port     string
}

func NewDatabase() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     os.Getenv("HOST"),
		Name:     os.Getenv("DATABASE_NAME"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		User:     os.Getenv("DATABASE_USER"),
		Port:     os.Getenv("DATABASE_PORT"),
	}
}
