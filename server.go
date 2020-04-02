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
	"google.golang.org/grpc/credentials"
)

type hostnameCoordinatorServer struct {
	contracts.UnimplementedHostnameCoordinatorServer
	hostnameQueue []*contracts.HostnamePaths
	mu sync.Mutex
}


var (
	crt = "ssl/server.crt"
	key = "ssl/server.key"
)

var (
	tls      = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile = flag.String("cert_file", "ssl/server.crt", "The TLS cert file")
	keyFile  = flag.String("key_file", "ssl/server.key", "The TLS key file")
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

func (s *hostnameCoordinatorServer) GetHostname(ctx context.Context, req *contracts.Empty) (*contracts.HostnamePaths, error) {
	return nil, nil
}

func (s *hostnameCoordinatorServer) AddHostnames(ctx context.Context,
	req *contracts.MultipleHostnamePaths) (*contracts.Empty, error) {
	return nil, nil
}

func newServer() *hostnameCoordinatorServer {
	return &hostnameCoordinatorServer{
		hostnameQueue: make([]*contracts.HostnamePaths, 0),
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
	var opts []grpc.ServerOption
	if *tls {
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			logger.Infof("Failed to generate credentials %v", err)
			os.Exit(1)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)
	contracts.RegisterHostnameCoordinatorServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
