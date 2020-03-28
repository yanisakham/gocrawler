package main

import (
	"flag"
	"os"

	"github.com/wangwalton/gocrawler/contracts"
	"github.com/wangwalton/gocrawler/scraper/scraper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

var (
	// tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	// caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr = flag.String("server_addr", "localhost:10000", "The server address in the format of host:port")
	// serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name use to verify the hostname returned by TLS handshake")
)

var logger, err = zap.NewDevelopment()

func initZapLog() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, _ := config.Build()
	return logger
}

func main() {
	loggerMgr := initZapLog()
	zap.ReplaceGlobals(loggerMgr)
	defer loggerMgr.Sync() // flushes buffer, if any
	logger := loggerMgr.Sugar()

	flag.Parse()
	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure())
	if err != nil {
		logger.Info("failed to dial: %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Seeding with Jobs
	scrape_url := "https://cnn.com"
	client := contracts.NewURLQueueClient(conn)
	scraper.Enqueue(client, &contracts.ScraperJob{Url: scrape_url, Requeue: false})

	for {
		job := scraper.Dequeue(client)
		scraper.ProcessURL(client, job.Url)
		logger.Debugf("Popped job of %s\n", job.Url)
	}
}
