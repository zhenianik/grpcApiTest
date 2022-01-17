package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/zhenianik/grpcApiTest/internal"
)

type UserStorage struct {
	pg *pgxpool.Pool
}

func New(db *pgxpool.Pool) *UserStorage {
	return &UserStorage{
		pg: db,
	}
}

func (db *UserStorage) GetUsers(ctx context.Context) (users []*internal.User, err error) {
	conn, err := db.pg.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching db connect: %w", err)
	}

	rows, err := conn.Query(ctx, "SELECT id, name, email FROM users")
	if err != nil {
		return nil, fmt.Errorf("select from db error")
	}
	for rows.Next() {
		var user internal.User

		err = rows.Scan(&user.Id, &user.Name, &user.Email)

		if err != nil {
			return nil, fmt.Errorf("error scaning row result: %w", err)
		}

		users = append(users, &user)
	}

	return users, nil
}

func (db *UserStorage) AddUser(ctx context.Context, user *internal.User) (internal.UserID, error) {
	tx, err := db.pg.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("db communication error: %w", err)
	}

	defer tx.Rollback(ctx)

	rows, _ := tx.Query(ctx, "SELECT id FROM users WHERE name = $1", user.Name)
	if rows.Next() {
		return 0, fmt.Errorf("service with name %s already exists", user.Name)
	}

	var id int64
	err = tx.QueryRow(ctx, "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", user.Name, user.Email).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("db communication error: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("db communication error: %w", err)
	}

	return internal.UserID(id), nil
}

func (db *UserStorage) RemoveUser(ctx context.Context, id internal.UserID) error {
	conn, err := db.pg.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("error fetching db connect: %w", err)
	}
	rows, _ := conn.Query(ctx, "SELECT id FROM users WHERE id = $1", id)
	if !rows.Next() {
		return fmt.Errorf("service with id %d doesn't exist", id)
	}

	tx, err := db.pg.Begin(ctx)
	if err != nil {
		return fmt.Errorf("db communication error: %w", err)
	}

	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "DELETE FROM users WHERE id = $1", id)

	if err != nil {
		return fmt.Errorf("db communication error: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("db communication error: %w", err)

	}

	return nil
}
