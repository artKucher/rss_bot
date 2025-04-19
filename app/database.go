package app

import (
	"database/sql"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectToDB() *sql.DB {
	db, err := sql.Open("sqlite3", Config.DatabasePath)
	if err != nil {
		panic(err)
	}
	db.Exec(`CREATE TABLE IF NOT EXISTS rss_sources (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT,
		latest_post_datetime DATETIME,
		telegram_chat_id INTEGER,
		keywords TEXT
	)`)
	return db
}

func GetRSSUrl(db *sql.DB) ([]RSSParsedInfo, error) {
	rows, err := db.Query("SELECT url, latest_post_datetime, telegram_chat_id, keywords FROM rss_sources")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var rss_infos []RSSParsedInfo

	for rows.Next() {
		var rss_info RSSParsedInfo
		var joinedKeywords string
		if err := rows.Scan(&rss_info.Url, &rss_info.LastPostDateTime, &rss_info.TelegramChatId, &joinedKeywords); err != nil {
			return rss_infos, err
		}
		rss_info.Keywords = strings.Split(joinedKeywords, "|")
		rss_infos = append(rss_infos, rss_info)
	}

	return rss_infos, nil
}

func checkRSSUrlAlreadyExists(db *sql.DB, url string, telegramChatId int64) bool {
	rows, err := db.Query("SELECT url FROM rss_sources WHERE url = ? AND telegram_chat_id = ?", url, telegramChatId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	return rows.Next()
}

func addRSSUrl(db *sql.DB, url string, telegramChatId int64, keywords []string) {
	joinedKeywords := strings.Join(keywords, "|")
	_, err := db.Exec("INSERT INTO rss_sources (url, latest_post_datetime, telegram_chat_id, keywords) VALUES (?, ?, ?, ?)", url, time.Now(), telegramChatId, joinedKeywords)
	if err != nil {
		panic(err)
	}
}

func deleteRSSUrl(db *sql.DB, url string, telegramChatId int64) {
	_, err := db.Exec("DELETE FROM rss_sources WHERE url = ? AND telegram_chat_id = ?", url, telegramChatId)
	if err != nil {
		panic(err)
	}
}

func SetLatestPostDateTime(db *sql.DB, url string, telegramChatId int64, latestPostDateTime time.Time) {
	_, err := db.Exec("UPDATE rss_sources SET latest_post_datetime = ? WHERE url = ? AND telegram_chat_id = ?", latestPostDateTime, url, telegramChatId)
	if err != nil {
		panic(err)
	}
}
