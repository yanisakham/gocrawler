package main

import (
	"github.com/wangwalton/gocrawler/crawler"
	"github.com/wangwalton/gocrawler/queue"
)

func main() {
	queue := make(queue.Channel, 1000)
	scrape_url := "https://cnn.com"
	queue.Enqueue(scrape_url)
	for {
		u := queue.Dequeue()
		go crawler.ProcessURL(queue, u)
	}
}
