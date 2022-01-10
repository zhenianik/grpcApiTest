package main

import (
	"fmt"
	"log"
	"net"

	"github.com/zhenianik/grpcApiTest/config"
	"github.com/zhenianik/grpcApiTest/internal/controller/api"
	"github.com/zhenianik/grpcApiTest/internal/usecase"
	"github.com/zhenianik/grpcApiTest/internal/usecase/cache"
	"github.com/zhenianik/grpcApiTest/internal/usecase/repository"
	"github.com/zhenianik/grpcApiTest/pkg/dbLogger"
	"github.com/zhenianik/grpcApiTest/pkg/logger"
	"github.com/zhenianik/grpcApiTest/pkg/postgres"
	"google.golang.org/grpc"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func NewGRPCServer(cfg *config.Config) *usecase.GRPCServer {
	logDB := dbLogger.New(cfg.KafkaAddress)
	db := postgres.Connect(cfg.PostgresUrl)
	redisCache := cache.NewRedisCache(cfg.RedisHost, cfg.RedisDb, cfg.RedisExpires)
	l := logger.New(cfg.LogLevel)
	repo := repository.New(db)
	return usecase.NewGRPCServer(repo, logDB, redisCache, l)
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
