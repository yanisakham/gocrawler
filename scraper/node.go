package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/wangwalton/gocrawler/contracts"
	"github.com/wangwalton/gocrawler/scraper/crawler"
	"github.com/wangwalton/gocrawler/scraper/queue"
	"google.golang.org/grpc"
)

var (
	// tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	// caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr = flag.String("server_addr", "localhost:10000", "The server address in the format of host:port")
	// serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name use to verify the hostname returned by TLS handshake")
)

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("faile to dial: %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Seeding with Jobs
	scrape_url := "https://cnn.com"
	client := contracts.NewURLQueueClient(conn)
	queue.Enqueue(client, &contracts.ScraperJob{Url: scrape_url})

	for {
		job := queue.Dequeue(client)
		go crawler.ProcessURL(client, job.Url)
		// fmt.Printf("Popped job of %s\n", job.Url)
	}
	// queue.Enqueue(scrape_url)
	// for {
	// 	u := queue.Dequeue()
	// 	go crawler.ProcessURL(queue, u)
	// }
}
