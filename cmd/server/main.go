package main

import (
	"fmt"
	"github.com/zhenianik/grpcApiTest/config"
	"github.com/zhenianik/grpcApiTest/internal/user"
	"github.com/zhenianik/grpcApiTest/internal/user/cache"
	"github.com/zhenianik/grpcApiTest/pkg/api"
	"github.com/zhenianik/grpcApiTest/pkg/database"
	"github.com/zhenianik/grpcApiTest/pkg/dbLogger"
	"github.com/zhenianik/grpcApiTest/pkg/logger"
	"google.golang.org/grpc"
	"net"
)

func main() {
	if err := run(); err != nil {
		logger.Logger.Fatal(err)
	}
}

func NewGRPCServer(cfg *config.Config) *user.GRPCServer {
	logDB := dbLogger.New(cfg.KafkaAddress)
	db := database.Connect(cfg.PostgresUrl)
	cache := cache.NewRedisCache(cfg.RedisHost, cfg.RedisDb, cfg.RedisExpires)
	logger := logger.NewLogger(cfg.LogLevel)
	return user.NewGRPCServer(db, logDB, cache, logger)
}

func run() error {
	// config
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("fetching config error: %w", err)
	}

	s := grpc.NewServer()
	srv := NewGRPCServer(cfg)
	api.RegisterUserServer(s, srv)

	l, err := net.Listen(cfg.GrpcNetwork, cfg.GrpcAddress)
	if err != nil {
		return fmt.Errorf("creating listener error: %w", err)
	}

	if err = s.Serve(l); err != nil {
		return fmt.Errorf("serve listener error: %w", err)
	}

	return nil
}
