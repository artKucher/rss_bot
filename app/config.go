package app

import (
	"os"
	"strings"
)

type configStruct struct {
	TelegramToken   string
	DatabasePath    string
	RssParsePeriod  int
	AllowedChatsIds []string
}

func loadConfig() *configStruct {
	config := &configStruct{
		DatabasePath:    "store.db",
		TelegramToken:   os.Getenv("TELEGRAM_TOKEN"),
		RssParsePeriod:  3,
		AllowedChatsIds: strings.Split(os.Getenv("ALLOWED_CHATS_IDS"), ","),
	}

	return config
}

var Config = loadConfig()
