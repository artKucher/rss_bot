package app

import "time"

type RSSParsedInfo struct {
	Url              string
	LastPostDateTime time.Time
	TelegramChatId   int64
	Keywords         []string
}

type Post struct {
	Title         string
	Description   string
	Link          string
	PublishedTime time.Time
}
