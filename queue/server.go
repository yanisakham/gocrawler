package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/wangwalton/gocrawler/contracts"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

type urlScraperServer struct {
	contracts.UnimplementedURLQueueServer
	urlQueue chan contracts.ScraperJob
	visited  map[string]bool
	mu       sync.Mutex // Protects visited
}

var (
	// tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	// certFile   = flag.String("cert_file", "", "The TLS cert file")
	// keyFile    = flag.String("key_file", "", "The TLS key file")
	// jsonDBFile = flag.String("json_db_file", "", "A json file containing a list of features")
	port = flag.Int("port", 10000, "The server port")
)

func initZapLog() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, _ := config.Build()
	return logger
}

func (s *urlScraperServer) isVisited(job *contracts.ScraperJob) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.visited[job.Url] {
		return true
	} else {
		s.visited[job.Url] = true
		return false
	}
}

func (s *urlScraperServer) Dequeue(ctx context.Context, req *contracts.Empty) (*contracts.ScraperJob, error) {
	job := <-s.urlQueue
	zap.S().Debugf("Receved dequeue request, sending %s\n", job.Url)
	return &job, nil
}

func (s *urlScraperServer) Enqueue(ctx context.Context, req *contracts.ScraperJob) (*contracts.Empty, error) {

	if req.Requeue || !s.isVisited(req) {
		// fmt.Printf("Enqueuing %s\n", req.Url)
		s.urlQueue <- *req
	} else {
		// fmt.Printf("Rejected %s\n", req.Url)
	}
	return &contracts.Empty{}, nil

}

func newServer() *urlScraperServer {
	return &urlScraperServer{
		urlQueue: make(chan contracts.ScraperJob, 1000000),
		visited:  make(map[string]bool),
	}
}

func main() {
	loggerMgr := initZapLog()
	zap.ReplaceGlobals(loggerMgr)
	defer loggerMgr.Sync()
	logger := loggerMgr.Sugar()

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		logger.Infof("failed to listen: %v", err)
		os.Exit(1)
	}
	grpcServer := grpc.NewServer()
	contracts.RegisterURLQueueServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
