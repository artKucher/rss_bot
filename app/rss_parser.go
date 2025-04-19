package app

import (
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

func GetNewPosts(url string) ([]Post, error) {
	rss_parser := gofeed.NewParser()
	feed, err := rss_parser.ParseURL(url)
	if err != nil {
		return nil, err
	}

	var posts []Post
	for _, item := range feed.Items {
		post := Post{
			Title:         item.Title,
			Description:   item.Description,
			Link:          item.Link,
			PublishedTime: *item.PublishedParsed,
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func FilterPosts(posts []Post, lastPostDateTime time.Time, keywords []string) []Post {
	var filteredPosts []Post

	for _, post := range posts {
		if !post.PublishedTime.After(lastPostDateTime) {
			continue
		}

		if isPostContainsKeywords(post, keywords) {
			filteredPosts = append(filteredPosts, post)
		}
	}
	return filteredPosts
}

func isPostContainsKeywords(post Post, keywords []string) bool {
	if len(keywords) == 0 {
		return true
	}

	for _, keyword := range keywords {
		if strings.Contains(post.Title, keyword) || strings.Contains(post.Description, keyword) {
			return true
		}
	}

	return false
}
