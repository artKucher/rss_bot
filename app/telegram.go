package app

import (
	"database/sql"
	"fmt"
	"slices"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func BuildBot() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(Config.TelegramToken)
	if err != nil {
		panic(err)
	}
	return bot
}

func HandleIncomingMessages(bot *tgbotapi.BotAPI, db *sql.DB) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60 * 3
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		fmt.Printf("[%s] %d %s\n", update.Message.From.UserName, update.Message.Chat.ID, update.Message.Text)
		if !slices.Contains(Config.AllowedChatsIds, fmt.Sprintf("%d", update.Message.Chat.ID)) {
			message := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, you are not allowed to use this bot")
			bot.Send(message)
			continue
		}

		parsedMessage := strings.Split(update.Message.Text, "\n")
		rssUrl := parsedMessage[0]

		if checkRSSUrlAlreadyExists(db, rssUrl, update.Message.Chat.ID) {
			deleteRSSUrl(db, rssUrl, update.Message.Chat.ID)
			message := tgbotapi.NewMessage(update.Message.Chat.ID, "RSS URL deleted")
			bot.Send(message)
			continue
		}

		_, err := GetNewPosts(rssUrl)
		if err != nil {
			message := tgbotapi.NewMessage(update.Message.Chat.ID, "Error getting posts")
			bot.Send(message)
			continue
		}

		keywords := []string{}
		if len(parsedMessage) == 2 {
			keywords = strings.Split(parsedMessage[1], "|")
		}

		addRSSUrl(db, rssUrl, update.Message.Chat.ID, keywords)
		message := tgbotapi.NewMessage(update.Message.Chat.ID, "RSS URL added")
		bot.Send(message)

	}
}
