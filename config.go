package main

import (
	"os"
)

type Config struct {
	TelegramToken  string
	DatabasePath   string
	RssParsePeriod int
}

func loadConfig() *Config {
	config := &Config{
		DatabasePath:   "store.db",
		TelegramToken:  os.Getenv("TELEGRAM_TOKEN"),
		RssParsePeriod: 3,
	}

	return config
}

var config = loadConfig()
