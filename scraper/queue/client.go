package queue

import (
	"context"
	"time"

	"github.com/wangwalton/gocrawler/contracts"
)

func Enqueue(client contracts.URLQueueClient, u *contracts.ScraperJob) {
	//fmt.Printf("Enqueuing %s\n", u.Url)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client.Enqueue(ctx, u)
}

func Dequeue(client contracts.URLQueueClient) (u *contracts.ScraperJob) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	job, _ := client.Dequeue(ctx, &contracts.Empty{})
	//fmt.Printf("Dequed %s\n", job.Url)
	return job
}
