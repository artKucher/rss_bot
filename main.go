package main

import (
	"database/sql"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db := connectToDB()
	bot := buildBot()
	go handleIncomingMessages(bot, db)

	for {
		fmt.Println("Start parsing RSS")
		processRSS(db, bot)
	}
}

func processRSS(db *sql.DB, bot *tgbotapi.BotAPI) {
	rss_infos, err := getRSSUrl(db)
	if err != nil {
		panic(err)
	}

	for _, rss_info := range rss_infos {
		current_time := time.Now()
		posts, err := getNewPosts(rss_info.url)
		if err != nil {
			panic(err)
		}
		posts = filterPosts(posts, rss_info.lastPostDateTime, rss_info.keywords)
		sendPostsToTelegram(bot, posts, rss_info.telegramChatId)
		setLatestPostDateTime(db, rss_info.url, rss_info.telegramChatId, current_time)
	}

	time.Sleep(time.Duration(config.RssParsePeriod) * time.Hour)
}

func sendPostsToTelegram(bot *tgbotapi.BotAPI, posts []Post, chatId int64) {
	for _, post := range posts {
		content := fmt.Sprintf("%s\n%s\n%s", post.Title, post.Description, post.Link)
		msg := tgbotapi.NewMessage(chatId, content)
		bot.Send(msg)
	}
}
