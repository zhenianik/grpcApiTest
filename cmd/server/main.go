package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/zhenianik/grpcApiTest/config"
	"github.com/zhenianik/grpcApiTest/internal/controller/grpc/api"
	"github.com/zhenianik/grpcApiTest/internal/usecase"
	eventsender "github.com/zhenianik/grpcApiTest/internal/usecase/kafka"
	"github.com/zhenianik/grpcApiTest/internal/usecase/repository"
	cache "github.com/zhenianik/grpcApiTest/pkg/cache/redis"
	logger "github.com/zhenianik/grpcApiTest/pkg/logger/zap"
	"google.golang.org/grpc"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func NewGRPCServer(cfg *config.Config) *usecase.Service {
	es := eventsender.New(cfg.KafkaAddress)
	redisCache := cache.NewRedisCache(cfg.RedisHost, cfg.RedisDb, cfg.RedisExpires)

	l, err := logger.New(cfg.LogLevel)
	if err != nil {
		log.Fatal("could not create logger: %w", err)
		return nil
	}

	conn, err := pgxpool.Connect(context.Background(), cfg.PostgresUrl)
	if err != nil {
		l.Error("could not connect to postgres: %w", err)
		return nil
	}

	repo := repository.New(conn)
	return usecase.NewService(repo, es, redisCache, l)
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
