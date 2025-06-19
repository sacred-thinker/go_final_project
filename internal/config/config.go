package config

import "os"

type Config struct {
	TodoPassword string
}

func LoadConfig() *Config {
	return &Config{
		TodoPassword: os.Getenv("TODO_PASSWORD"),
	}
}
