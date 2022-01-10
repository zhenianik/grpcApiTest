package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/zhenianik/grpcApiTest/pkg/logger"
)

func Connect(url string) *pgxpool.Pool {
	conn, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		logger.Logger.Error("could not connect to postgres: %w", err)
		return nil
	}

	return conn
}
