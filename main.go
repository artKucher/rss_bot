package main

import (
	"database/sql"
	"fmt"
	"rss_bot/app"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db := app.ConnectToDB()
	bot := app.BuildBot()
	go app.HandleIncomingMessages(bot, db)

	for {
		fmt.Println("Start parsing RSS")
		processRSS(db, bot)
	}
}

func processRSS(db *sql.DB, bot *tgbotapi.BotAPI) {
	rss_infos, err := app.GetRSSUrl(db)
	if err != nil {
		panic(err)
	}

	for _, rss_info := range rss_infos {
		current_time := time.Now()
		posts, err := app.GetNewPosts(rss_info.Url)
		if err != nil {
			panic(err)
		}
		posts = app.FilterPosts(posts, rss_info.LastPostDateTime, rss_info.Keywords)
		sendPostsToTelegram(bot, posts, rss_info.TelegramChatId)
		app.SetLatestPostDateTime(db, rss_info.Url, rss_info.TelegramChatId, current_time)
		time.Sleep(10 * time.Minute)
	}

	time.Sleep(time.Duration(app.Config.RssParsePeriod) * time.Hour)
}

func sendPostsToTelegram(bot *tgbotapi.BotAPI, posts []app.Post, chatId int64) {
	for _, post := range posts {
		content := fmt.Sprintf("%s\n%s\n%s", post.Title, post.Description, post.Link)
		msg := tgbotapi.NewMessage(chatId, content)
		bot.Send(msg)
	}
}
