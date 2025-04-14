package main

import "time"

type RSSParsedInfo struct {
	url              string
	lastPostDateTime time.Time
	telegramChatId   int64
	keywords         []string
}

type Post struct {
	Title         string
	Description   string
	Link          string
	PublishedTime time.Time
}
