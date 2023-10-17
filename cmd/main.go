package main

import (
	"log"

	scraper "github.com/knrd/scraper_task"
	"github.com/knrd/scraper_task/cache"
	"github.com/knrd/scraper_task/queue"
)

func main() {
	urls := []string{
		"http://example.com",
		"http://example.net",
		"http://example.org",
	}

	cache := cache.New()
	queue := queue.New(urls)
	out := scraper.Run(queue, 2, cache)

	log.Println(out)
}
