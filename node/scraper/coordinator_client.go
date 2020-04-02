package scraper

import (
	"context"
	"time"

	"github.com/wangwalton/gocrawler/contracts"
)

func GetHostname(client contracts.HostnameCoordinatorClient, u *contracts.ScraperJob) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client.GetHostname(ctx, u)
}

func Dequeue(client contracts.HostnameCoordinatorClient) (u *contracts.ScraperJob) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	job, _ := client.Dequeue(ctx, &contracts.Empty{})
	return job
}
