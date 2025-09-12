package main

import (
	"github.com/iosifache/sec-feeds/internal/feeds"
	"github.com/iosifache/sec-feeds/internal/logger"
)

func main() {
	logger.InitLogger()
	feeds.GenerateAllFeeds()
}
