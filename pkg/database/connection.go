package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/zhenianik/grpcApiTest/pkg/logger"
)

func Connect(url string) *pgxpool.Pool {
	conn, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		logger.Logger.Error(fmt.Errorf("could not connect to database: %w", err))
		return nil
	}

	return conn
}
