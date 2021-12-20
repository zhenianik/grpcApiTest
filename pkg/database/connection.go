package database

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

func Connect(url string) *pgxpool.Pool {
	conn, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err.Error())
		return nil
	}

	return conn
}
