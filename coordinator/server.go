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
	hostnameMap   map[string]*contracts.HostnamePaths
	cv            *sync.Cond
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
	logger_, _ := config.Build()
	return logger_
}

func (s *hostnameCoordinatorServer) GetHostname(ctx context.Context, req *contracts.Empty) (*contracts.HostnamePaths, error) {
	s.cv.L.Lock()
	for len(s.hostnameQueue) == 0 {
		s.cv.Wait()
	}

	h := s.hostnameQueue[0]
	s.hostnameQueue = s.hostnameQueue[1:]

	if _, ok := s.hostnameMap[h.Hostname]; ok {
		delete(s.hostnameMap, h.Hostname)
	} else {
		zap.S().Errorf("Can't find %s in map", h.Hostname)
	}

	s.cv.L.Unlock()
	return h, nil
}

func (s *hostnameCoordinatorServer) AddHostnames(ctx context.Context, req *contracts.MultipleHostnamePaths) (*contracts.Empty, error) {
	s.cv.L.Lock()
	for _, hp := range req.Urls {
		if len(s.hostnameQueue) == 0 {
			s.cv.Signal()
		}

		if val, ok := s.hostnameMap[hp.Hostname]; ok {
			for key := range hp.Paths {
				if _, ok = val.Paths[key]; !ok {
					val.Paths[key] = &contracts.Empty{}
				}
			}
		} else {
			s.hostnameMap[hp.Hostname] = hp
			s.hostnameQueue = append(s.hostnameQueue, hp)
		}
	}
	s.cv.L.Unlock()
	return nil, nil
}

func newServer() *hostnameCoordinatorServer {
	return &hostnameCoordinatorServer{
		hostnameQueue: make([]*contracts.HostnamePaths, 0),
		hostnameMap:   make(map[string]*contracts.HostnamePaths),
		cv:            sync.NewCond(&sync.Mutex{}),
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
