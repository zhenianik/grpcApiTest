package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/zhenianik/grpcApiTest/config"
	"github.com/zhenianik/grpcApiTest/internal"
	"github.com/zhenianik/grpcApiTest/internal/controller/grpc/api"
	grpcservice "github.com/zhenianik/grpcApiTest/internal/grpc"
	eventsender "github.com/zhenianik/grpcApiTest/internal/kafka"
	"github.com/zhenianik/grpcApiTest/internal/postgres"
	cache "github.com/zhenianik/grpcApiTest/pkg/cache/redis"
	logger "github.com/zhenianik/grpcApiTest/pkg/logger/zap"
	"google.golang.org/grpc"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func NewGRPCServer(cfg *config.Config) (*grpcservice.Service, error) {
	es := eventsender.New(cfg.KafkaAddress)
	redisCache := cache.NewRedisCache(cfg.RedisHost, cfg.RedisDb, cfg.RedisExpires)

	l, err := logger.New(cfg.LogLevel)
	if err != nil {
		log.Fatal("could not create logger: %w", err)
	}

	conn, err := pgxpool.Connect(context.Background(), cfg.PostgresUrl)
	if err != nil {
		return nil, fmt.Errorf("could not connect to postgres: %w", err)
	}

	storage := postgres.New(conn)
	userRepo := internal.New(storage, es, redisCache)

	return grpcservice.New(userRepo, l), nil
}

func run() error {
	// config
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("fetching config error: %w", err)
	}

	s := grpc.NewServer()
	srv, err := NewGRPCServer(cfg)
	if err != nil {
		return fmt.Errorf("error creating grpc server: %w", err)
	}
	api.RegisterUsersServer(s, srv)

	l, err := net.Listen(cfg.GrpcNetwork, cfg.GrpcAddress)
	if err != nil {
		return fmt.Errorf("creating listener error: %w", err)
	}

	if err = s.Serve(l); err != nil {
		return fmt.Errorf("serve listener error: %w", err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	fmt.Println("closing")

	return nil
}
