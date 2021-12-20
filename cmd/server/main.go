package main

import (
	"fmt"
	"github.com/zhenianik/grpcApiTest/config"
	"github.com/zhenianik/grpcApiTest/internal/user"
	"github.com/zhenianik/grpcApiTest/internal/user/cache"
	"github.com/zhenianik/grpcApiTest/pkg/api"
	"github.com/zhenianik/grpcApiTest/pkg/database"
	"github.com/zhenianik/grpcApiTest/pkg/logger"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func NewGRPCServer(cfg *config.Config) *user.GRPCServer {
	logDB := logger.New(cfg.KafkaAddress)
	db := database.Connect(cfg.PostgresUrl)
	cache := cache.NewRedisCache(cfg.RedisHost, cfg.RedisDb, cfg.RedisExpires)
	return user.NewGRPCServer(db, logDB, cache)
}

func run() error {
	// config
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("ошибка получения конфига: %w", err)
	}

	s := grpc.NewServer()
	srv := NewGRPCServer(cfg)
	api.RegisterUserServer(s, srv)

	l, err := net.Listen(cfg.GrpcNetwork, cfg.GrpcAddress)
	if err != nil {
		log.Fatal(err)
	}

	if err = s.Serve(l); err != nil {
		return err
	}

	fmt.Printf("listening %s port %s", cfg.GrpcNetwork, cfg.GrpcAddress)

	return nil
}
