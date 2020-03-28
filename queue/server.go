package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"sync"

	"github.com/wangwalton/gocrawler/contracts"
	"google.golang.org/grpc"
)

type urlScraperServer struct {
	contracts.UnimplementedURLQueueServer
	urlQueue chan contracts.ScraperJob

	mu sync.Mutex
}

var (
	// tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	// certFile   = flag.String("cert_file", "", "The TLS cert file")
	// keyFile    = flag.String("key_file", "", "The TLS key file")
	// jsonDBFile = flag.String("json_db_file", "", "A json file containing a list of features")
	port = flag.Int("port", 10000, "The server port")
)

func (s *urlScraperServer) Dequeue(ctx context.Context, req *contracts.Empty) (*contracts.ScraperJob, error) {
	job := <-s.urlQueue
	fmt.Printf("Receved dequeue request, sending %s\n", job.Url)
	return &job, nil
}
func (s *urlScraperServer) Enqueue(ctx context.Context, req *contracts.ScraperJob) (*contracts.Empty, error) {
	fmt.Printf("Received enqueue request of %s\n", req.Url)
	s.urlQueue <- *req
	return &contracts.Empty{}, nil
}

func newServer() *urlScraperServer {
	return &urlScraperServer{urlQueue: make(chan contracts.ScraperJob, 1000000)}
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	contracts.RegisterURLQueueServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
